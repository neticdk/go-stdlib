package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	resetGlobalUnitState()

	tests := []struct {
		input   string
		want    ParseResult
		wantErr bool
		errMsg  string
	}{
		{
			input: "1.5 KB",
			want:  ParseResult{Raw: 1.5, Scale: Kilo, Unit: Byte},
		},
		{
			input: "1 MiB",
			want:  ParseResult{Raw: 1, Scale: Mebi, Unit: Byte},
		},
		{
			input: "500 MHz",
			want:  ParseResult{Raw: 500, Scale: Mega, Unit: Hertz},
		},
		{
			input: "2.7 kb",
			want:  ParseResult{Raw: 2.7, Scale: Kilo, Unit: Bit},
		},
		{
			input: "1024",
			want:  ParseResult{Raw: 1024, Scale: 1, Unit: None},
		},
		{
			input:   "invalid",
			wantErr: true,
			errMsg:  "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParse_Custom(t *testing.T) {
	resetGlobalUnitState()

	// Register Custom Unit: Pixel
	UnitPixel, errPx := Register(Descriptor{
		Symbol:   "px",
		Singular: "pixel",
		Plural:   "pixels",
	})
	require.NoError(t, errPx, "Setup: Failed to register UnitPixel")
	require.NotEqual(t, None, UnitPixel, "Setup: UnitPixel should not be None")

	// Register Custom Unit: Cycle
	UnitCycle, errCyc := Register(Descriptor{
		Symbol:   "cyc",
		Singular: "cycle",
		Plural:   "cycles",
	})
	require.NoError(t, errCyc, "Setup: Failed to register UnitCycle")
	require.NotEqual(t, None, UnitCycle, "Setup: UnitCycle should not be None")

	// Register Custom Format System: Widgets
	_, err := RegisterFormatSystem(
		"widgets",
		[]float64{1000.0}, // Boundaries (just one for simplicity)
		map[float64]Prefix{
			1000.0: {Name: "kiloWidget", Symbol: "kWid"},
		},
	)
	assert.NoError(t, err, "Setup: Failed to register FormatSystem 'widgets'")

	testCases := []struct {
		name           string
		input          string
		want           ParseResult
		wantErr        bool
		errMsgContains string
	}{
		// Tests with Custom Unit 'Pixel'
		{
			name:  "custom unit symbol, standard decimal prefix",
			input: "2 kpx",
			want:  ParseResult{Raw: 2, Scale: Kilo, Unit: UnitPixel},
		},
		{
			name:  "custom unit symbol, standard decimal prefix (uppercase K alias)",
			input: "1.5 Kpx",
			want:  ParseResult{Raw: 1.5, Scale: Kilo, Unit: UnitPixel},
		},
		{
			name:  "custom unit symbol, standard binary prefix",
			input: "4 Kipx",
			want:  ParseResult{Raw: 4, Scale: Kibi, Unit: UnitPixel},
		},
		{
			name:  "custom unit symbol, no prefix",
			input: "100 px",
			want:  ParseResult{Raw: 100, Scale: 1, Unit: UnitPixel},
		},
		{
			name:  "custom unit singular name, no prefix",
			input: "50 pixel",
			want:  ParseResult{Raw: 50, Scale: 1, Unit: UnitPixel},
		},
		{
			name:  "custom unit plural name, no prefix",
			input: "75 pixels",
			want:  ParseResult{Raw: 75, Scale: 1, Unit: UnitPixel},
		},
		{
			name:  "custom unit symbol, with space before",
			input: " 100 px",
			want:  ParseResult{Raw: 100, Scale: 1, Unit: UnitPixel},
		},
		{
			name:  "custom unit symbol, with space after",
			input: "100 px ",
			want:  ParseResult{Raw: 100, Scale: 1, Unit: UnitPixel},
		},
		{
			name:  "custom unit symbol, standard prefix, no space",
			input: "3Mpx",
			want:  ParseResult{Raw: 3, Scale: Mega, Unit: UnitPixel},
		},

		// Tests with Custom Unit 'Cycle'
		{
			name:  "custom unit 'cycle', standard prefix",
			input: "5 Mcyc",
			want:  ParseResult{Raw: 5, Scale: Mega, Unit: UnitCycle},
		},
		{
			name:  "custom unit 'cycle', standard prefix (K alias)",
			input: "1.5 Kcyc", // Clearer test for K alias
			want:  ParseResult{Raw: 1.5, Scale: Kilo, Unit: UnitCycle},
		},
		{
			name:  "custom unit 'cycle', plural name",
			input: "10 cycles",
			want:  ParseResult{Raw: 10, Scale: 1, Unit: UnitCycle},
		},

		// Tests with Custom Format System 'Widgets'
		{
			name:  "custom system prefix symbol, custom unit symbol",
			input: "3 kWid cyc",
			want:  ParseResult{Raw: 3, Scale: 1000.0, Unit: UnitCycle},
		},
		{
			name:  "custom system prefix symbol, custom unit name (singular)",
			input: "3 kWid cycle",
			want:  ParseResult{Raw: 3, Scale: 1000.0, Unit: UnitCycle},
		},
		{
			name:  "custom system prefix symbol, custom unit name (plural)",
			input: "3 kWid cycles",
			want:  ParseResult{Raw: 3, Scale: 1000.0, Unit: UnitCycle},
		},
		{
			name:  "custom system prefix name, custom unit symbol",
			input: "1.5 kiloWidget cyc",
			want:  ParseResult{Raw: 1.5, Scale: 1000.0, Unit: UnitCycle},
		},
		{
			name:  "custom system prefix symbol only",
			input: "7 kWid",
			want:  ParseResult{Raw: 7, Scale: 1000.0, Unit: None},
		},
		{
			name:  "custom system prefix name only",
			input: "8 kiloWidget",
			want:  ParseResult{Raw: 8, Scale: 1000.0, Unit: None},
		},
		{
			name:  "custom system prefix symbol, standard unit",
			input: "9 kWid B", // kiloWidget Bytes
			want:  ParseResult{Raw: 9, Scale: 1000.0, Unit: Byte},
		},

		// Error Cases
		{
			name:           "error - known custom prefix, unknown unit",
			input:          "2 kWid foobar",
			wantErr:        true,
			errMsgContains: `unknown unit or prefix combination: "kWid foobar"`,
		},
		{
			name:           "error - unknown prefix, known custom unit",
			input:          "9 Zzz px",
			wantErr:        true,
			errMsgContains: `unknown unit or prefix combination: "Zzz px"`,
		},
		{
			name:           "error - known custom unit, invalid number",
			input:          "abc px",
			wantErr:        true,
			errMsgContains: `invalid format - no valid number found`,
		},
		{
			name:           "error - only unknown string",
			input:          "foobar",
			wantErr:        true,
			errMsgContains: `invalid format - no valid number found`, // Because it first fails to parse as just a number
		},
		{
			name:           "error - number followed by unknown string",
			input:          "123 foobar",
			wantErr:        true,
			errMsgContains: `unknown unit or prefix combination: "foobar"`,
		},
		{
			name:           "error - ambiguous KBpx format (Kilo Byte Pixel)",
			input:          "1.5 KBpx",
			wantErr:        true,
			errMsgContains: `unknown unit or prefix combination: "KBpx"`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsgContains != "" {
					assert.Contains(t, err.Error(), tt.errMsgContains, "Error message mismatch")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Raw, got.Raw, "Raw value mismatch")
				assert.Equal(t, tt.want.Scale, got.Scale, "Scale mismatch")
				assert.Equal(t, tt.want.Unit, got.Unit, "Unit mismatch")
			}
		})
	}
}

// Package-level variables to store results and prevent compiler optimization.
var (
	benchParseResult ParseResult
	benchParseErr    error
)

// BenchmarkParse measures the performance of the unit.Parse function.
func BenchmarkParse(b *testing.B) {
	resetGlobalUnitState()

	// Setup: Ensure Custom Elements are Registered
	// Use MustRegister for simplicity in benchmark setup, assuming registration won't fail here.
	// If registration could potentially fail mid-benchmark run (not typical), more robust error handling needed.
	MustRegister(Descriptor{
		Symbol:   "px",
		Singular: "pixel",
		Plural:   "pixels",
	})
	MustRegister(Descriptor{
		Symbol:   "cyc",
		Singular: "cycle",
		Plural:   "cycles",
	})
	MustRegisterFormatSystem(
		"widgets",
		[]float64{1000.0},
		map[float64]Prefix{
			1000.0: {Name: "kiloWidget", Symbol: "kWid"},
		},
	)
	// Ensure lookups are rebuilt (implicitly done by Register funcs)

	// Input Data: A mix of valid and invalid cases ---
	inputs := []string{
		// Standard SI
		"1024",
		"1.5k",
		"2 kB",
		"500 MHz",
		"1GB",
		"0.9 TB",
		"10 KB",

		// Standard IEC
		"1 KiB",
		"1024 MiB",
		"2.5GiB",

		// Custom Units with Standard Prefixes
		"1 Mpx",
		"10kcyc",
		"2 Kpx",
		"500 Gipx",
		"3 Tibytes",
		"75 pixels",
		"1 cycle",

		// Custom Prefixes
		"3 kWid cyc",
		"1.5 kiloWidget px",
		"7 kWid",
		"8 kiloWidget",

		// Edge cases / Whitespace
		" 100 B ",
		"   5 Mpx   ",

		// Invalid Cases
		"invalid",
		"abc px",
		"1.2.3 GB",
		"5 kWid unknown",
		"10 Zzz",
	}

	// Ensure input slice is not empty
	if len(inputs) == 0 {
		b.Fatal("Input slice for benchmark is empty")
	}

	b.ReportAllocs() // Report memory allocations per operation
	b.ResetTimer()   // Start timing precisely here, excluding setup

	// --- The Benchmark Loop ---
	// The loop runs b.N times. We cycle through the inputs using the modulo operator.
	// This gives a performance average across the different input types.
	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)] // Cycle through the inputs
		// Assign to package-level vars to prevent optimization
		benchParseResult, benchParseErr = Parse(input)
	}
}

// WARNING: This function manipulates global state and should ONLY be used in tests
// to ensure isolation between test functions that modify the unit/prefix registries.
func resetGlobalUnitState() {
	// Reset Unit Registry to initial built-in state
	unitRegistryMutex.Lock()
	nextUnitID = unitMaxBuiltin // Reset the ID counter
	// Clear the map and re-add built-ins
	unitRegistry = make(map[Unit]Descriptor) // Assign a new empty map
	unitRegistry[None] = Descriptor{}        // Not strictly needed if Describe handles nil lookup
	unitRegistry[Byte] = Descriptor{Byte, "B", "byte", "bytes"}
	unitRegistry[Bit] = Descriptor{Bit, "b", "bit", "bits"}
	unitRegistry[Hertz] = Descriptor{Hertz, "Hz", "hertz", "hertz"}
	rebuildUnitLookups() // Rebuild lookups based ONLY on built-ins now
	unitRegistryMutex.Unlock()

	// Reset Format Systems Registry to empty state
	formatSystemsMutex.Lock()
	formatSystems = make(map[string]*FormatSystem)
	rebuildPrefixLookups() // Rebuild prefix lookup (will only contain default SI/IEC now)
	formatSystemsMutex.Unlock()
}
