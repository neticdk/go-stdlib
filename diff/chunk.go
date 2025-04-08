package diff

import "sort"

type chunkRange struct {
	start      int  // Start index (inclusive) in the original edits slice
	end        int  // End index (exclusive) in the original edits slice
	isNewGroup bool // Indicates if this range starts a new group (for separators)
}

// calculateChunkRanges calculates the chunk ranges for the given edits and
// context lines.
func calculateChunkRanges(edits []Line, contextLines int) []chunkRange {
	if contextLines < 0 {
		contextLines = 0
	}

	// Handle empty edits case.
	if len(edits) == 0 {
		return []chunkRange{}
	}

	// Find the indices of the lines that are not equal (i.e., changes).
	changeIndices := findChangeIndices(edits)

	// If there are no changes and contextLines is 0, return an empty slice.  Otherwise, return a single chunk of the entire edit.
	if len(changeIndices) == 0 {
		if contextLines > 0 {
			return []chunkRange{{start: 0, end: len(edits), isNewGroup: false}}
		}
		return []chunkRange{}
	}

	// If context is 0, then simply generate a chunk for each diff.
	if contextLines == 0 {
		ranges := make([]chunkRange, 0, len(changeIndices))
		for _, idx := range changeIndices {
			ranges = append(ranges, chunkRange{start: idx, end: idx + 1, isNewGroup: true})
		}
		return ranges
	}

	// Group the change indices based on the context lines.
	groups := groupIndices(changeIndices, contextLines*2)

	// Create chunk ranges from the grouped change indices.
	ranges := createChunkRanges(edits, groups, contextLines)

	// Merge overlapping chunk ranges to form larger chunks.
	mergedRanges := mergeChunkRanges(ranges)

	// Post-process to ensure groups are properly marked
	for i := 1; i < len(mergedRanges); i++ {
		mergedRanges[i].isNewGroup = true
	}

	return mergedRanges
}

// createChunkRanges creates ChunkRange structs from grouped change indices.
func createChunkRanges(edits []Line, groups [][]int, contextLines int) []chunkRange {
	ranges := make([]chunkRange, 0, len(groups))
	for i, group := range groups {
		firstChange := group[0]
		lastChange := group[len(group)-1]

		startIdx := max(0, firstChange-contextLines)
		endIdx := min(len(edits), lastChange+contextLines+1)

		isNew := i > 0
		ranges = append(ranges, chunkRange{start: startIdx, end: endIdx, isNewGroup: isNew})
	}
	return ranges
}

// mergeChunkRanges merges overlapping ChunkRange structs.
func mergeChunkRanges(ranges []chunkRange) []chunkRange {
	if len(ranges) <= 1 {
		return ranges
	}

	// Sort ranges by start index to enable merging.
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	mergedRanges := make([]chunkRange, 0, len(ranges))
	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		next := ranges[i]

		// If ranges overlap significantly (more than just touching at edges)
		if next.start < current.end {
			// Merge the ranges
			current.end = max(current.end, next.end)
		} else {
			// Add current range and start new one
			mergedRanges = append(mergedRanges, current)
			current = next
		}
	}
	// Add the final range
	mergedRanges = append(mergedRanges, current)

	return mergedRanges
}

// findChangeIndices returns the indices of non-Equal edits.
func findChangeIndices(edits []Line) []int {
	indices := make([]int, 0)
	for i, edit := range edits {
		if edit.Kind != Equal {
			indices = append(indices, i)
		}
	}
	return indices
}

// groupIndices groups consecutive indices where the gap between them is <= maxGap.
func groupIndices(indices []int, maxGap int) [][]int {
	if len(indices) == 0 {
		return [][]int{}
	}

	groups := make([][]int, 0)
	currentGroup := []int{indices[0]}

	for i := 1; i < len(indices); i++ {
		if indices[i]-indices[i-1] <= maxGap {
			// Add to current group
			currentGroup = append(currentGroup, indices[i])
		} else {
			// Start a new group
			groups = append(groups, currentGroup)
			currentGroup = []int{indices[i]}
		}
	}
	// Add the last group
	groups = append(groups, currentGroup)

	return groups
}
