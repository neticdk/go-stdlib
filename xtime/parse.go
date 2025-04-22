package xtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// DurationParseError provides specific details about a duration parsing failure.
type DurationParseError struct {
	// The original input string that caused the error.
	Input string
	// The specific token that caused the error, if applicable.
	Token string
	// The index of the problematic token, if applicable.
	TokenIndex int
	// Description of the error.
	Message string
}

func (e *DurationParseError) Error() string {
	if e.TokenIndex >= 0 && e.Token != "" {
		return fmt.Sprintf("xtime: %s (token='%s' at index %d in input='%s')", e.Message, e.Token, e.TokenIndex, e.Input)
	}
	return fmt.Sprintf("xtime: %s (input='%s')", e.Message, e.Input)
}

func newParseError(input, message string) error {
	return &DurationParseError{Input: input, Message: message}
}

func newParseTokenError(input, token, message string, index int) error {
	return &DurationParseError{Input: input, Token: token, TokenIndex: index, Message: message}
}

// ParseDuration converts a human-readable duration string into
// a time.Duration.
//
// It first cleans the input string using the clean() function (removing extra
// words like "and", replacing commas with spaces, trimming whitespace).
// Then, it tokenizes the cleaned string using the tokenize() function, which
// splits the string by spaces and also by transitions between numbers and
// letters (e.g., "1h30m" becomes ["1", "h", "30", "m"]).
//
// The function handles the following cases:
//   - Empty or whitespace-only input: Returns an error.
//   - Single token input: Attempts to parse it as a float64 representing
//     seconds. Returns an error if the token is not a valid number.
//   - Multiple tokens: Expects an even number of tokens representing pairs of
//     (value, unit). It iterates through these pairs, parses the value as a
//     float64, looks up the unit in the predefined units map, and accumulates
//     the total duration.
//
// It returns an error if:
//   - The input is empty after cleaning.
//   - A single token cannot be parsed as a number.
//   - There is an odd number of tokens (expecting pairs).
//   - A token expected to be a value cannot be parsed as a float64.
//   - A token expected to be a unit is not found in the units map.
func ParseDuration(s string) (time.Duration, error) {
	return parse(s)
}

var units = map[string]time.Duration{
	"ns":           time.Nanosecond,
	"nanosecond":   time.Nanosecond,
	"nanoseconds":  time.Nanosecond,
	"μs":           time.Microsecond,
	"us":           time.Microsecond,
	"microsecond":  time.Microsecond,
	"microseconds": time.Microsecond,
	"ms":           time.Millisecond,
	"millisecond":  time.Millisecond,
	"milliseconds": time.Millisecond,
	"s":            time.Second,
	"sec":          time.Second,
	"secs":         time.Second,
	"second":       time.Second,
	"seconds":      time.Second,
	"m":            time.Minute,
	"min":          time.Minute,
	"mins":         time.Minute,
	"minute":       time.Minute,
	"minutes":      time.Minute,
	"h":            time.Hour,
	"hr":           time.Hour,
	"hour":         time.Hour,
	"hours":        time.Hour,
	"d":            24 * time.Hour,
	"day":          24 * time.Hour,
	"days":         24 * time.Hour,
	"w":            7 * 24 * time.Hour,
	"week":         7 * 24 * time.Hour,
	"weeks":        7 * 24 * time.Hour,
	"mo":           MonthApprox, // Approximate
	"month":        MonthApprox, // Approximate
	"months":       MonthApprox, // Approximate
	"y":            YearApprox,  // Approximate
	"year":         YearApprox,  // Approximate
	"years":        YearApprox,  // Approximate
}

func parse(input string) (time.Duration, error) {
	cleanedInput := clean(input)
	tokens := tokenize(cleanedInput)

	var duration time.Duration

	if len(tokens) == 0 {
		return 0, newParseError(input, "string resulted in zero tokens after cleaning")
	}

	// Handle single number input (interpreted as seconds)
	if len(tokens) == 1 {
		token := tokens[0]
		seconds, err := strconv.ParseFloat(token, floatBitSize)
		if err == nil {
			duration = time.Duration(seconds * float64(time.Second))
			return duration, nil
		}
		msg := "single token is not a valid number (expected seconds)"
		return 0, newParseTokenError(input, token, msg, 0)

	}

	// Expect pairs of (value, unit) from here on
	if len(tokens)%2 != 0 {
		msg := fmt.Sprintf("expected pairs of value and unit, but got %d tokens", len(tokens))
		return 0, newParseError(input, msg)
	}

	for i := 0; i < len(tokens); i += 2 {
		valueStr := tokens[i]
		unitStr := tokens[i+1]

		value, err := strconv.ParseFloat(valueStr, floatBitSize)
		if err != nil {
			msg := "expected a number"
			return 0, newParseTokenError(input, valueStr, msg, i)
		}

		unitMultiplier, ok := units[unitStr]
		if !ok {
			msg := "unknown unit"
			return 0, newParseTokenError(input, unitStr, msg, i+1)
		}

		duration += time.Duration(value * float64(unitMultiplier))
	}

	return duration, nil
}

// clean cleans the input string
func clean(input string) string {
	s := strings.TrimSpace(input)
	s = strings.ReplaceAll(s, "and ", "")
	s = strings.ReplaceAll(s, ",", " ")

	return s
}

const (
	typeUnknown = iota
	typeDigitOrDecimal
	typeLetter
)

// tokenize splits a cleaned duration string into potential number and unit
// tokens.
//
// It operates in two stages:
// 1. Splits the input string by whitespace into fields.
// 2. Processes each field:
//   - If a field starts with '.', a '0' is prepended (e.g., ".5h" -> "0.5h").
//   - The field is then split internally whenever the character type changes
//     between a digit/decimal point and a letter. For example:
//     "10years" -> ["10", "years"]
//     "1h30m" -> ["1", "h", "30", "m"]
//     "h1" -> ["h", "1"]
//
// This function does not perform validation on the content or sequence of
// tokens (e.g., it doesn't check if "1.2.3" is a valid number or if "xyz" is
// a known unit, or if tokens are in number-unit pairs). It aims to never
// return an error, deferring all validation to the caller (typically the parse
// function). Input is assumed to be already cleaned (e.g., no commas, extra
// words).
func tokenize(cleanedInput string) []string {
	fields := strings.Fields(cleanedInput)
	finalTokens := make([]string, 0, len(fields))

	for _, field := range fields {
		var currentToken strings.Builder
		lastType := typeUnknown

		// Handle leading decimal for the whole field ".5h" -> "0.5h"
		// Or ".h" -> "0.h"
		processedField := field
		if field[0] == '.' {
			processedField = "0" + field
		}

		for i, r := range processedField {
			currentType := typeUnknown
			if unicode.IsDigit(r) || r == '.' {
				currentType = typeDigitOrDecimal
			} else if unicode.IsLetter(r) {
				currentType = typeLetter
			}

			// Start of the first token in the field
			if i == 0 {
				currentToken.WriteRune(r)
				lastType = currentType
				continue
			}

			// If type changes between Letter and Digit/Decimal, split
			if (lastType == typeLetter && currentType == typeDigitOrDecimal) ||
				(lastType == typeDigitOrDecimal && currentType == typeLetter) {
				// Add the completed token
				finalTokens = append(finalTokens, currentToken.String())
				// Start the new token
				currentToken.Reset()
			}

			currentToken.WriteRune(r)
			lastType = currentType
		}

		// Add the last token from the field
		if currentToken.Len() > 0 {
			finalTokens = append(finalTokens, currentToken.String())
		}
	}

	return finalTokens
}
