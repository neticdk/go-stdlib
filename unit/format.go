package unit

// PrefixFormat represents the format to use for unit prefixes
type PrefixFormat string

const (
	PrefixFormatShort PrefixFormat = "short" // k, M, G, T, Ki, Mi, Gi, Ti, etc
	PrefixFormatLong  PrefixFormat = "long"  // kilo, mega, giga, tera, kibi, mebi, gibi, tebi, etc
)

// FormatStyle represents the format to use for units
type FormatStyle string

const (
	DecimalFormat      FormatStyle = "decimal"      // k, kilo, M, mega, G, giga, T, tera, etc
	DecimalBytesFormat FormatStyle = "decimalBytes" // kB, kilobytes, MB, megabytes, GB, gigabytes, TB, terabytes, etc
	DecimalBitsFormat  FormatStyle = "decimalBits"  // b, bits, Kb, kilobits, Mb, megabits, Gb, gigabits, Tb, terabits, etc
	DecimalHertzFormat FormatStyle = "decimalHertz" // Hz, hertz, kHz, kilohertz, MHz, megahertz, GHz, gigahertz, THz, terahertz, etc
	BinaryFormat       FormatStyle = "binary"       // Ki, kibi, Mi, mebi, Gi, gibi, Ti, tebi, etc
	BinaryBytesFormat  FormatStyle = "binaryBytes"  // B, bytes, KiB, kibibytes, MiB, mebibytes, GiB, gibibytes, TiB, tebibytes, etc
	BinaryBitsFormat   FormatStyle = "binaryBits"   // b, bits, Kib, kibibits, Mib, mebibits, Gib, gibibits, Tib, tebibits, etc
)

// FormatOptions customizes the output of the formatting functions
// Note that UsePluralUnit implies using a long prefix format
type FormatOptions struct {
	Precision     int          // Number of decimal places (default: 0)
	UseSpace      bool         // Add space between value and unit (default: true)
	UsePluralUnit bool         // Use plural form of units (default: true)
	PrefixFormat  PrefixFormat // Format to use for unit prefixes (default: "short")
	FormatStyle   FormatStyle  // Format style to use for units (default: "decimal")
	Unit          Unit         // The unit to format with

	// Support for custom format systems
	Boundaries []float64          // Custom boundaries for unit scaling
	PrefixMap  map[float64]Prefix // Custom prefixes for unit formatting
	SystemName string             // Name of the formatting system
}

// DefaultFormatOptions provides default formatting settings
var DefaultFormatOptions = FormatOptions{
	Precision:     0,
	UseSpace:      true,
	UsePluralUnit: false,
	PrefixFormat:  PrefixFormatShort,
	FormatStyle:   DecimalFormat,
}

// FormatOption represents a functional option for the Format function.
type FormatOption func(*FormatOptions)

// Format applies formatting options to a numeric value and returns the
// value. It handles both binary and decimal formatting with various
// unit types.
//
// Examples:
//
//	// Basic formatting
//	Format(1024)                           // "1024"
//	Format(1024, Binary())                 // "1 Ki"
//	Format(1024, Binary(), WithUnit(Byte)) // "1 KiB"
//
//	// Combining multiple options
//	Format(1024, Binary(), WithUnit(Byte), WithPrecision(2))  // "1.00 KiB"
//	Format(1500, Decimal(), WithUnit(Bit), WithPlural())      // "1.5 kilobits"
//
//	// Advanced formatting
//	Format(1024, Binary(), WithUnit(Byte), WithoutSpace())    // "1KiB"
//	Format(1024, Binary(), WithUnit(Byte), WithLongPrefix())  // "1 kibibyte"
//
// The first option that specifies a unit system (Binary/Decimal) takes
// precedence. If no unit system is specified, the default is Decimal.
func Format[T number](value T, opts ...FormatOption) Value {
	options := DefaultFormatOptions
	for _, opt := range opts {
		opt(&options)
	}

	// Validate and normalize options
	validateOptions(&options)

	// Determine unit:
	// 1. Use the explicitly set unit from options if available
	// 2. Otherwise derive it from the format
	unit := options.Unit
	if unit == None {
		unit = fromFormat(options.FormatStyle)
	}

	return formatValueWithUnit(value, unit, options)
}

// Binary configures the Format function to use IEC binary units (base-2).
// Units increase by powers of 1024 with prefixes like Ki, Mi, Gi.
func Binary() FormatOption {
	return func(o *FormatOptions) {
		// Only set if not already set to a binary format
		if o.FormatStyle != BinaryFormat &&
			o.FormatStyle != BinaryBytesFormat &&
			o.FormatStyle != BinaryBitsFormat {
			o.FormatStyle = BinaryFormat
		}

		// Set binary boundaries and prefixes if no custom system is defined
		if len(o.Boundaries) == 0 {
			o.Boundaries = IECBinaryBoundaries
		}
		if o.PrefixMap == nil {
			o.PrefixMap = binaryPrefixes
		}
	}
}

// Decimal configures the Format function to use SI decimal units (base-10).
// Units increase by powers of 1000 with prefixes like k, M, G.
func Decimal() FormatOption {
	return func(o *FormatOptions) {
		// Only set if not already set to a decimal format
		if o.FormatStyle != DecimalFormat &&
			o.FormatStyle != DecimalBytesFormat &&
			o.FormatStyle != DecimalBitsFormat &&
			o.FormatStyle != DecimalHertzFormat {
			o.FormatStyle = DecimalFormat
		}

		// Set decimal boundaries and prefixes if no custom system is defined
		if len(o.Boundaries) == 0 {
			o.Boundaries = SIDecimalBoundaries
		}
		if o.PrefixMap == nil {
			o.PrefixMap = decimalPrefixes
		}
	}
}

// WithUnit specifies the unit type to use (Byte, Bit, Hz).
// This overrides any previous unit type settings. For built-in units, we also
// set appropriate format potentially overriding any previous format settings.
// For custom units, we don't change the format.
//
// Examples:
//
//	Format(1024, Binary(), WithUnit(Byte))  // "1 KiB"
//	Format(1024, Binary(), WithUnit(Bit))   // "1 Kib"
//	Format(1000, Decimal(), WithUnit(Hz))   // "1 kHz"
func WithUnit(unit Unit) FormatOption {
	return func(o *FormatOptions) {
		// Store the unit directly in the options
		o.Unit = unit

		switch unit {
		case Byte:
			if isFormatBinary(o.FormatStyle) {
				o.FormatStyle = BinaryBytesFormat
			} else {
				o.FormatStyle = DecimalBytesFormat
			}
		case Bit:
			if isFormatBinary(o.FormatStyle) {
				o.FormatStyle = BinaryBitsFormat
			} else {
				o.FormatStyle = DecimalBitsFormat
			}
		case Hertz:
			o.FormatStyle = DecimalHertzFormat
		}
	}
}

// WithPrecision sets the number of decimal places to display.
//
// Setting a negative precision will result in no decimal places.
//
// Examples:
//
//	Format(1.5, WithPrecision(2))  // "1.50"
//	Format(1.5, WithPrecision(0))  // "2"
func WithPrecision(precision int) FormatOption {
	return func(o *FormatOptions) {
		if precision < 0 {
			precision = 0
		}
		o.Precision = precision
	}
}

// WithPlural formats units using their plural names.
// This automatically uses long prefix formats.
//
// Examples:
//
//	Format(1024, Binary(), WithUnit(Byte), WithPlural())  // "1 kibibytes"
//	Format(1000, Decimal(), WithUnit(Bit), WithPlural())  // "1 kilobits"
func WithPlural() FormatOption {
	return func(o *FormatOptions) {
		o.UsePluralUnit = true
		o.PrefixFormat = PrefixFormatLong
	}
}

// WithLongPrefix uses the full name of unit prefixes instead of symbols.
//
// Examples:
//
//	Format(1024, Binary(), WithLongPrefix())  // "1 kibi"
//	Format(1000, Decimal(), WithLongPrefix()) // "1 kilo"
func WithLongPrefix() FormatOption {
	return func(o *FormatOptions) {
		o.PrefixFormat = PrefixFormatLong
	}
}

// WithoutSpace removes the space between the value and unit.
//
// Examples:
//
//	Format(1024, Binary(), WithUnit(Byte), WithoutSpace())  // "1KiB"
//	Format(1000, Decimal(), WithUnit(Bit), WithoutSpace())  // "1kb"
func WithoutSpace() FormatOption {
	return func(o *FormatOptions) {
		o.UseSpace = false
	}
}

// WithSystem allows using a custom unit system for formatting
func WithSystem(system *FormatSystem) FormatOption {
	return func(o *FormatOptions) {
		if system != nil {
			o.Boundaries = system.Boundaries
			o.PrefixMap = system.Prefixes
			o.SystemName = system.Name
		}
	}
}

// WithSystemByName allows using a registered custom unit system by name
func WithSystemByName(name string) FormatOption {
	return func(o *FormatOptions) {
		if system, found := GetFormatSystem(name); found {
			o.Boundaries = system.Boundaries
			o.PrefixMap = system.Prefixes
			o.SystemName = name
		}
	}
}

// FormatBinary converts a numeric value to a human-readable string using IEC
// 80000-13 binary (base2) units. It picks the unit that converts to the
// smallest value and returns value and shorthand unit prefix. Examples of
// values: bytes, bits, hertz, etc.
//
// Example:
//
//	FormatBinary(42)   // Value{Scaled: 42.0, Prefix: Prefix{}} -> String(): "42"
//	FormatBinary(1000) // Value{Scaled: 1000.0, Prefix: Prefix{}} -> String(): "1000"
//	FormatBinary(1024) // Value{Scaled: 1.0, Prefix: Prefix{Name: "kibi", Symbol: "Ki"}} -> String(): "1 Ki"
func FormatBinary[T number](value T) Value {
	return Format(value, Binary())
}

// FormatBinaryUnit converts a numeric value with an associated unit to
// a human-readable string using IEC 80000-13 binary (base2) units. It picks
// the unit that converts to the smallest value and returns the value and
// shorthand unit prefix.
//
// Examples:
//
//	FormatBinaryUnit(42, Byte)    // Value{Scaled: 42.0, Prefix: Prefix{}, Unit: unit.Describe(Byte)} -> String(): "42 B"
//	FormatBinaryUnit(1000, Bit)   // Value{Scaled: 1000.0, Prefix: Prefix{}, Unit: unit.Describe(Bit)} -> String(): "1000 b"
//	FormatBinaryUnit(1024, Byte)  // Value{Scaled: 1.0, Prefix: unit.PrefixFor(Kibi), Unit: unit.Describe(Byte)} -> String(): "1 KiB"
//
// Displaying file sizes in a human-readable way
//
//	fileSize := 1572864 // bytes
//	value := unit.FormatBinaryUnit(fileSize, unit.Byte)
//	fmt.Printf("File size: %s\n", value.String()) // Output: File size: 1.5 MiB
func FormatBinaryUnit[T number](value T, unit Unit) Value {
	return Format(value, Binary(), WithUnit(unit))
}

// FormatDecimal converts a numeric value to human-readable string using ISO
// 1000 metric (base10) SI units. It picks the unit that converts to the
// smallest value and returns the value and shorthand unit symbol.
// Examples of values: bytes, bits, hertz, etc.
//
// Example:
//
//	FormatDecimal(42)      // Value{Scaled: 42.0, Prefix: Prefix{}} -> String(): "42"
//	FormatDecimal(1000)    // Value{Scaled: 1.0, Prefix: Prefix{Name: "kilo", Symbol: "k"}} -> String(): "1 k"
//	FormatDecimal(1024)    // Value{Scaled: 1.024, Prefix: Prefix{Name: "kilo", Symbol: "k"}} -> String(): "1 k" (due to default precision 0)
func FormatDecimal[T number](value T) Value {
	return Format(value, Decimal())
}

// FormatDecimalUnit converts a numeric value with an associated unit to
// a human-readable string using ISO 1000 metric (base10) SI units. It picks
// the unit that converts to the smallest value and returns the value and
// shorthand unit prefix.
//
// Example:
//
//	FormatDecimalUnit(42, Byte)    // Value{Scaled: 42.0, Prefix: Prefix{}, Unit: unit.Describe(Byte)} -> String(): "42 B"
//	FormatDecimalUnit(1000, Bit)   // Value{Scaled: 1.0, Prefix: unit.PrefixFor(Kilo), Unit: unit.Describe(Bit)} -> String(): "1 kb"
//	FormatDecimalUnit(1024, Byte)  // Value{Scaled: 1.024, Prefix: unit.PrefixFor(Kilo), Unit: unit.Describe(Byte)} -> String(): "1 kB" (due to default precision 0)
//
// Network bandwidth display
//
//	bandwidth := 1500000 // bits per second
//	value := unit.FormatDecimalUnit(bandwidth, unit.Bit)
//	fmt.Printf("Network speed: %s/s\n", value.String()) // Output: Network speed: 1.5 Mb/s
func FormatDecimalUnit[T number](value T, unit Unit) Value {
	return Format(value, Decimal(), WithUnit(unit))
}

// formatValueWithUnit formats a numeric value with the appropriate unit prefix
// based on its magnitude.
// It scales the value to the largest unit that keeps the number readable
// (typically between 1-999).
//
// Parameters:
//   - value: The numeric value to format
//   - unit: The base unit type (Byte, Bit, Hz, or None)
//   - boundaries: The unit boundaries to use for scaling (typically
//     SIDecimalBoundaries or IECBinaryBoundaries)
//   - opts: Formatting options to control precision, spacing, and unit display
//
// The function handles special cases including:
//   - Zero values (returns "0" or "0 [unit]" depending on unit type)
//   - Negative values (preserves sign and applies proper rounding)
//   - Values smaller than the smallest boundary (displays with no prefix)
//   - Values larger than the largest boundary (uses the largest available prefix)
//
// The return value contains both the scaled numeric value and the appropriate
// prefix information for further formatting.
func formatValueWithUnit[T number](value T, unit Unit, opts FormatOptions) Value {
	v := float64(value)

	returnedValue := Value{
		Unit:          Describe(unit),
		formatOptions: opts,
	}

	// Handle zero value explicitly
	if v == 0 {
		returnedValue.Scaled = 0
		returnedValue.Prefix = Prefix{}
		return returnedValue
	}

	// Determine which boundaries to use (custom or standard)
	boundaries := opts.Boundaries
	if len(boundaries) == 0 {
		// Fall back to standard boundaries based on format
		boundaries = boundariesForFormat(opts.FormatStyle)
		if len(boundaries) == 0 {
			// Default to decimal if no valid boundaries found
			boundaries = SIDecimalBoundaries
		}
	}

	// Determine which prefix map to use (custom or standard)
	prefixMap := opts.PrefixMap
	if len(prefixMap) == 0 {
		// Use global prefix map
		prefixMap = prefixes
	}

	// Handle negative values by using absolute value for unit calculation
	// then restoring the sign later
	isNegative := v < 0
	if isNegative {
		v = -v
	}

	// Handle different value ranges
	switch {
	// Handle large values - if larger than the largest boundary
	case v >= boundaries[0]:
		returnedValue.Scaled = v / boundaries[0]
		returnedValue.Prefix = prefixFromMap(boundaries[0], prefixMap)

	// Handle small values - if smaller than the smallest boundary
	case v < boundaries[len(boundaries)-1]:
		returnedValue.Scaled = v
		returnedValue.Prefix = Prefix{}

	// Handle values in the normal range
	default:
		for _, b := range boundaries {
			if v >= b {
				returnedValue.Scaled = v / b
				returnedValue.Prefix = prefixFromMap(b, prefixMap)
				break
			}
		}
	}

	// If no match was found, use the original value
	if returnedValue.Scaled == 0 {
		returnedValue.Scaled = v
		returnedValue.Prefix = Prefix{}
	}

	// Restore negative sign if needed
	if isNegative {
		returnedValue.Scaled = -returnedValue.Scaled
	}

	return returnedValue
}

// boundariesForFormat returns the boundaries for a given format.
func boundariesForFormat(style FormatStyle) []float64 {
	switch style {
	case DecimalFormat, DecimalBytesFormat, DecimalBitsFormat, DecimalHertzFormat:
		return SIDecimalBoundaries
	case BinaryFormat, BinaryBytesFormat, BinaryBitsFormat:
		return IECBinaryBoundaries
	default:
		return nil
	}
}

// fromFormat returns the unit for a given format.
func fromFormat(style FormatStyle) Unit {
	switch style {
	case DecimalBytesFormat, BinaryBytesFormat:
		return Byte
	case DecimalBitsFormat, BinaryBitsFormat:
		return Bit
	case DecimalHertzFormat:
		return Hertz
	default:
		return None
	}
}

// validateOptions performs validation and default value setting on format
// options. It ensures that all necessary options have appropriate values.
func validateOptions(opts *FormatOptions) {
	if opts.Precision < 0 {
		opts.Precision = 0
	}

	// If user specified plural units, automatically use long prefix format
	if opts.UsePluralUnit && opts.PrefixFormat == PrefixFormatShort {
		opts.PrefixFormat = PrefixFormatLong
	}

	// Ensure we have a valid format
	if opts.FormatStyle == "" {
		opts.FormatStyle = DecimalFormat
	}
}

// Helper function to check if a format is binary
func isFormatBinary(style FormatStyle) bool {
	return style == BinaryFormat ||
		style == BinaryBytesFormat ||
		style == BinaryBitsFormat
}
