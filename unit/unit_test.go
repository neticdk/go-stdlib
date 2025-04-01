package unit

import (
	"slices"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		name     string
		value    float64
		from     float64
		to       float64
		expected float64
	}{
		{
			name:     "bytes to kilobytes",
			value:    1024,
			from:     1,
			to:       Kilo,
			expected: 1.024,
		},
		{
			name:     "kilobytes to bytes",
			value:    1,
			from:     Kilo,
			to:       1,
			expected: 1000,
		},
		{
			name:     "megabytes to kilobytes",
			value:    2,
			from:     Mega,
			to:       Kilo,
			expected: 2000,
		},
		{
			name:     "kibibytes to bytes",
			value:    1,
			from:     Kibi,
			to:       1,
			expected: 1024,
		},
		{
			name:     "bytes to mebibytes",
			value:    2097152,
			from:     1,
			to:       Mebi,
			expected: 2,
		},
		{
			name:     "gigabytes to bytes",
			value:    1,
			from:     Giga,
			to:       1,
			expected: 1e9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Convert(tc.value, tc.from, tc.to)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInfo(t *testing.T) {
	testCases := []struct {
		name     string
		unitKind Unit
		expected Descriptor
	}{
		{
			name:     "byte unit",
			unitKind: Byte,
			expected: Descriptor{Byte, "B", "byte", "bytes"},
		},
		{
			name:     "bit unit",
			unitKind: Bit,
			expected: Descriptor{Bit, "b", "bit", "bits"},
		},
		{
			name:     "hertz unit",
			unitKind: Hertz,
			expected: Descriptor{Hertz, "Hz", "hertz", "hertz"},
		},
		{
			name:     "none/unknown unit",
			unitKind: None,
			expected: Descriptor{},
		},
		{
			name:     "invalid unit type",
			unitKind: Unit(999),
			expected: Descriptor{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Describe(tc.unitKind)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValueMethods(t *testing.T) {
	testCases := []struct {
		name             string
		value            Value
		expectedPlural   string
		expectedSingular string
		expectedSymbol   string
	}{
		{
			name: "kilobyte",
			value: Value{
				Scaled: 1.5,
				Prefix: prefixes[Kilo],
				Unit:   unitRegistry[Byte],
			},
			expectedPlural:   "kilobytes",
			expectedSingular: "kilobyte",
			expectedSymbol:   "kB",
		},
		{
			name: "megabit",
			value: Value{
				Scaled: 2.0,
				Prefix: prefixes[Mega],
				Unit:   unitRegistry[Bit],
			},
			expectedPlural:   "megabits",
			expectedSingular: "megabit",
			expectedSymbol:   "Mb",
		},
		{
			name: "gigahertz",
			value: Value{
				Scaled: 3.5,
				Prefix: prefixes[Giga],
				Unit:   unitRegistry[Hertz],
			},
			expectedPlural:   "gigahertz",
			expectedSingular: "gigahertz",
			expectedSymbol:   "GHz",
		},
		{
			name: "mebibyte",
			value: Value{
				Scaled: 1.0,
				Prefix: prefixes[Mebi],
				Unit:   unitRegistry[Byte],
			},
			expectedPlural:   "mebibytes",
			expectedSingular: "mebibyte",
			expectedSymbol:   "MiB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPlural, tc.value.Plural())
			assert.Equal(t, tc.expectedSingular, tc.value.Singular())
			assert.Equal(t, tc.expectedSymbol, tc.value.Symbol())
		})
	}
}

func TestValueFormat(t *testing.T) {
	testCases := []struct {
		name           string
		value          Value
		formatOptions  FormatOptions
		expectedFormat string
	}{
		{
			name: "default format",
			value: Value{
				Scaled: 1.5,
				Prefix: prefixes[Kilo],
				Unit:   unitRegistry[Byte],
				formatOptions: FormatOptions{
					Precision:     0,
					UseSpace:      true,
					UsePluralUnit: true,
					PrefixFormat:  PrefixFormatShort,
				},
			},
			expectedFormat: "2 kilobytes",
		},
		{
			name: "with precision",
			value: Value{
				Scaled: 1.5,
				Prefix: prefixes[Kilo],
				Unit:   unitRegistry[Byte],
				formatOptions: FormatOptions{
					Precision:     2,
					UseSpace:      true,
					UsePluralUnit: true,
					PrefixFormat:  PrefixFormatShort,
				},
			},
			expectedFormat: "1.5 kilobytes",
		},
		{
			name: "no space",
			value: Value{
				Scaled: 1.5,
				Prefix: prefixes[Kilo],
				Unit:   unitRegistry[Byte],
				formatOptions: FormatOptions{
					Precision:     0,
					UseSpace:      false,
					UsePluralUnit: true,
					PrefixFormat:  PrefixFormatShort,
				},
			},
			expectedFormat: "2kilobytes",
		},
		{
			name: "symbolic format",
			value: Value{
				Scaled: 1.5,
				Prefix: prefixes[Kilo],
				Unit:   unitRegistry[Byte],
				formatOptions: FormatOptions{
					Precision:     0,
					UseSpace:      true,
					UsePluralUnit: false,
					PrefixFormat:  PrefixFormatShort,
				},
			},
			expectedFormat: "2 kB",
		},
		{
			name: "singular format",
			value: Value{
				Scaled: 1.5,
				Prefix: prefixes[Kilo],
				Unit:   unitRegistry[Byte],
				formatOptions: FormatOptions{
					Precision:     0,
					UseSpace:      true,
					UsePluralUnit: false,
					PrefixFormat:  PrefixFormatLong,
				},
			},
			expectedFormat: "2 kilobyte",
		},
		{
			name: "zero value",
			value: Value{
				Scaled: 0,
				Prefix: Prefix{},
				Unit:   Descriptor{},
				formatOptions: FormatOptions{
					Precision:     0,
					UseSpace:      true,
					UsePluralUnit: true,
				},
			},
			expectedFormat: "0",
		},
		{
			name: "large negative value",
			value: Value{
				Scaled: -2.5,
				Prefix: prefixes[Giga],
				Unit:   unitRegistry[Byte],
				formatOptions: FormatOptions{
					Precision:     1,
					UseSpace:      true,
					UsePluralUnit: true,
				},
			},
			expectedFormat: "-2.5 gigabytes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.value.String()
			assert.Equal(t, tc.expectedFormat, result)
			// Test String() which should call String()
			assert.Equal(t, tc.expectedFormat, tc.value.String())
		})
	}
}

func TestBoundariesForFormat(t *testing.T) {
	testCases := []struct {
		name     string
		format   FormatStyle
		expected []float64
	}{
		{
			name:     "decimal format",
			format:   DecimalFormat,
			expected: SIDecimalBoundaries,
		},
		{
			name:     "decimal bytes format",
			format:   DecimalBytesFormat,
			expected: SIDecimalBoundaries,
		},
		{
			name:     "decimal bits format",
			format:   DecimalBitsFormat,
			expected: SIDecimalBoundaries,
		},
		{
			name:     "decimal hertz format",
			format:   DecimalHertzFormat,
			expected: SIDecimalBoundaries,
		},
		{
			name:     "binary format",
			format:   BinaryFormat,
			expected: IECBinaryBoundaries,
		},
		{
			name:     "binary bytes format",
			format:   BinaryBytesFormat,
			expected: IECBinaryBoundaries,
		},
		{
			name:     "binary bits format",
			format:   BinaryBitsFormat,
			expected: IECBinaryBoundaries,
		},
		{
			name:     "unknown format",
			format:   FormatStyle("unknown"),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := boundariesForFormat(tc.format)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFromFormat(t *testing.T) {
	testCases := []struct {
		name     string
		format   FormatStyle
		expected Unit
	}{
		{
			name:     "decimal bytes format",
			format:   DecimalBytesFormat,
			expected: Byte,
		},
		{
			name:     "binary bytes format",
			format:   BinaryBytesFormat,
			expected: Byte,
		},
		{
			name:     "decimal bits format",
			format:   DecimalBitsFormat,
			expected: Bit,
		},
		{
			name:     "binary bits format",
			format:   BinaryBitsFormat,
			expected: Bit,
		},
		{
			name:     "decimal hertz format",
			format:   DecimalHertzFormat,
			expected: Hertz,
		},
		{
			name:     "decimal format",
			format:   DecimalFormat,
			expected: None,
		},
		{
			name:     "binary format",
			format:   BinaryFormat,
			expected: None,
		},
		{
			name:     "unknown format",
			format:   FormatStyle("unknown"),
			expected: None,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fromFormat(tc.format)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFormatFloat(t *testing.T) {
	testCases := []struct {
		name      string
		value     float64
		precision int
		expected  string
	}{
		{
			name:      "integer",
			value:     42,
			precision: 0,
			expected:  "42",
		},
		{
			name:      "decimal with trailing zeros",
			value:     1.500,
			precision: 3,
			expected:  "1.5",
		},
		{
			name:      "decimal with multiple trailing zeros",
			value:     1.5000,
			precision: 4,
			expected:  "1.5",
		},
		{
			name:      "exact precision",
			value:     1.234,
			precision: 3,
			expected:  "1.234",
		},
		{
			name:      "rounded decimal",
			value:     1.2345,
			precision: 3,
			expected:  "1.235",
		},
		{
			name:      "negative value",
			value:     -1.5,
			precision: 1,
			expected:  "-1.5",
		},
		{
			name:      "zero",
			value:     0,
			precision: 2,
			expected:  "0",
		},
		{
			name:      "exact integer after rounding",
			value:     1.999,
			precision: 2,
			expected:  "2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := roundAndFormatFloat(tc.value, tc.precision)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNames(t *testing.T) {
	// Get all names - skip other custom units added in the test
	allNames := Names()
	names := make([]string, 0, 3)
	for _, name := range allNames {
		if slices.Contains([]string{"byte", "bit", "hertz"}, name) {
			names = append(names, name)
			continue
		}
	}

	// Expected unit names
	expected := []string{
		"byte",
		"bit",
		"hertz",
	}

	// Assert that all expected names are in the result
	// (not checking order since map iteration isn't guaranteed)
	assert.ElementsMatch(t, expected, names)
	assert.Len(t, names, 3, "Should have exactly 3 unit names")
}

func TestList(t *testing.T) {
	// Get all units - skip other custom units added in the test
	allUnits := List()
	units := make([]Descriptor, 0, 3)
	for _, unit := range allUnits {
		if slices.Contains([]Unit{Byte, Bit, Hertz}, unit.Unit) {
			units = append(units, unit)
			continue
		}
	}

	// Create a map of expected unit descriptors
	expectedUnits := []Descriptor{
		{Unit: Byte, Symbol: "B", Singular: "byte", Plural: "bytes"},
		{Unit: Bit, Symbol: "b", Singular: "bit", Plural: "bits"},
		{Unit: Hertz, Symbol: "Hz", Singular: "hertz", Plural: "hertz"},
	}

	// Test that all expected units are in the result
	// (not checking order since map iteration isn't guaranteed)
	assert.ElementsMatch(t, expectedUnits, units)
	assert.Len(t, units, 3, "Should have exactly 3 units")
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		name     string
		unit     Unit
		expected Descriptor
	}{
		{
			name:     "Byte unit",
			unit:     Byte,
			expected: Descriptor{Unit: Byte, Symbol: "B", Singular: "byte", Plural: "bytes"},
		},
		{
			name:     "Bit unit",
			unit:     Bit,
			expected: Descriptor{Unit: Bit, Symbol: "b", Singular: "bit", Plural: "bits"},
		},
		{
			name:     "Hertz unit",
			unit:     Hertz,
			expected: Descriptor{Unit: Hertz, Symbol: "Hz", Singular: "hertz", Plural: "hertz"},
		},
		{
			name:     "None unit",
			unit:     None,
			expected: Descriptor{}, // Empty descriptor for None
		},
		{
			name:     "Unregistered unit",
			unit:     Unit(999),
			expected: Descriptor{}, // Empty descriptor for unknown units
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Describe(tt.unit)
			assert.Equal(t, tt.expected, result)
		})
	}
}
