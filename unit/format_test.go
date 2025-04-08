package unit

import (
	"fmt"
	"slices"
	"sync"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

// Test custom unit registration and lookup
func TestCustomUnits(t *testing.T) {
	// Create a custom unit
	UnitPixel, err := Register(Descriptor{
		Symbol:   "px",
		Singular: "pixel",
		Plural:   "pixels",
	})

	assert.NoError(t, err)

	// Test that we can retrieve the unit descriptor
	desc := Describe(UnitPixel)
	assert.Equal(t, "px", desc.Symbol)
	assert.Equal(t, "pixel", desc.Singular)
	assert.Equal(t, "pixels", desc.Plural)

	// Check unit discovery
	allUnits := Names()
	found := slices.Contains(allUnits, "pixel")
	assert.True(t, found, "Custom unit 'pixel' should be discoverable in Names()")

	// Test formatting with the custom unit
	value := FormatDecimalUnit(1234, UnitPixel)
	assert.Equal(t, 1.234, value.Scaled)
	assert.Equal(t, prefixes[Kilo], value.Prefix)
	assert.Equal(t, "1 kpx", value.String())
}

// Test custom format system
func TestCustomFormatSystem(t *testing.T) {
	// Register a custom time format system
	timeSystem, err := RegisterFormatSystem(
		"time",
		[]float64{3600, 60, 1}, // hours, minutes, seconds
		map[float64]Prefix{
			3600: {"hour", "h"},
			60:   {"minute", "m"},
			1:    {"second", "s"},
		},
	)

	assert.NoError(t, err)

	// Register a time unit
	UnitTime, err := Register(Descriptor{
		Symbol:   "",
		Singular: "",
		Plural:   "s",
	})
	assert.NoError(t, err)

	// Format a duration
	duration := 3725 // 1 hour, 2 minutes, 5 seconds
	value := Format(duration,
		WithSystem(timeSystem),
		WithUnit(UnitTime))

	// Use approximate comparison for floating point
	assert.InDelta(t, 1.035, value.Scaled, 0.001)

	assert.Equal(t, "h", value.Prefix.Symbol)
	assert.Equal(t, "hour", value.Prefix.Name)
	assert.Equal(t, "1 h", value.String())

	// Test with precision
	preciseFormat := Format(duration,
		WithSystem(timeSystem),
		WithUnit(UnitTime),
		WithPrecision(2),
		WithPlural())

	assert.InDelta(t, 1.035, preciseFormat.Scaled, 0.001)
	assert.True(t, preciseFormat.formatOptions.UsePluralUnit,
		"UsePluralUnit should be true when WithPlural() is used")
	assert.Equal(t, PrefixFormatLong, preciseFormat.formatOptions.PrefixFormat,
		"PrefixFormat should be Long when WithPlural() is used")

	// With plural and long prefix, it should show full names
	assert.Equal(t, "1.03 hours", preciseFormat.String())
}

func TestSystemRegistryThreadSafety(t *testing.T) {
	// Test concurrent registration of systems
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("system-%d", i)
			RegisterFormatSystem(
				name,
				[]float64{float64(i) * 10, float64(i)},
				map[float64]Prefix{
					float64(i) * 10: {Symbol: "X", Name: "deca"},
					float64(i):      {Symbol: "Y", Name: "unit"},
				},
			)
		}(i)
	}
	wg.Wait()

	// Verify some of the registered systems
	system, found := GetFormatSystem("system-42")
	assert.True(t, found)
	assert.Equal(t, "system-42", system.Name)
	assert.Equal(t, 2, len(system.Boundaries))
	assert.Equal(t, float64(420), system.Boundaries[0])
}

func TestUnitRegistry(t *testing.T) {
	// Test basic registry functions
	units := Names()
	assert.Contains(t, units, "byte")
	assert.Contains(t, units, "bit")

	// Get information about units
	unitDescriptors := List()
	var byteDescriptor Descriptor
	for _, info := range unitDescriptors {
		if info.Singular == "byte" {
			byteDescriptor = info
			break
		}
	}
	assert.Equal(t, Byte, byteDescriptor.Unit)
	assert.Equal(t, "B", byteDescriptor.Symbol)
	assert.Equal(t, "bytes", byteDescriptor.Plural)

	// Test Unit.Info() method
	byteUnitInfo := Byte.Info()
	assert.Equal(t, "B", byteUnitInfo.Symbol)
	assert.Equal(t, "byte", byteUnitInfo.Singular)
}

// Update the existing formatValueWithUnit test to work with updated function signature
func TestFormatValueWithUnit(t *testing.T) {
	pluralOpts := DefaultFormatOptions
	pluralOpts.UsePluralUnit = true

	singularOpts := DefaultFormatOptions
	singularOpts.PrefixFormat = PrefixFormatLong

	precisionOptions := DefaultFormatOptions
	precisionOptions.Precision = 2

	testCases := []struct {
		name                string
		input               float64
		unit                Unit
		options             FormatOptions
		expectedValue       float64
		expectedPrefix      Prefix
		expectedUnit        Descriptor
		expectedValueString string
		description         string
	}{
		// Zero value
		{
			name:                "zero_value",
			input:               0,
			unit:                None,
			options:             DefaultFormatOptions,
			expectedValue:       0,
			expectedPrefix:      Prefix{},
			expectedValueString: "0",
			description:         "Zero should return zero with no prefix",
		},

		// Negative values
		{
			name:                "negative_small",
			input:               -42,
			unit:                None,
			options:             DefaultFormatOptions,
			expectedValue:       -42,
			expectedPrefix:      Prefix{},
			expectedValueString: "-42",
			description:         "Small negative value should preserve sign with no prefix",
		},
		{
			name:                "negative_kilo",
			input:               -1500,
			unit:                None,
			options:             DefaultFormatOptions,
			expectedValue:       -1.5,
			expectedPrefix:      prefixes[Kilo],
			expectedValueString: "-2 k",
			description:         "Negative kilos should round and preserve sign with kilo prefix",
		},

		// Test with custom format system
		{
			name:  "custom_system",
			input: 90,
			unit:  None,
			options: FormatOptions{
				Precision:     0,
				UseSpace:      true,
				UsePluralUnit: false,
				PrefixFormat:  PrefixFormatShort,
				Boundaries:    []float64{100, 10, 1},
				PrefixMap: map[float64]Prefix{
					100: {"hundred", "h"},
					10:  {"ten", "t"},
					1:   {"", ""},
				},
			},
			expectedValue:       9,
			expectedPrefix:      Prefix{"ten", "t"},
			expectedValueString: "9 t",
			description:         "Custom system should use its own boundaries and prefixes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := formatValueWithUnit(tc.input, tc.unit, tc.options)

			assert.Equal(t, tc.expectedValue, value.Scaled,
				fmt.Sprintf("Value mismatch for %s: expected %f, got %f (%s)",
					tc.name, tc.expectedValue, value.Scaled, tc.description))

			assert.Equal(t, tc.expectedPrefix, value.Prefix,
				fmt.Sprintf("Prefix mismatch for %s: expected %v, got %v (%s)",
					tc.name, tc.expectedPrefix, value.Prefix, tc.description))

			assert.Equal(t, tc.expectedUnit, value.Unit,
				fmt.Sprintf("Unit mismatch for %s: expected %v, got %v (%s)",
					tc.name, tc.expectedUnit, value.Unit, tc.description))

			assert.Equal(t, tc.expectedValueString, value.String(),
				fmt.Sprintf("Value mismatch for %s: expected %v, got %v (%s)",
					tc.name, tc.expectedValue, value.String(), tc.description))
		})
	}
}

// Test the function with multiple types to ensure generic type handling works
func TestFormatDecimalOrBinaryMultipleTypes(t *testing.T) {
	// Test with uint64
	t.Run("uint64", func(t *testing.T) {
		var input uint64 = 1500
		value := formatValueWithUnit(input, None, DefaultFormatOptions)
		assert.Equal(t, 1.5, value.Scaled)
		assert.Equal(t, prefixes[Kilo], value.Prefix)
	})

	// Test with int
	t.Run("int", func(t *testing.T) {
		var input int = 1024
		opts := DefaultFormatOptions
		opts.Boundaries = IECBinaryBoundaries
		value := formatValueWithUnit(input, None, opts)
		assert.Equal(t, 1.0, value.Scaled)
		assert.Equal(t, prefixes[Kibi], value.Prefix)
	})
}

func TestFormatBinaryUnit(t *testing.T) {
	testCases := []struct {
		name           string
		input          float64
		unit           Unit
		expectedValue  float64
		expectedPrefix Prefix
		expectedOutput string
	}{
		{
			name:           "bytes - zero value",
			input:          0,
			unit:           Byte,
			expectedValue:  0,
			expectedPrefix: Prefix{},
			expectedOutput: "0 B",
		},
		{
			name:           "bytes - small value",
			input:          42,
			unit:           Byte,
			expectedValue:  42,
			expectedPrefix: Prefix{},
			expectedOutput: "42 B",
		},
		{
			name:           "bytes - large value",
			input:          1024 * 1024 * 1024,
			unit:           Byte,
			expectedValue:  1.0,
			expectedPrefix: prefixes[Gibi],
			expectedOutput: "1 GiB",
		},
		{
			name:           "bytes - very large value",
			input:          1024 * 1024 * 1024 * 1024,
			unit:           Byte,
			expectedValue:  1.0,
			expectedPrefix: prefixes[Tebi],
			expectedOutput: "1 TiB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := FormatBinaryUnit(tc.input, tc.unit)
			assert.Equal(t, tc.expectedValue, value.Scaled, "Value mismatch")
			assert.Equal(t, tc.expectedPrefix, value.Prefix, "Prefix mismatch")
			assert.Equal(t, tc.expectedOutput, value.String(), "Output mismatch")
		})
	}
}

// Add tests for the new unit discovery functions
func TestUnitDiscovery(t *testing.T) {
	// Test basic retrieval
	names := Names()
	assert.Contains(t, names, "byte")
	assert.Contains(t, names, "bit")
	assert.Contains(t, names, "hertz")

	// Test with custom units
	UnitCentimeter, err := Register(Descriptor{
		Symbol:   "cm",
		Singular: "centimeter",
		Plural:   "centimeters",
	})
	assert.NoError(t, err)

	updatedNames := Names()
	assert.Contains(t, updatedNames, "centimeter")

	// Test unit info
	info := UnitCentimeter.Info()
	assert.Equal(t, "cm", info.Symbol)
	assert.Equal(t, "centimeter", info.Singular)
	assert.Equal(t, "centimeters", info.Plural)

	// Test discovery of all units
	units := List()
	found := false
	for _, u := range units {
		if u.Singular == "centimeter" {
			found = true
			assert.Equal(t, "cm", u.Symbol)
			assert.Equal(t, "centimeters", u.Plural)
			break
		}
	}
	assert.True(t, found, "Custom unit should be discoverable through List()")
}

func TestBinaryDecimalWithCustomUnits(t *testing.T) {
	// Register a custom unit
	UnitPixel, err := Register(Descriptor{
		Symbol:   "px",
		Singular: "pixel",
		Plural:   "pixels",
	})
	assert.NoError(t, err)

	tests := []struct {
		name     string
		value    float64
		options  []FormatOption
		expected string
	}{
		{
			name:     "binary custom unit",
			value:    1024,
			options:  []FormatOption{Binary(), WithUnit(UnitPixel)},
			expected: "1 Kipx",
		},
		{
			name:     "decimal custom unit",
			value:    1000,
			options:  []FormatOption{Decimal(), WithUnit(UnitPixel)},
			expected: "1 kpx",
		},
		{
			name:  "custom system takes precedence",
			value: 1000,
			options: []FormatOption{
				WithSystem(MustRegisterFormatSystem(
					"test",
					[]float64{1000},
					map[float64]Prefix{1000: {"thousand", "T"}},
				)),
				Binary(), // Should be ignored
				WithUnit(UnitPixel),
			},
			expected: "1 Tpx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.value, tt.options...).String()
			assert.Equal(t, tt.expected, result)
		})
	}
}
