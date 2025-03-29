package version

import "testing"

func TestFirst(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		expected string
	}{
		{"No versions", []string{}, ""},
		{"All empty versions", []string{"", "", ""}, ""},
		{"First non-empty version", []string{"", "1.0.0", "2.0.0"}, "1.0.0"},
		{"First version non-empty", []string{"1.0.0", "2.0.0", "3.0.0"}, "1.0.0"},
		{"Middle version non-empty", []string{"", "2.0.0", ""}, "2.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := First(tt.versions...)
			if result != tt.expected {
				t.Errorf("First(%v) = %v; want %v", tt.versions, result, tt.expected)
			}
		})
	}
}

