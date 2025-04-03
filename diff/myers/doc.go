// Package myers implements the Myers' diff algorithm.
//
// The Myers algorithm is an efficient method for computing the shortest edit
// script between/ two sequences (typically lines of text). This implementation
// provides both string-based/ and string slice-based diff functions with
// various configuration options.
//
// # Basic Usage
//
// To compare two strings:
//
//	diff, err := myers.Diff("hello\nworld", "hello\nthere\nworld")
//
// To compare slices of strings directly:
//
//	diff, err := myers.DiffStrings([]string{"hello", "world"}, []string{"hello", "there", "world"})
//
// # Output Format
//
// The output is a formatted string showing the differences between inputs with line numbers
// and change indicators:
//   - Lines prefixed with spaces are unchanged (present in both inputs)
//   - Lines prefixed with "+" were added (present only in the second input)
//   - Lines prefixed with "-" were deleted (present only in the first input)
//
// By default, line numbers from both inputs are shown in the output.
//
// # Configuration Options
//
// Several options can be used to customize the diff output:
//
//	// Show 5 context lines surrounding changes
//	diff, err := myers.Diff(a, b, myers.WithContextLines(5))
//
//	// Hide line numbers
//	diff, err := myers.Diff(a, b, myers.WithShowLineNumbers(false))
//
//	// Set maximum edit distance (for performance on large inputs)
//	diff, err := myers.Diff(a, b, myers.WithMaxEditDistance(100))
//
//	// Use linear space algorithm variant
//	diff, err := myers.Diff(a, b, myers.WithLinearSpace(true))
//
//	// Combine multiple options
//	diff, err := myers.Diff(a, b,
//	            myers.WithContextLines(3),
//	            myers.WithShowLineNumbers(false))
//
// # Algorithm Details
//
// This implementation uses:
// 1. Myers' greedy algorithm for finding the shortest edit script
// 2. A fallback to a simpler LCS-based algorithm for very large edit distances
// 3. Context-aware output formatting similar to unified diff format
//
// The time complexity is O(ND) where N is the sum of input lengths and D is the
// size of/ the minimum edit script. Space complexity is O(N) or O(D^2)
// depending on the chosen algorithm variant.

package myers
