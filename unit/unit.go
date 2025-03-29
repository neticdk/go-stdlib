package unit

import (
	"math"
	"strconv"
)

// Unit represents the type of unit.
type Unit int

const (
	None           Unit = iota // Raw number, no unit
	Byte                       // Bytes (B)
	Bit                        // Bits (b)
	Hertz                      // Hertz (Hz)
	unitMaxBuiltin             // Used internally to track custom units

	// magic numbers
	float64BitSize = 64 // used with strconv functions
)

// String returns the name of the unit (for stringer interface)
func (u Unit) String() string {
	return Describe(u).Singular
}

// Info returns information about a unit
func (u Unit) Info() Descriptor {
	return Describe(u)
}

const (
	// IEC binary prefixes as defined in IEC 80000-13
	Kibi  float64 = 1024
	Mebi  float64 = 1048576
	Gibi  float64 = 1073741824
	Tebi  float64 = 1099511627776
	Pebi  float64 = 1125899906842624
	Exbi  float64 = 1152921504606846976
	Zebi  float64 = 1180591620717411303424
	Yobi  float64 = 1208925819614629174706176
	Robi  float64 = 1237940039285380274899124224
	Quebi float64 = 1267650600228229401496703205376

	// SI decimal prefixes as defined in ISO 1000
	Kilo   float64 = 1000
	Mega   float64 = 1000000
	Giga   float64 = 1000000000
	Tera   float64 = 1000000000000
	Peta   float64 = 1000000000000000
	Exa    float64 = 1000000000000000000
	Zetta  float64 = 1000000000000000000000
	Yotta  float64 = 1000000000000000000000000
	Ronna  float64 = 1000000000000000000000000000
	Quetta float64 = 1000000000000000000000000000000
)

var (
	// SIDecimalBoundaries represents the boundaries for SI decimal prefixes
	// sorted in descending order.
	SIDecimalBoundaries = []float64{Quetta, Ronna, Yotta, Zetta, Exa, Peta, Tera, Giga, Mega, Kilo}

	// IECBinaryBoundaries represents the boundaries for IEC binary prefixes
	// sorted in descending order.
	IECBinaryBoundaries = []float64{Quebi, Robi, Yobi, Zebi, Exbi, Pebi, Tebi, Gibi, Mebi, Kibi}

	// DataUnits contains units related to data storage/transfer
	DataUnits = []Unit{Byte, Bit}

	// FrequencyUnits contains units related to frequency
	FrequencyUnits = []Unit{Hertz}

	// AllUnits contains all built-in units
	AllUnits = []Unit{Byte, Bit, Hertz}
)

// number represents a numeric type that can be used for unit conversions.
type number interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// Descriptor contains information about how a unit is represented textually
type Descriptor struct {
	// Unit is the identifier for this unit type, returned by Register and used
	// internally.
	Unit Unit
	// Symbol is the short textual representation (e.g., "B", "b", "Hz"). Used
	// in short format strings and parsing. Can be empty.
	Symbol string
	// Singular is the full textual representation (e.g., "byte", "bit",
	// "hertz"). Used in long singular format strings and parsing. Can be empty.
	Singular string
	// Plural is the full textual representation for quantities other than 1
	// (e.g., "bytes", "bits", "hertz"). Used in long plural format strings and
	// parsing. Can be empty (often is if Singular == Plural, like "hertz"). If
	// Singular is empty, Plural may act as a suffix (e.g., custom time unit
	// "s").
	Plural string
}

// unitRegistry holds the mapping between Unit identifiers and their textual
// representations (Descriptor). It stores both the built-in units (like Byte,
// Bit, Hertz) and any custom units added via Register/MustRegister.
//
// Access to this map is synchronized via unitRegistryMutex, as it's read and
// written concurrently by registration and lookup functions (Describe, List,
// Names, and the internal reverse lookup builders used by Parse).
// It is initialized with the standard built-in units.
var unitRegistry = map[Unit]Descriptor{
	Byte:  {Byte, "B", "byte", "bytes"},
	Bit:   {Bit, "b", "bit", "bits"},
	Hertz: {Hertz, "Hz", "hertz", "hertz"},
}

// Describe returns the unit descriptor for the given unit.
func Describe(unit Unit) Descriptor {
	unitRegistryMutex.RLock()
	defer unitRegistryMutex.RUnlock()

	if descriptor, found := unitRegistry[unit]; found {
		return descriptor
	}
	return Descriptor{}
}

// Names returns the names of all registered units
func Names() []string {
	unitRegistryMutex.RLock()
	defer unitRegistryMutex.RUnlock()

	names := make([]string, 0, len(unitRegistry))
	for _, descriptor := range unitRegistry {
		names = append(names, descriptor.Singular)
	}
	return names
}

// List returns information about all registered units
func List() []Descriptor {
	unitRegistryMutex.RLock()
	defer unitRegistryMutex.RUnlock()

	infos := make([]Descriptor, 0, len(unitRegistry))
	for u, descriptor := range unitRegistry {
		infos = append(infos, Descriptor{
			Unit:     u,
			Symbol:   descriptor.Symbol,
			Singular: descriptor.Singular,
			Plural:   descriptor.Plural,
		})
	}
	return infos
}

// Value represents a formatted unit value.
type Value struct {
	// Scaled is the numeric value after being adjusted (divided) by the
	// appropriate prefix scale (for example, Kilo, Mebi) selected during
	// formatting. The formatter chooses the prefix based on the magnitude of
	// the original value and the chosen format options (like Binary() or
	// Decimal()). For/ example, if the input was 1024 and Binary() format was
	// used, Scaled would be 1.0.
	Scaled float64

	// Prefix contains the details (Name and Symbol) of the unit prefix (like
	// "kilo", "k" or "Mebi", "Mi") that corresponds to the scaling factor
	// applied to the original value to get the Scaled value. If the original
	// value was too small to warrant a prefix, or if formatting options
	// disabled prefixes, this will be an empty Prefix struct (`{}`),
	// indicating the base unit scale (1.0) was used.
	Prefix Prefix

	// Unit holds the descriptive information (Symbol, Singular, Plural names)
	// for the fundamental base unit type (for example, Byte, Bit, Hertz, or a
	// custom registered unit). This descriptor is determined by the `WithUnit`
	// option passed to the formatting function or inferred from specific format
	// styles (like `DecimalBytesFormat`). It always represents the base unit,
	// not the combined prefix and unit (for example, it represents "byte", even
	// if the Prefix is "Kilo").
	Unit Descriptor

	// formatOptions holds the options used to generate this Value.
	formatOptions FormatOptions
}

// Plural returns the plural form of the prefix and unit.
//
// If unit singular is empty but has a plural suffix add the plural suffix
// to the prefix name, e.g., Prefix="hour", Unit={Singular:"", Plural:"s"} ->
// "hours".
//
// If there's no prefix name, just use the unit's plural name, e.g. "bytes",
// "hertz".
//
// Default: Combine prefix name and unit plural name, e.g. "kilo" + "bytes" ->
// "kilobytes".
//
// Examples: kilobytes, kibibytes, megahertz
func (u Value) Plural() string {
	if u.Unit.Singular == "" && u.Unit.Plural != "" {
		return u.Prefix.Name + u.Unit.Plural
	}

	if u.Prefix.Name == "" {
		return u.Unit.Plural
	}

	return u.Prefix.Name + u.Unit.Plural
}

// Singular returns the singular form of the prefix and unit.
//
// If unit singular is empty, the prefix name represents the whole unit, e.g.,
// Prefix="hour", Unit={Singular:"", Plural:"s"} -> "hour".
//
// If unit singular is empty and there's no prefix name, just use the unit's
// singular name, e.g. "byte", "hertz".
//
// Default: Combine prefix name and unit singular name, e.g. "kilo" + "byte" ->
// "kilobyte".
//
// Examples: kilobyte, kibibyte, megahertz
func (u Value) Singular() string {
	if u.Unit.Singular == "" {
		return u.Prefix.Name
	}

	if u.Prefix.Name == "" {
		return u.Unit.Singular
	}

	return u.Prefix.Name + u.Unit.Singular
}

// Symbol returns the symbol form of the prefix and unit.
//
// If unit symbol is empty, just use the prefix symbol, e.g., Prefix="h",
// Unit={Symbol:""} -> "h" .
//
// If there's no prefix symbol, just use the unit symbol, e.g. "B", "Hz".
//
// Default: Combine prefix symbol and unit symbol, e.g. "k" + "B" -> "kB".
//
// Examples: kB, KiB, MHz
func (u Value) Symbol() string {
	if u.Unit.Symbol == "" {
		return u.Prefix.Symbol
	}

	if u.Prefix.Symbol == "" {
		return u.Unit.Symbol
	}

	return u.Prefix.Symbol + u.Unit.Symbol
}

// String returns the formatted string representation of the unit value,
// fulfilling the fmt.Stringer interface.
//
// The output format is determined by the FormatOptions that were used when
// this Value struct was created (typically via the Format function). It
// combines the Scaled numeric value (formatted according to the Precision
// option) with a textual representation of the combined prefix and unit.
//
// Logic:
//  1. The Scaled value is rounded and formatted to a string using the
//     Precision specified in formatOptions.
//  2. A combined prefix/unit string is generated based on formatOptions:
//     - If UsePluralUnit is true, the Plural() method is called.
//     - If UsePluralUnit is false:
//     - If PrefixFormat is PrefixFormatShort, the Symbol() method is called.
//     - If PrefixFormat is PrefixFormatLong, the Singular() method is called.
//  3. If a non-empty prefix/unit string was generated:
//     - A space is inserted between the formatted number and the prefix/unit string if UseSpace is true.
//     - The formatted number, optional space, and prefix/unit string are concatenated.
//  4. If no prefix/unit string was generated (e.g., formatting a raw number
//     with no unit specified), only the formatted number string is returned.
//
// Examples (assuming appropriate Value was generated by Format):
//   - Value{Scaled: 1.5, Prefix: Kilo, Unit: Byte, opts: {Precision: 1, UsePluralUnit: true, UseSpace: true}} -> "1.5 kilobytes"
//   - Value{Scaled: 1.5, Prefix: Kilo, Unit: Byte, opts: {Precision: 0, PrefixFormat: PrefixFormatShort, UseSpace: true}} -> "2 kB"
//   - Value{Scaled: 1.5, Prefix: Kilo, Unit: Byte, opts: {Precision: 1, PrefixFormat: PrefixFormatLong, UseSpace: false}} -> "1.5kilobyte"
//   - Value{Scaled: 1024, Prefix: {}, Unit: {}, opts: {}} -> "1024"
func (u Value) String() string {
	formatted := roundAndFormatFloat(u.Scaled, u.formatOptions.Precision)

	var separator, prefixUnit string

	// Determine the combined prefix/unit string
	if u.formatOptions.UsePluralUnit {
		prefixUnit = u.Plural()
	} else {
		if u.formatOptions.PrefixFormat == PrefixFormatShort {
			prefixUnit = u.Symbol()
		} else {
			prefixUnit = u.Singular()
		}
	}

	// If no prefix/unit information, just return the number
	if prefixUnit == "" {
		return formatted
	}

	if u.formatOptions.UseSpace {
		separator = " "
	}

	return formatted + separator + prefixUnit
}

// Convert transforms a value from one unit scale to another by applying the
// appropriate conversion factor. This allows explicit conversion between
// different unit magnitudes.
//
// Parameters:
//   - value: The value to convert
//   - fromUnit: The unit scale of the input value (e.g., 1 for base unit, Kilo
//     for kilos, Mebi for mebibytes)
//   - toUnit: The unit scale to convert to (e.g., Mega for megabytes, Gibi for gibibytes)
//
// Examples:
//
//	Convert(1, Mega, Kilo)    // Convert 1 MB to KB = 1000 KB
//	Convert(1, Mebi, 1)       // Convert 1 MiB to bytes = 1048576 bytes
//	Convert(1024, 1, Kibi)    // Convert 1024 bytes to KiB = 1 KiB
//
// Converting between units
//
//	bytes := 1024 * 1024 // 1 MiB in bytes
//	megabytes := unit.Convert(bytes, 1, unit.Mebi) // Convert from bytes to MiB
//	fmt.Println(megabytes) // Output: 1
//
// Returns the converted value in the requested unit scale.
func Convert[T number](value T, fromUnit, toUnit float64) float64 {
	return float64(value) * fromUnit / toUnit
}

// roundAndFormatFloat converts a float64 to a string with the specified precision.
//
// It applies appropriate rounding and removes trailing zeros:
//   - For positive numbers, it rounds to the nearest value (0.5 rounds up)
//   - For negative numbers, it rounds to maintain the expected magnitude relationship
//   - Trailing zeros after the decimal point are removed
//   - If no decimal digits remain after removing zeros, the decimal point is also removed
//
// Examples:
//
//	roundAndFormatFloat(1.5, 0)    => "2"
//	roundAndFormatFloat(1.5, 1)    => "1.5"
//	roundAndFormatFloat(1.50, 2)   => "1.5"
//	roundAndFormatFloat(-1.5, 0)   => "-2"
//	roundAndFormatFloat(1.0, 2)    => "1"
func roundAndFormatFloat(num float64, precision int) string {
	// Ensure precision is non-negative
	if precision < 0 {
		precision = 0
	}

	// Calculate scale factor for rounding
	scale := math.Pow10(precision)

	// For formatting purposes, it's reasonable to want `-1.5` to round to `-2`
	// rather than `-1`. This preserves the magnitude relationship, which is
	// what users typically expect when viewing formatted units.
	var roundedNum float64
	if num < 0 {
		// For negative numbers, round toward negative infinity for ties
		roundedNum = float64(int(num*scale-0.5)) / scale
	} else {
		// For positive numbers, round toward positive infinity for ties
		roundedNum = float64(int(num*scale+0.5)) / scale
	}

	formattedBytes := strconv.AppendFloat(nil, roundedNum, 'f', precision, float64BitSize)

	// Trim trailing zeros and the decimal point
	trimmedBytes := trimFloatBytes(formattedBytes, precision)

	return string(trimmedBytes)
}

// trimFloatBytes removes trailing zeros and the decimal point (if necessary)
// from a byte slice representing a float formatted with 'f'.
// It assumes the input slice was generated by strconv.AppendFloat using 'f'
// format with the given precision.
//
// It is an optimized alternative to strings.TrimRight()
func trimFloatBytes(b []byte, precision int) []byte {
	if precision <= 0 {
		// If precision is 0, strconv might have added a trailing ".", remove
		// it.
		if len(b) > 0 && b[len(b)-1] == '.' {
			return b[:len(b)-1]
		}
		// Otherwise, no fractional part to trim.
		return b
	}

	// Precision > 0: Find the decimal point.
	// strconv.AppendFloat with 'f' and precision > 0 guarantees a decimal
	// point.
	dotIndex := -1
	for i, char := range b {
		if char == '.' {
			dotIndex = i
			break
		}
	}

	// If no dot found (unexpected for precision > 0), return original.
	if dotIndex == -1 {
		return b
	}

	// Find the index of the last non-zero digit after the decimal point.
	trimIndex := len(b) - 1
	for trimIndex > dotIndex && b[trimIndex] == '0' {
		trimIndex--
	}

	// If the last non-zero digit is the decimal point itself, remove the point too.
	if trimIndex == dotIndex {
		return b[:dotIndex]
	}

	// Otherwise, trim up to the last non-zero digit.
	return b[:trimIndex+1]
}
