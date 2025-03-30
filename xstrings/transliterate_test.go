package xstrings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransliterate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "ascii only",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "accented characters",
			input:    "Héllò Wórld",
			expected: "Hello World",
		},
		{
			name:     "cyrillic",
			input:    "Привет мир",
			expected: "Privet mir",
		},
		{
			name:     "chinese",
			input:    "你好，世界",
			expected: "Ni Hao ,Shi Jie ",
		},
		{
			name:     "mixed characters",
			input:    "Café परीक्षा 测试 123",
			expected: "Cafe priikssaa Ce Shi  123",
		},
		{
			name:     "symbols",
			input:    "©®™",
			expected: "(c)(r)(tm)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Transliterate(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
