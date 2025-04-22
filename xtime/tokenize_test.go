package xtime

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []string
	}{
		// Basic cases
		{"empty string", "", []string{}},
		{"single space", " ", []string{}},
		{"multiple spaces", "   ", []string{}},
		{"single number integer", "123", []string{"123"}},
		{"single number float", "123.45", []string{"123.45"}},
		{"single unit", "h", []string{"h"}},

		// Simple pairs (no space)
		{"pair no space int", "1h", []string{"1", "h"}},
		{"pair no space float", "1.5m", []string{"1.5", "m"}},
		{"pair no space long unit", "10years", []string{"10", "years"}},

		// Simple pairs (with space)
		{"pair with space int", "1 h", []string{"1", "h"}},
		{"pair with space float", "1.5 m", []string{"1.5", "m"}},
		{"pair with space long unit", "10 years", []string{"10", "years"}},

		// Multiple pairs
		{"multiple pairs no spaces", "1h30m15s", []string{"1", "h", "30", "m", "15", "s"}},
		{"multiple pairs with spaces", "1 h 30 m 15 s", []string{"1", "h", "30", "m", "15", "s"}},
		{"multiple pairs mixed spaces 1", "1h 30m 15s", []string{"1", "h", "30", "m", "15", "s"}},
		{"multiple pairs mixed spaces 2", " 1h30m 15s ", []string{"1", "h", "30", "m", "15", "s"}},
		{"multiple pairs floats", "1.5h 30.2m 0.5s", []string{"1.5", "h", "30.2", "m", "0.5", "s"}},

		// Unit followed by number (should split)
		{"unit then number", "h1", []string{"h", "1"}},
		{"unit then float", "m1.5", []string{"m", "1.5"}},
		{"unit then decimal start", "s.5", []string{"s", ".5"}},
		{"multiple unit then number", "h1m2s3", []string{"h", "1", "m", "2", "s", "3"}},

		// Decimal handling
		{"leading decimal number", ".5", []string{"0.5"}},
		{"leading decimal pair", ".5h", []string{"0.5", "h"}},
		{"number with trailing decimal", "1.", []string{"1."}},
		{"pair with trailing decimal", "1.h", []string{"1.", "h"}},
		{"multiple decimals", "1.2.3", []string{"1.2.3"}},
		{"multiple decimals in pair", "1.2.3h", []string{"1.2.3", "h"}},
		{"decimal only", ".", []string{"0."}},
		{"decimal only unit", ".h", []string{"0.", "h"}},

		// Edge cases with spaces
		{"trailing space", "1h ", []string{"1", "h"}},
		{"leading space", " 1h", []string{"1", "h"}},
		{"spaces around", " 1h30m ", []string{"1", "h", "30", "m"}},
		{"multiple internal spaces", "1  h   30    m", []string{"1", "h", "30", "m"}},

		// Longer units mixed
		{"long units mixed", "1hour30mins", []string{"1", "hour", "30", "mins"}},
		{"long units spaces", "1 hour 30 mins", []string{"1", "hour", "30", "mins"}},
		{"long units mixed three units", "1hour30mins10seconds", []string{"1", "hour", "30", "mins", "10", "seconds"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Note: These tests assume the input to tokenize is already "cleaned"
			// (no commas, no "and ").
			got := tokenize(tc.input)

			assert.Equal(t, got, tc.want, "tokenize/%s", tc.name)
		})
	}
}
