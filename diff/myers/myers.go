package myers

import (
	"fmt"
	"strings"

	"github.com/neticdk/go-stdlib/diff"
)

var (
	DefaultContextLines    = 3
	DefaultLinearSpace     = false
	DefaultShowLineNumbers = true
	DefaultMaxEditDistance = 10
)

// options configures the myers diff algorithm
type options struct {
	maxEditDistance int  // MaxEditDistance specifies the maximum edit distance to consider.
	linearSpace     bool // LinearSpace specifies whether to use linear space algorithm.

	// Formatting options
	contextLines    int
	showLineNumbers bool
}

// Diff computes differences between two values using the Myers diff algorithm.
func Diff(a, b string, opts ...Option) (string, error) {
	if a == "" && b == "" {
		return "", nil
	}

	aLines := splitLines(a)
	bLines := splitLines(b)
	return DiffStrings(aLines, bLines, opts...)
}

// DiffStrings computes differences between string slices using the Myers diff algorithm.
func DiffStrings(a, b []string, opts ...Option) (string, error) {
	return myersDiffStrings(a, b, applyOptions(opts...))
}

func myersDiffStrings(a, b []string, opts options) (string, error) {
	// Compute edit script using Myers' algorithm
	script := computeEditScript(a, b, opts.maxEditDistance)

	// Format the diff output
	var sb strings.Builder

	// Track line numbers if enabled
	aLineNum := 1
	bLineNum := 1

	// Group edits by type for context-aware output
	chunks := groupEditsByContext(script, opts.contextLines)

	for i, chunk := range chunks {
		// Add separator between chunks
		if i > 0 && opts.contextLines > 0 {
			sb.WriteString("...\n")
		}

		for _, edit := range chunk {
			switch edit.Kind {
			case diff.Equal:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d %4d   ", aLineNum, bLineNum))
				} else {
					sb.WriteString("  ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				aLineNum++
				bLineNum++
			case diff.Delete:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("%4d      - ", aLineNum))
				} else {
					sb.WriteString("- ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				aLineNum++
			case diff.Insert:
				if opts.showLineNumbers {
					sb.WriteString(fmt.Sprintf("     %4d + ", bLineNum))
				} else {
					sb.WriteString("+ ")
				}
				sb.WriteString(edit.Text)
				sb.WriteString("\n")
				bLineNum++
			}
		}
	}

	return sb.String(), nil
}

// computeEditScript implements Myers' diff algorithm to find the shortest edit script
func computeEditScript(a, b []string, maxEditDistance int) []diff.Line {
	n, m := len(a), len(b)
	maxDist := n + m

	// Handle special cases
	if maxDist == 0 {
		return []diff.Line{}
	}

	if n == 0 {
		result := make([]diff.Line, m)
		for i := range m {
			result[i] = diff.Line{Kind: diff.Insert, Text: b[i]}
		}
		return result
	}

	if m == 0 {
		result := make([]diff.Line, n)
		for i := range n {
			result[i] = diff.Line{Kind: diff.Delete, Text: a[i]}
		}
		return result
	}

	// Limit max edit distance if specified
	if maxEditDistance > 0 && maxEditDistance < maxDist {
		maxDist = maxEditDistance
	}

	// Initialize the edit graph
	v := make([]int, 2*maxDist+1)
	vs := make([][]int, 0, maxDist)

	// Find the shortest edit script path
	var x, y int
	var foundPath bool

	for d := 0; d <= maxDist; d++ {
		// Clone the current v for backtracking
		vClone := make([]int, len(v))
		copy(vClone, v)
		vs = append(vs, vClone)

		for k := -d; k <= d; k += 2 {
			// Determine whether to move down or right
			if k == -d || (k != d && v[maxDist+k-1] < v[maxDist+k+1]) {
				x = v[maxDist+k+1] // Move down
			} else {
				x = v[maxDist+k-1] + 1 // Move right
			}

			y = x - k

			// Follow diagonal (matching) paths as far as possible
			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}

			v[maxDist+k] = x

			// Check if we've reached the bottom-right corner
			if x >= n && y >= m {
				foundPath = true
				break
			}
		}

		if foundPath {
			break
		}
	}

	// If no path found within constraints, fall back to a simpler approach
	if !foundPath {
		return simpleDiff(a, b)
	}

	// Backtrack to construct the edit script
	return backtrack(a, b, vs, maxDist)
}

// backtrack constructs the edit script by walking backward through the edit graph
func backtrack(a, b []string, vs [][]int, midpoint int) []diff.Line {
	edits := []diff.Line{}
	x, y := len(a), len(b)

	for d := len(vs) - 1; d >= 0; d-- {
		v := vs[d]
		k := x - y

		var prevK, prevX, prevY int

		if k == -d || (k != d && v[midpoint+k-1] < v[midpoint+k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = v[midpoint+prevK]
		prevY = prevX - prevK

		// Handle diagonal moves (matching lines)
		for x > prevX && y > prevY {
			edits = append([]diff.Line{{Kind: diff.Equal, Text: a[x-1]}}, edits...)
			x--
			y--
		}

		// Handle vertical/horizontal moves (insertions/deletions)
		if d > 0 {
			if x > prevX {
				edits = append([]diff.Line{{Kind: diff.Delete, Text: a[x-1]}}, edits...)
				x--
			} else if y > prevY {
				edits = append([]diff.Line{{Kind: diff.Insert, Text: b[y-1]}}, edits...)
				y--
			}
		}
	}

	return edits
}

// simpleDiff is a fallback diff algorithm for when Myers becomes too expensive
func simpleDiff(a, b []string) []diff.Line {
	edits := []diff.Line{}

	// Use a simple longest common subsequence approach
	lcs := longestCommonSubsequence(a, b)
	aIndex, bIndex := 0, 0

	for _, item := range lcs {
		// Add deletions for unmatched items in a
		for aIndex < item.AIndex {
			edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
			aIndex++
		}

		// Add insertions for unmatched items in b
		for bIndex < item.BIndex {
			edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
			bIndex++
		}

		// Add the matching item
		edits = append(edits, diff.Line{Kind: diff.Equal, Text: a[aIndex]})
		aIndex++
		bIndex++
	}

	// Handle remaining items in a
	for aIndex < len(a) {
		edits = append(edits, diff.Line{Kind: diff.Delete, Text: a[aIndex]})
		aIndex++
	}

	// Handle remaining items in b
	for bIndex < len(b) {
		edits = append(edits, diff.Line{Kind: diff.Insert, Text: b[bIndex]})
		bIndex++
	}

	return edits
}

// MatchItem represents a matching item in two sequences
type MatchItem struct {
	AIndex, BIndex int
}

// longestCommonSubsequence finds the longest common subsequence between two string slices
func longestCommonSubsequence(a, b []string) []MatchItem {
	// Create a map of values to positions in B for faster lookup
	bValues := make(map[string][]int)
	for i, val := range b {
		bValues[val] = append(bValues[val], i)
	}

	var lcs []MatchItem

	// Find matches
	for i, aVal := range a {
		if positions, ok := bValues[aVal]; ok {
			// Try to find the best match position
			bestMatch := -1

			for _, pos := range positions {
				// Check if this match extends the current LCS
				valid := true
				for j := len(lcs) - 1; j >= 0; j-- {
					if lcs[j].BIndex >= pos {
						valid = false
						break
					}
				}

				if valid && (bestMatch == -1 || pos < bestMatch) {
					bestMatch = pos
				}
			}

			if bestMatch != -1 {
				lcs = append(lcs, MatchItem{i, bestMatch})
			}
		}
	}

	return lcs
}

// groupEditsByContext groups edits into chunks for context-aware display
func groupEditsByContext(edits []diff.Line, contextLines int) [][]diff.Line {
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

// splitLines splits a string into lines, handling empty strings and trailing newlines
func splitLines(s string) []string {
	// Special case for completely empty string
	if s == "" {
		return []string{}
	}

	// Split the string by newlines
	lines := strings.Split(s, "\n")

	// If the string ends with a newline, the split will produce an empty string
	// at the end - we should remove it to avoid confusing diff output
	if s[len(s)-1] == '\n' {
		lines = lines[:len(lines)-1]
	}

	return lines
}
