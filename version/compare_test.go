package version

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestFirst(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		expected string
	}{
		{
			name:     "No versions",
			versions: []string{},
			expected: "",
		},

		{
			name:     "All empty versions",
			versions: []string{"", "", ""},
			expected: "",
		},
		{
			name:     "First non-empty version",
			versions: []string{"", "1.0.0", "2.0.0"},
			expected: "1.0.0",
		},
		{
			name:     "First version non-empty",
			versions: []string{"1.0.0", "2.0.0", "3.0.0"},
			expected: "1.0.0",
		},
		{
			name:     "Middle version non-empty",
			versions: []string{"", "2.0.0", ""},
			expected: "2.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := First(tt.versions...)
			assert.Equal(t, result, tt.expected, "First/%q", tt.name)
		})
	}
}
