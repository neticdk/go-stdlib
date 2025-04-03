package diffcore

import (
	"github.com/neticdk/go-stdlib/diff"
)

// matchItem represents a matching item in two sequences
type matchItem struct {
	AIndex, BIndex int
}

// LongestCommonSubsequence finds the longest common subsequence between two string slices
// using a standard dynamic programming approach.
// Time complexity: O(len(a) * len(b))
// Space complexity: O(len(a) * len(b)) for the DP table.
func LongestCommonSubsequence(a, b []string) []matchItem {
	n := len(a)
	m := len(b)

	// Handle empty slice cases
	if n == 0 || m == 0 {
		return []matchItem{}
	}

	// dp[i][j] stores the length of the LCS of a[:i] and b[:j]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	// Build the DP table
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				// If elements match, extend the LCS from the diagonal
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				// If elements don't match, take the max LCS length
				// from excluding the current element of 'a' or 'b'
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// Backtrack from dp[n][m] to reconstruct the LCS match items
	// Allocate with estimated capacity (max length of LCS is min(n, m))
	lcs := make([]matchItem, 0, min(n, m))
	i, j := n, m
	for i > 0 && j > 0 {
		// If characters match, it's part of the LCS
		if a[i-1] == b[j-1] { //nolint:gocritic // ifElseChain is more readable here
			// Append the match (indices are i-1, j-1)
			lcs = append(lcs, matchItem{AIndex: i - 1, BIndex: j - 1})
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			// Move up (equivalent to deleting a[i-1])
			i--
		} else {
			// Move left (equivalent to inserting b[j-1])
			// Note: Tie-breaking choice (here favoring move left) doesn't affect LCS length.
			j--
		}
	}

	// Reverse the slice since we appended items in reverse order during backtracking.
	// This is generally more efficient than prepending in the loop.
	for k, l := 0, len(lcs)-1; k < l; k, l = k+1, l-1 {
		lcs[k], lcs[l] = lcs[l], lcs[k]
	}

	return lcs
}

// ComputeEditsLCS implements a simple diff algorithm based on longest common subsequence
func ComputeEditsLCS(a, b []string) []diff.Line {
	edits := []diff.Line{}

	// Use a simple longest common subsequence approach
	longest := LongestCommonSubsequence(a, b)
	aIndex, bIndex := 0, 0

	for _, item := range longest {
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
