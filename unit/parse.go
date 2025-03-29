package unit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseError represents an error that occurred during parsing
type ParseError struct {
	Input string
	Msg   string
	Pos   int // position in input where error occurred
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at position %d: %s (input: %q)", e.Pos, e.Msg, e.Input)
}

// ParseResult represents a parsed value with its unit
type ParseResult struct {
	// Raw is the numeric value parsed directly from the input string, before
	// applying the prefix scale factor.
	Raw float64
	// Scale is the numeric scale factor determined by the parsed prefix (for
	// example, 1 for base unit, 1000 for Kilo, 1024 for Kibi).
	Scale float64
	// Unit is the base unit type identified (for example, Byte, Bit, Hertz, or
	// a custom unit). Will be None if only a number or number+prefix was
	// parsed.
	Unit Unit
}

// Regex to find leading number (adjust as needed for scientific notation etc.)
var numberRegex = regexp.MustCompile(`^([+-]?\s*[0-9]+(?:\.[0-9]*)?|\.[0-9]+)\s*(.*)$`)

// Parse attempts to parse a string containing a value with units.
// It handles standard and registered custom units and prefixes.
//
// Examples of valid formats:
//   - "1.5 KB"    -> {Raw: 1.5, Scale: 1000, Unit: Byte}
//   - "1 MiB"     -> {Raw: 1, Scale: 1048576, Unit: Byte}
//   - "500 MHz"   -> {Raw: 500, Scale: 1000000, Unit: Hertz}
//   - "2.7 kb"    -> {Raw: 2.7, Scale: 1000, Unit: Bit}
//   - "1024"      -> {Raw: 1024, Scale: 1, Unit: None}
func Parse(s string) (ParseResult, error) {
	// Validation
	s = strings.TrimSpace(s)
	if s == "" {
		return ParseResult{}, &ParseError{Input: s, Msg: "empty input"}
	}

	// Extract number part and remainder of string
	rawValue, remainder, err := parseNumberAndRemainder(s)
	if err != nil {
		return ParseResult{}, err
	}

	// Handle number-only case
	if remainder == "" {
		return ParseResult{Raw: rawValue, Scale: 1, Unit: None}, nil
	}

	// Remainder Parsing Logic

	// Acquire locks needed for the matching helpers
	unitRegistryMutex.RLock()
	formatSystemsMutex.RLock()
	defer unitRegistryMutex.RUnlock()
	defer formatSystemsMutex.RUnlock()

	// Strategy 1: Try matching Prefix + Unit (longest prefix first)
	scale, unit, found := matchPrefixAndUnit(remainder)
	if found {
		return ParseResult{Raw: rawValue, Scale: scale, Unit: unit}, nil
	}

	// Strategy 2: Try matching Unit Only
	unit, found = matchUnitOnly(remainder)
	if found {
		// Scale is 1 (base unit) when only a unit is found
		return ParseResult{Raw: rawValue, Scale: 1, Unit: unit}, nil
	}

	// Strategy 3: Try matching Prefix Only
	scale, found = matchPrefixOnly(remainder)
	if found {
		// Unit is None when only a prefix is found
		return ParseResult{Raw: rawValue, Scale: scale, Unit: None}, nil
	}

	// If none of the strategies worked, the remainder is unrecognized
	// Calculate position of the remainder for the error message
	remainderPos := 0
	if idx := strings.Index(s, remainder); idx != -1 {
		remainderPos = idx
	}
	return ParseResult{}, &ParseError{
		Input: s, Msg: fmt.Sprintf("unknown unit or prefix combination: %q", remainder), Pos: remainderPos,
	}
}

// MustParse is a convenience function that parses a string and panics if an
// error occurs.
func MustParse(s string) ParseResult {
	result, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return result
}

// Convenience methods for Result
func (r ParseResult) Value() float64 {
	return r.Raw * r.Scale
}

// Format formats the value with the given options.
func (r ParseResult) Format(opts ...FormatOption) Value {
	return Format(r.Value(), append(opts, WithUnit(r.Unit))...)
}

// parseNumberAndRemainder extracts the initial numeric value and the rest of
// the string. It returns the parsed float, the remainder, and an error if
// parsing fails.
func parseNumberAndRemainder(s string) (rawValue float64, remainder string, err error) {
	matches := numberRegex.FindStringSubmatch(s)
	var numberStr string

	if matches != nil {
		// Found number and remainder using regex
		numberStr = strings.ReplaceAll(matches[1], " ", "") // Remove spaces like in "+ 1.5"
		remainder = strings.TrimSpace(matches[2])
		rawValue, err = strconv.ParseFloat(numberStr, float64BitSize)
		if err != nil {
			// Wrap strconv error, provide position if possible
			pos := strings.Index(s, numberStr)
			if pos == -1 {
				pos = 0
			}
			err = &ParseError{Input: s, Msg: fmt.Sprintf("invalid number: %v", err), Pos: pos}
			return
		}
	} else {
		// Didn't match number+remainder regex, try parsing the whole string as a number
		rawValue, err = strconv.ParseFloat(s, float64BitSize)
		if err != nil {
			// Not a valid number alone either
			err = &ParseError{Input: s, Msg: "invalid format - no valid number found", Pos: 0}
			return
		}
		// It was just a number, no remainder
		remainder = ""
	}

	return
}

// matchPrefixAndUnit attempts to match the remainder as a known prefix followed
// by a known unit. It requires read locks on the registries to be held by the
// caller.
func matchPrefixAndUnit(remainder string) (scale float64, unit Unit, found bool) {
	for _, prefixKey := range sortedPrefixKeys { // Iterate sorted prefixes (longest first)
		if strings.HasPrefix(remainder, prefixKey) {
			potentialUnitPart := strings.TrimSpace(remainder[len(prefixKey):])

			// Allow matching just a prefix if the remainder *is* the prefixKey
			if potentialUnitPart == "" {
				// This case (prefix only) is handled later by matchPrefixOnly,
				// but checking here avoids unnecessary unit lookups below.
				// Continue to the next shorter prefixKey if any.
				// We prioritize Unit-only match over Prefix-only match later.
				continue // Let matchPrefixOnly handle this later
			}

			// Check if the potentialUnitPart matches a known unit
			var foundUnit Unit
			var unitMatch bool
			if u, ok := unitSymbolLookup[potentialUnitPart]; ok {
				foundUnit = u
				unitMatch = true
			} else if u, ok := unitSingularLookup[potentialUnitPart]; ok {
				foundUnit = u
				unitMatch = true
			} else if u, ok := unitPluralLookup[potentialUnitPart]; ok {
				foundUnit = u
				unitMatch = true
			}

			if unitMatch {
				// Success: Found Prefix + Unit
				scale = prefixLookup[prefixKey] // Get scale for the matched prefixKey
				unit = foundUnit
				found = true
				return // Return the first successful Prefix+Unit match
			}
			// If potentialUnitPart didn't match a unit, continue the loop
			// to try shorter/other prefixes that might also be prefixes of the
			// remainder.
		}
	}

	// No Prefix+Unit combination found
	return
}

// matchUnitOnly attempts to match the entire remainder as a known unit (symbol,
// singular, or plural). It requires read lock on the unit registry to be held
// by the caller.
func matchUnitOnly(remainder string) (unit Unit, found bool) {
	if u, ok := unitSymbolLookup[remainder]; ok {
		unit = u
		found = true
	} else if u, ok := unitSingularLookup[remainder]; ok {
		unit = u
		found = true
	} else if u, ok := unitPluralLookup[remainder]; ok {
		unit = u
		found = true
	}
	return
}

// matchPrefixOnly attempts to match the entire remainder as a known prefix
// (symbol or name). It requires read lock on the format systems registry to be
// held by the caller.
func matchPrefixOnly(remainder string) (scale float64, found bool) {
	scale, found = prefixLookup[remainder]
	return
}
