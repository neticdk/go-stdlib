package diffcore

import "github.com/neticdk/go-stdlib/diff"

// GroupEditsByContext groups edits into chunks for context-aware display
func GroupEditsByContext(edits []diff.Line, contextLines int) [][]diff.Line {
	// If no context requested, only include non-equal edits
	if contextLines <= 0 {
		var changedEdits []diff.Line
		for _, edit := range edits {
			if edit.Kind != diff.Equal {
				changedEdits = append(changedEdits, edit)
			}
		}
		return [][]diff.Line{changedEdits}
	}

	// If the total edit is small enough, return as one chunk
	if len(edits) <= contextLines*2 {
		return [][]diff.Line{edits}
	}

	// Find indices of non-equal edits (changes)
	var changeIndices []int
	for i, edit := range edits {
		if edit.Kind != diff.Equal {
			changeIndices = append(changeIndices, i)
		}
	}

	// If no changes, return all as one chunk
	if len(changeIndices) == 0 {
		return [][]diff.Line{edits}
	}

	// Pre-allocate chunks - in worst case, each change could be in its own chunk
	// A more conservative estimate would be len(changeIndices)/2 + 1
	chunks := make([][]diff.Line, 0, min(len(changeIndices), len(edits)/2+1))

	// Group changes that are close to each other
	var currentGroup []int
	var groups [][]int

	for i, idx := range changeIndices {
		if i == 0 || idx-changeIndices[i-1] > contextLines*2 {
			// Start a new group
			if len(currentGroup) > 0 {
				groups = append(groups, currentGroup)
			}
			currentGroup = []int{idx}
		} else {
			// Add to current group
			currentGroup = append(currentGroup, idx)
		}
	}

	// Add the last group
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	// Create chunks based on groups with context
	for _, group := range groups {
		firstChange := group[0]
		lastChange := group[len(group)-1]

		// Calculate chunk boundaries with context
		startIdx := max(0, firstChange-contextLines)
		endIdx := min(len(edits), lastChange+contextLines+1)

		chunk := make([]diff.Line, endIdx-startIdx)
		copy(chunk, edits[startIdx:endIdx])
		chunks = append(chunks, chunk)
	}

	return chunks
}
