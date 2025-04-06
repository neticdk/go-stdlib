package diffcore_test

import (
	"testing"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/diff"
	"github.com/neticdk/go-stdlib/diff/internal/diffcore"
)

func TestGroupEditsByContext(t *testing.T) {
	tests := []struct {
		name     string
		edits    []diff.Line
		context  int
		expected [][]diff.Line
	}{
		{
			name:     "empty edits",
			edits:    []diff.Line{},
			context:  2,
			expected: [][]diff.Line{{}},
		},
		{
			name:     "empty diff",
			edits:    []diff.Line{{Kind: diff.Equal, Text: "line1"}},
			context:  1,
			expected: [][]diff.Line{{{Kind: diff.Equal, Text: "line1"}}},
		},
		{
			name: "no context",
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Insert, Text: "line2"},
				{Kind: diff.Delete, Text: "line3"},
				{Kind: diff.Equal, Text: "line4"},
			},
			context: 0,
			expected: [][]diff.Line{
				{
					{Kind: diff.Insert, Text: "line2"},
					{Kind: diff.Delete, Text: "line3"},
				},
			},
		},
		{
			name: "all equals",
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Equal, Text: "line2"},
				{Kind: diff.Equal, Text: "line3"},
			},
			context: 2,
			expected: [][]diff.Line{
				{
					{Kind: diff.Equal, Text: "line1"},
					{Kind: diff.Equal, Text: "line2"},
					{Kind: diff.Equal, Text: "line3"},
				},
			},
		},
		{
			name: "single change",
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Insert, Text: "line2"},
				{Kind: diff.Equal, Text: "line3"},
			},
			context: 1,
			expected: [][]diff.Line{
				{
					{Kind: diff.Equal, Text: "line1"},
					{Kind: diff.Insert, Text: "line2"},
					{Kind: diff.Equal, Text: "line3"},
				},
			},
		},
		{
			name: "multiple changes within context",
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Insert, Text: "line2"},
				{Kind: diff.Delete, Text: "line3"},
				{Kind: diff.Equal, Text: "line4"},
			},
			context: 2,
			expected: [][]diff.Line{
				{
					{Kind: diff.Equal, Text: "line1"},
					{Kind: diff.Insert, Text: "line2"},
					{Kind: diff.Delete, Text: "line3"},
					{Kind: diff.Equal, Text: "line4"},
				},
			},
		},
		{
			name: "multiple changes outside context",
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Insert, Text: "line2"},
				{Kind: diff.Equal, Text: "line3"},
				{Kind: diff.Delete, Text: "line4"},
				{Kind: diff.Equal, Text: "line5"},
			},
			context: 0,
			expected: [][]diff.Line{
				{
					{Kind: diff.Insert, Text: "line2"},
					{Kind: diff.Delete, Text: "line4"},
				},
			},
		},
		{
			name: "multiple changes outside context with context lines",
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Insert, Text: "line2"},
				{Kind: diff.Equal, Text: "line3"},
				{Kind: diff.Equal, Text: "line4"},
				{Kind: diff.Equal, Text: "line5"},
				{Kind: diff.Delete, Text: "line6"},
				{Kind: diff.Equal, Text: "line7"},
			},
			context: 1,
			expected: [][]diff.Line{
				{
					{Kind: diff.Equal, Text: "line1"},
					{Kind: diff.Insert, Text: "line2"},
					{Kind: diff.Equal, Text: "line3"},
				},
				{
					{Kind: diff.Equal, Text: "line5"},
					{Kind: diff.Delete, Text: "line6"},
					{Kind: diff.Equal, Text: "line7"},
				},
			},
		},
		{
			name: "multiple groups",
			edits: []diff.Line{
				{Kind: diff.Equal, Text: "line1"},
				{Kind: diff.Insert, Text: "line2"},
				{Kind: diff.Equal, Text: "line3"},
				{Kind: diff.Equal, Text: "line4"},
				{Kind: diff.Equal, Text: "line5"},
				{Kind: diff.Delete, Text: "line6"},
				{Kind: diff.Equal, Text: "line7"},
				{Kind: diff.Equal, Text: "line8"},
				{Kind: diff.Insert, Text: "line9"},
				{Kind: diff.Equal, Text: "line10"},
			},
			context: 1,
			expected: [][]diff.Line{
				{
					{Kind: diff.Equal, Text: "line1"},
					{Kind: diff.Insert, Text: "line2"},
					{Kind: diff.Equal, Text: "line3"},
				},
				{
					{Kind: diff.Equal, Text: "line5"},
					{Kind: diff.Delete, Text: "line6"},
					{Kind: diff.Equal, Text: "line7"},
				},
				{
					{Kind: diff.Equal, Text: "line8"},
					{Kind: diff.Insert, Text: "line9"},
					{Kind: diff.Equal, Text: "line10"},
				},
			},
		},
		{
			name: "all changes",
			edits: []diff.Line{
				{Kind: diff.Insert, Text: "line1"},
				{Kind: diff.Delete, Text: "line2"},
				{Kind: diff.Equal, Text: "line3"},
			},
			context: 1,
			expected: [][]diff.Line{
				{
					{Kind: diff.Insert, Text: "line1"},
					{Kind: diff.Delete, Text: "line2"},
					{Kind: diff.Equal, Text: "line3"},
				},
			},
		},
		{
			name: "context exceeds edits length",
			edits: []diff.Line{
				{Kind: diff.Insert, Text: "line1"},
				{Kind: diff.Delete, Text: "line2"},
			},
			context: 5,
			expected: [][]diff.Line{
				{
					{Kind: diff.Insert, Text: "line1"},
					{Kind: diff.Delete, Text: "line2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := diffcore.GroupEditsByContext(tt.edits, tt.context)
			assert.Equal(t, result, tt.expected, "GroupEditsByContext/%q", tt.name)
		})
	}
}
