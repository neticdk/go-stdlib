package xstrings_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/xstrings"
)

func TestCoalesce(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "Empty slice",
			args: []string{},
			want: "",
		},
		{
			name: "All empty strings",
			args: []string{"", "", ""},
			want: "",
		},
		{
			name: "First string non-empty",
			args: []string{"first", "second", "third"},
			want: "first",
		},
		{
			name: "Second string non-empty",
			args: []string{"", "second", "third"},
			want: "second",
		},
		{
			name: "Last string non-empty",
			args: []string{"", "", "third"},
			want: "third",
		},
		{
			name: "Mixed empty and non-empty",
			args: []string{"", "second", "", "fourth"},
			want: "second",
		},
		{
			name: "Single non-empty string",
			args: []string{"single"},
			want: "single",
		},
		{
			name: "Single empty string",
			args: []string{""},
			want: "",
		},
		{
			name: "Nil slice (variadic converts nil to empty)",
			args: nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xstrings.Coalesce(tt.args...)
			assert.Equal(t, got, tt.want, "Coalesce()/%s", tt.name)
		})
	}
}
