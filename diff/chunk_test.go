package diff

import (
	"reflect"
	"testing"
)

func TestCalculateChunkRanges(t *testing.T) {
	tests := []struct {
		name     string
		edits    []Line
		context  int
		expected []chunkRange
	}{
		{
			name:     "empty edits",
			edits:    []Line{},
			context:  2,
			expected: []chunkRange{},
		},
		{
			name:     "empty diff",
			edits:    []Line{{Kind: Equal, Text: "line1"}},
			context:  1,
			expected: []chunkRange{{start: 0, end: 1, isNewGroup: false}},
		},
		{
			name: "no context",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Delete, Text: "line3"},
				{Kind: Equal, Text: "line4"},
			},
			context: 0,
			expected: []chunkRange{
				{start: 1, end: 2, isNewGroup: true},
				{start: 2, end: 3, isNewGroup: true},
			},
		},
		{
			name: "all equals",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Equal, Text: "line2"},
				{Kind: Equal, Text: "line3"},
			},
			context:  2,
			expected: []chunkRange{{start: 0, end: 3, isNewGroup: false}},
		},
		{
			name: "single change",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Equal, Text: "line3"},
			},
			context:  1,
			expected: []chunkRange{{start: 0, end: 3, isNewGroup: false}},
		},
		{
			name: "multiple changes within context",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Delete, Text: "line3"},
				{Kind: Equal, Text: "line4"},
			},
			context:  2,
			expected: []chunkRange{{start: 0, end: 4, isNewGroup: false}},
		},
		{
			name: "multiple changes outside context",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Equal, Text: "line3"},
				{Kind: Delete, Text: "line4"},
				{Kind: Equal, Text: "line5"},
			},
			context: 0,
			expected: []chunkRange{
				{start: 1, end: 2, isNewGroup: true},
				{start: 3, end: 4, isNewGroup: true},
			},
		},
		{
			name: "multiple changes outside context with context lines",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Equal, Text: "line3"},
				{Kind: Equal, Text: "line4"},
				{Kind: Equal, Text: "line5"},
				{Kind: Delete, Text: "line6"},
				{Kind: Equal, Text: "line7"},
			},
			context: 1,
			expected: []chunkRange{
				{start: 0, end: 3, isNewGroup: false},
				{start: 4, end: 7, isNewGroup: true},
			},
		},
		{
			name: "multiple groups",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Equal, Text: "line3"},
				{Kind: Equal, Text: "line4"},
				{Kind: Equal, Text: "line5"},
				{Kind: Delete, Text: "line6"},
				{Kind: Equal, Text: "line7"},
				{Kind: Equal, Text: "line8"},
				{Kind: Insert, Text: "line9"},
				{Kind: Equal, Text: "line10"},
			},
			context: 1,
			expected: []chunkRange{
				{start: 0, end: 3, isNewGroup: false},
				{start: 4, end: 7, isNewGroup: true},
				{start: 7, end: 10, isNewGroup: true},
			},
		},
		{
			name: "all changes",
			edits: []Line{
				{Kind: Insert, Text: "line1"},
				{Kind: Delete, Text: "line2"},
				{Kind: Equal, Text: "line3"},
			},
			context:  1,
			expected: []chunkRange{{start: 0, end: 3, isNewGroup: false}},
		},
		{
			name: "context exceeds edits length",
			edits: []Line{
				{Kind: Insert, Text: "line1"},
				{Kind: Delete, Text: "line2"},
			},
			context:  5,
			expected: []chunkRange{{start: 0, end: 2, isNewGroup: false}},
		},
		{
			name: "no actual edits; just context",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Equal, Text: "line2"},
				{Kind: Equal, Text: "line3"},
			},
			context:  1,
			expected: []chunkRange{{start: 0, end: 3, isNewGroup: false}},
		},
		{
			name: "long sequence with all changes and no context",
			edits: []Line{
				{Kind: Insert, Text: "insert01"},
				{Kind: Delete, Text: "delete02"},
				{Kind: Insert, Text: "insert03"},
				{Kind: Delete, Text: "delete04"},
				{Kind: Insert, Text: "insert05"},
				{Kind: Delete, Text: "delete06"},
				{Kind: Insert, Text: "insert07"},
				{Kind: Delete, Text: "delete08"},
				{Kind: Insert, Text: "insert09"},
			},
			context: 0,
			expected: []chunkRange{
				{start: 0, end: 1, isNewGroup: true},
				{start: 1, end: 2, isNewGroup: true},
				{start: 2, end: 3, isNewGroup: true},
				{start: 3, end: 4, isNewGroup: true},
				{start: 4, end: 5, isNewGroup: true},
				{start: 5, end: 6, isNewGroup: true},
				{start: 6, end: 7, isNewGroup: true},
				{start: 7, end: 8, isNewGroup: true},
				{start: 8, end: 9, isNewGroup: true},
			},
		},
		{
			name: "long sequence with multiple groups and context",
			edits: []Line{
				{Kind: Equal, Text: "line01"},
				{Kind: Equal, Text: "line02"},
				{Kind: Equal, Text: "line03"},
				{Kind: Insert, Text: "insert04"},
				{Kind: Equal, Text: "line05"},
				{Kind: Equal, Text: "line06"},
				{Kind: Equal, Text: "line07"},
				{Kind: Delete, Text: "delete08"},
				{Kind: Equal, Text: "line09"},
				{Kind: Equal, Text: "line10"},
				{Kind: Insert, Text: "insert11"},
				{Kind: Equal, Text: "line12"},
				{Kind: Equal, Text: "line13"},
				{Kind: Equal, Text: "line14"},
				{Kind: Equal, Text: "line15"},
				{Kind: Delete, Text: "delete16"},
				{Kind: Equal, Text: "line17"},
				{Kind: Equal, Text: "line18"},
				{Kind: Equal, Text: "line19"},
				{Kind: Insert, Text: "insert20"},
				{Kind: Equal, Text: "line21"},
				{Kind: Equal, Text: "line22"},
				{Kind: Equal, Text: "line23"},
				{Kind: Equal, Text: "line24"},
				{Kind: Delete, Text: "delete25"},
				{Kind: Equal, Text: "line26"},
				{Kind: Equal, Text: "line27"},
			},
			context: 2,
			expected: []chunkRange{
				{start: 1, end: 13, isNewGroup: false}, // Covers first two changes and their context
				{start: 13, end: 22, isNewGroup: true}, // Covers middle changes and their context
				{start: 22, end: 27, isNewGroup: true}, // Covers final changes and their context
			},
		},
		{
			name: "overlapping changes with small context",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Delete, Text: "line2"},
				{Kind: Insert, Text: "line3"},
				{Kind: Equal, Text: "line4"},
				{Kind: Delete, Text: "line5"},
				{Kind: Insert, Text: "line6"},
				{Kind: Equal, Text: "line7"},
			},
			context: 1,
			expected: []chunkRange{
				{start: 0, end: 7, isNewGroup: false},
			},
		},
		{
			name: "overlapping changes with large context",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Delete, Text: "line2"},
				{Kind: Insert, Text: "line3"},
				{Kind: Equal, Text: "line4"},
				{Kind: Delete, Text: "line5"},
				{Kind: Insert, Text: "line6"},
				{Kind: Equal, Text: "line7"},
			},
			context:  3,
			expected: []chunkRange{{start: 0, end: 7, isNewGroup: false}},
		},
		{
			name: "tightly spaced changes with minimal context",
			edits: []Line{
				{Kind: Delete, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Delete, Text: "line3"},
				{Kind: Insert, Text: "line4"},
				{Kind: Delete, Text: "line5"},
				{Kind: Insert, Text: "line6"},
			},
			context:  1,
			expected: []chunkRange{{start: 0, end: 6, isNewGroup: false}},
		},
		{
			name: "tightly spaced changes with no context",
			edits: []Line{
				{Kind: Delete, Text: "line1"},
				{Kind: Insert, Text: "line2"},
				{Kind: Delete, Text: "line3"},
				{Kind: Insert, Text: "line4"},
				{Kind: Delete, Text: "line5"},
				{Kind: Insert, Text: "line6"},
			},
			context: 0,
			expected: []chunkRange{
				{start: 0, end: 1, isNewGroup: true},
				{start: 1, end: 2, isNewGroup: true},
				{start: 2, end: 3, isNewGroup: true},
				{start: 3, end: 4, isNewGroup: true},
				{start: 4, end: 5, isNewGroup: true},
				{start: 5, end: 6, isNewGroup: true},
			},
		},
		{
			name: "many small overlapping changes",
			edits: []Line{
				{Kind: Equal, Text: "a"},
				{Kind: Delete, Text: "b"},
				{Kind: Insert, Text: "c"},
				{Kind: Equal, Text: "d"},
				{Kind: Delete, Text: "e"},
				{Kind: Insert, Text: "f"},
				{Kind: Equal, Text: "g"},
				{Kind: Delete, Text: "h"},
				{Kind: Insert, Text: "i"},
				{Kind: Equal, Text: "j"},
			},
			context:  1,
			expected: []chunkRange{{start: 0, end: 10, isNewGroup: false}},
		},
		{
			name: "ensure group marking with overlap",
			edits: []Line{
				{Kind: Equal, Text: "a"},
				{Kind: Delete, Text: "b"},
				{Kind: Insert, Text: "c"},
				{Kind: Equal, Text: "d"},
				{Kind: Equal, Text: "e"},
				{Kind: Equal, Text: "f"},
				{Kind: Equal, Text: "g"},
				{Kind: Delete, Text: "h"},
				{Kind: Insert, Text: "i"},
				{Kind: Equal, Text: "j"},
			},
			context: 1,
			expected: []chunkRange{
				{start: 0, end: 4, isNewGroup: false},
				{start: 6, end: 10, isNewGroup: true},
			},
		},
		{
			name: "edge case where merge condition triggers but shouldn't",
			edits: []Line{
				{Kind: Equal, Text: "line1"},
				{Kind: Equal, Text: "line2"},
				{Kind: Equal, Text: "line3"},
				{Kind: Equal, Text: "line4"},
				{Kind: Delete, Text: "line5"},
				{Kind: Equal, Text: "line6"},
				{Kind: Equal, Text: "line7"},
				{Kind: Equal, Text: "line8"},
				{Kind: Equal, Text: "line9"},
				{Kind: Equal, Text: "line10"},
				{Kind: Delete, Text: "line11"},
				{Kind: Equal, Text: "line12"},
			},
			context: 3,
			expected: []chunkRange{
				{start: 1, end: 12, isNewGroup: false},
			},
		},
		{
			name: "extreme overlap with context = 1",
			edits: []Line{
				{Kind: Delete, Text: "1"},
				{Kind: Insert, Text: "2"},
				{Kind: Delete, Text: "3"},
				{Kind: Insert, Text: "4"},
				{Kind: Delete, Text: "5"},
				{Kind: Insert, Text: "6"},
				{Kind: Delete, Text: "7"},
				{Kind: Insert, Text: "8"},
				{Kind: Delete, Text: "9"},
				{Kind: Insert, Text: "10"},
			},
			context:  1,
			expected: []chunkRange{{start: 0, end: 10, isNewGroup: false}},
		},
		{
			name: "context = 2, many overlapping edits",
			edits: []Line{
				{Kind: Equal, Text: "0"},
				{Kind: Delete, Text: "1"},
				{Kind: Insert, Text: "2"},
				{Kind: Delete, Text: "3"},
				{Kind: Insert, Text: "4"},
				{Kind: Delete, Text: "5"},
				{Kind: Insert, Text: "6"},
				{Kind: Equal, Text: "7"},
			},
			context:  2,
			expected: []chunkRange{{start: 0, end: 8, isNewGroup: false}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateChunkRanges(tt.edits, tt.context)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("calculateChunkRanges(%q) = \n%#v\nwant\n%#v", tt.name, result, tt.expected)
			}
		})
	}
}
