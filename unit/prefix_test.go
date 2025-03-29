package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixFor(t *testing.T) {
	testCases := []struct {
		name     string
		value    float64
		expected Prefix
	}{
		{
			name:     "kilo prefix",
			value:    Kilo,
			expected: Prefix{Name: "kilo", Symbol: "k"},
		},
		{
			name:     "mega prefix",
			value:    Mega,
			expected: Prefix{Name: "mega", Symbol: "M"},
		},
		{
			name:     "giga prefix",
			value:    Giga,
			expected: Prefix{Name: "giga", Symbol: "G"},
		},
		{
			name:     "tera prefix",
			value:    Tera,
			expected: Prefix{Name: "tera", Symbol: "T"},
		},
		{
			name:     "peta prefix",
			value:    Peta,
			expected: Prefix{Name: "peta", Symbol: "P"},
		},
		{
			name:     "exa prefix",
			value:    Exa,
			expected: Prefix{Name: "exa", Symbol: "E"},
		},
		{
			name:     "zetta prefix",
			value:    Zetta,
			expected: Prefix{Name: "zetta", Symbol: "Z"},
		},
		{
			name:     "yotta prefix",
			value:    Yotta,
			expected: Prefix{Name: "yotta", Symbol: "Y"},
		},
		{
			name:     "ronna prefix",
			value:    Ronna,
			expected: Prefix{Name: "ronna", Symbol: "R"},
		},
		{
			name:     "quetta prefix",
			value:    Quetta,
			expected: Prefix{Name: "quetta", Symbol: "Q"},
		},
		{
			name:     "kibi prefix",
			value:    Kibi,
			expected: Prefix{Name: "kibi", Symbol: "Ki"},
		},
		{
			name:     "mebi prefix",
			value:    Mebi,
			expected: Prefix{Name: "mebi", Symbol: "Mi"},
		},
		{
			name:     "gibi prefix",
			value:    Gibi,
			expected: Prefix{Name: "gibi", Symbol: "Gi"},
		},
		{
			name:     "tebi prefix",
			value:    Tebi,
			expected: Prefix{Name: "tebi", Symbol: "Ti"},
		},
		{
			name:     "pebi prefix",
			value:    Pebi,
			expected: Prefix{Name: "pebi", Symbol: "Pi"},
		},
		{
			name:     "exbi prefix",
			value:    Exbi,
			expected: Prefix{Name: "exbi", Symbol: "Ei"},
		},
		{
			name:     "zebi prefix",
			value:    Zebi,
			expected: Prefix{Name: "zebi", Symbol: "Zi"},
		},
		{
			name:     "yobi prefix",
			value:    Yobi,
			expected: Prefix{Name: "yobi", Symbol: "Yi"},
		},
		{
			name:     "robi prefix",
			value:    Robi,
			expected: Prefix{Name: "robi", Symbol: "Ri"},
		},
		{
			name:     "quebi prefix",
			value:    Quebi,
			expected: Prefix{Name: "quebi", Symbol: "Qi"},
		},
		{
			name:     "unknown prefix",
			value:    123456,
			expected: Prefix{},
		},
		{
			name:     "zero value",
			value:    0,
			expected: Prefix{},
		},
		{
			name:     "negative value",
			value:    -1000,
			expected: Prefix{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PrefixFor(tc.value)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPrefixProperties(t *testing.T) {
	// Test that all prefixes in the map have expected properties
	for value, prefix := range prefixes {
		t.Run(prefix.Name, func(t *testing.T) {
			// Name should not be empty
			assert.NotEmpty(t, prefix.Name, "Prefix name should not be empty")

			// Symbol should not be empty
			assert.NotEmpty(t, prefix.Symbol, "Prefix symbol should not be empty")

			// Test lookup consistency
			lookupResult := PrefixFor(value)
			assert.Equal(t, prefix, lookupResult, "Prefix lookup should be consistent")
		})
	}
}

func TestPrefixMapping(t *testing.T) {
	// Make sure our decimal boundaries map to the right prefixes
	assert.Equal(t, prefixes[Kilo], PrefixFor(SIDecimalBoundaries[9]))
	assert.Equal(t, prefixes[Mega], PrefixFor(SIDecimalBoundaries[8]))
	assert.Equal(t, prefixes[Giga], PrefixFor(SIDecimalBoundaries[7]))
	assert.Equal(t, prefixes[Tera], PrefixFor(SIDecimalBoundaries[6]))
	assert.Equal(t, prefixes[Peta], PrefixFor(SIDecimalBoundaries[5]))
	assert.Equal(t, prefixes[Exa], PrefixFor(SIDecimalBoundaries[4]))
	assert.Equal(t, prefixes[Zetta], PrefixFor(SIDecimalBoundaries[3]))
	assert.Equal(t, prefixes[Yotta], PrefixFor(SIDecimalBoundaries[2]))
	assert.Equal(t, prefixes[Ronna], PrefixFor(SIDecimalBoundaries[1]))
	assert.Equal(t, prefixes[Quetta], PrefixFor(SIDecimalBoundaries[0]))

	// Make sure our binary boundaries map to the right prefixes
	assert.Equal(t, prefixes[Kibi], PrefixFor(IECBinaryBoundaries[9]))
	assert.Equal(t, prefixes[Mebi], PrefixFor(IECBinaryBoundaries[8]))
	assert.Equal(t, prefixes[Gibi], PrefixFor(IECBinaryBoundaries[7]))
	assert.Equal(t, prefixes[Tebi], PrefixFor(IECBinaryBoundaries[6]))
	assert.Equal(t, prefixes[Pebi], PrefixFor(IECBinaryBoundaries[5]))
	assert.Equal(t, prefixes[Exbi], PrefixFor(IECBinaryBoundaries[4]))
	assert.Equal(t, prefixes[Zebi], PrefixFor(IECBinaryBoundaries[3]))
	assert.Equal(t, prefixes[Yobi], PrefixFor(IECBinaryBoundaries[2]))
	assert.Equal(t, prefixes[Robi], PrefixFor(IECBinaryBoundaries[1]))
	assert.Equal(t, prefixes[Quebi], PrefixFor(IECBinaryBoundaries[0]))
}

// Benchmark PrefixFor function
func BenchmarkPrefixFor(b *testing.B) {
	values := []float64{
		Kilo, Mega, Giga, Tera, Peta,
		Kibi, Mebi, Gibi, Tebi, Pebi,
		123456, // unknown value
	}

	b.ResetTimer()
	for b.Loop() {
		for _, v := range values {
			PrefixFor(v)
		}
	}
}
