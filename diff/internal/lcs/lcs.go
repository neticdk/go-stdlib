package lcs

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
	var lcs []matchItem
	i, j := n, m
	for i > 0 && j > 0 {
		// If characters match, it's part of the LCS
		if a[i-1] == b[j-1] { //nolint:gocritic // ifElseChain is more readable here
			// Prepend the match (indices are i-1, j-1 because slices are 0-based)
			// Note: Appending and reversing later is more efficient, but
			//       prepending here matches the original output order need
			//       without an explicit reverse step. Given the scale where
			//       LCS is likely used (fallback/simple), this might be acceptable.
			//       For extreme performance, append then reverse.
			lcs = append([]matchItem{{AIndex: i - 1, BIndex: j - 1}}, lcs...)
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			// Move up (equivalent to deleting a[i-1])
			i--
		} else {
			// Move left (equivalent to inserting b[j-1])
			// Note: >= handles cases where lengths are equal - choice doesn't affect LCS length.
			j--
		}
	}

	return lcs
}
