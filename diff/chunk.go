package diff

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

	// Handle empty input case
	if len(edits) == 0 {
		return []chunkRange{}
	}

	changeIndices := findChangeIndices(edits)
	if len(changeIndices) == 0 {
		if contextLines > 0 {
			return []chunkRange{{start: 0, end: len(edits), isNewGroup: false}}
		}
		return []chunkRange{}
	}

	// Handle context=0 case
	if contextLines == 0 {
		ranges := make([]chunkRange, 0, len(changeIndices))
		for _, idx := range changeIndices {
			ranges = append(ranges, chunkRange{start: idx, end: idx + 1, isNewGroup: true})
		}
		return ranges
	}

	// Group changes that are close based on contextLines
	groups := groupIndices(changeIndices, contextLines*2)

	// Create initial ranges without merging
	ranges := make([]chunkRange, 0, len(groups))
	for i, group := range groups {
		firstChange := group[0]
		lastChange := group[len(group)-1]

		startIdx := max(0, firstChange-contextLines)
		endIdx := min(len(edits), lastChange+contextLines+1)

		isNew := i > 0
		ranges = append(ranges, chunkRange{start: startIdx, end: endIdx, isNewGroup: isNew})
	}

	// Return early if no merging needed
	if len(ranges) <= 1 {
		return ranges
	}

	// Process ranges to merge only when necessary while preserving group structure
	result := make([]chunkRange, 0, len(ranges))
	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		next := ranges[i]

		// If ranges overlap significantly (more than just touching at edges)
		if next.start < current.end-contextLines {
			// Merge the ranges
			current.end = max(current.end, next.end)
		} else {
			// Add current range and start new one
			result = append(result, current)
			current = next
		}
	}
	// Add the final range
	result = append(result, current)

	// Post-process to ensure groups are properly marked
	for i := 1; i < len(result); i++ {
		result[i].isNewGroup = true
	}

	return result
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

// calculateStartLinesForIndex determines the original a and b line numbers
// corresponding to the start of the edit at the given targetIndex.
func calculateStartLinesForIndex(edits []Line, targetIndex int) (aLineNum, bLineNum int) {
	aLine, bLine := 1, 1 // Start from line 1
	for i := 0; i < targetIndex; i++ {
		// Ensure we don't go out of bounds if targetIndex is invalid
		if i >= len(edits) {
			break
		}
		switch edits[i].Kind {
		case Equal:
			aLine++
			bLine++
		case Delete:
			aLine++
		case Insert:
			bLine++
		}
	}
	return aLine, bLine
}
