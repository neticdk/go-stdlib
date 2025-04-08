// Package myers implements the Myers' diff algorithm.
//
// The Myers algorithm is an efficient method for computing the shortest edit
// script between two sequences (typically lines of text). This implementation
// provides both string-based and string slice-based diff functions with
// configuration options.
//
// For large inputs, a linear space variant of Myers algorithm is used to reduce
// memory consumption.  A fallback to a simpler LCS-based algorithm is also
// employed for large inputs (default > 10000 lines) or when recursion depth
// limits (default 30) are reached.
//
// # Basic Usage
//
// To compare two strings using the default settings:
//
//	differ := myers.NewDiffer()
//	diff, err := differ.Diff("hello\nworld", "hello\nthere\nworld")
//
// To compare slices of strings directly:
//
//	differ := myers.NewDiffer()
//	diff, err := differ.DiffStrings([]string{"hello", "world"}, []string{"hello", "there", "world"})
//
// You can also use the Diff or DiffStrings functions directly:
//
//	diff, err := myers.Diff("hello\nworld", "hello\nthere\nworld")
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
// By default, line numbers from both inputs are shown in the output. The output
// format defaults to a context diff, which includes surrounding lines to provide
// context for the changes. The specific formatting is handled by a `Formatter`
// interface.
//
// # Configuration
//
// Options can be used to customize the diff output and algorithm behavior. Key
// among these is the ability to select a `Formatter` to control the output
// style:
//
//   - `WithContextFormatter`:  Selects the context diff format (default).
//   - `WithUnifiedFormatter`:  Selects the unified diff format.
//   - `WithFormatter`: Allows specifying a custom `Formatter` implementation.
//
// Other options include:
//
// Create a custom differ with 5 context lines surrounding changes:
//
//	differ := myers.NewCustomDiffer(myers.WithContextLines(5))
//
// Create a custom differ that hides line numbers:
//
//	differ := myers.NewCustomDiffer(myers.WithShowLineNumbers(false))
//
// Create a custom differ that produces a unified diff output:
//
//	differ := myers.NewCustomDiffer(myers.WithUnifiedFormatter())
//
// Create a custom differ with a maximum edit distance (for performance on large inputs):
//
//	differ := myers.NewCustomDiffer(myers.WithMaxEditDistance(100))
//
// Create a custom differ that uses the linear space algorithm variant:
//
//	differ := myers.NewCustomDiffer(myers.WithLinearSpace(true))
//
// Combine options:
//
//	differ := myers.NewCustomDiffer(
//	     myers.WithContextLines(3),
//	     myers.WithShowLineNumbers(false))
//
// Create a custom formatter:
//
//	customFormatter := &MyCustomFormatter{}  // Replace with your custom implementation
//	differ := myers.NewCustomDiffer(myers.WithFormatter(customFormatter))
//
// You can also use the Diff or DiffStrings functions directly with options:
//
//	// With strings:
//	diff, err := myers.Diff("hello\nworld", "hello\nthere\nworld",
//		myers.WithContextLines(3), myers.WithShowLineNumbers(false))
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(diff)
//	}
//
//	// With string slices:
//	diff, err = myers.DiffStrings([]string{"hello", "world"}, []string{"hello", "there", "world"},
//		myers.WithContextLines(3), myers.WithShowLineNumbers(false))
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(diff)
//	}
//
// # Algorithm Details
//
// This implementation uses the following strategies:
//   - Myers' greedy algorithm for finding the shortest edit script.
//   - A linear space variant of Myers' algorithm (Hirschberg's algorithm
//     principle) to reduce memory usage for large inputs.
//   - A fallback to a simpler LCS-based algorithm (implemented in the
//     internal/diffcore package) for large edit distances, large
//     inputs, or when linear space recursion depth limits are reached.
//   - Context-aware output formatting.
//
// The time complexity is typically O(ND) where N is the sum of input lengths
// and D is the size of the minimum edit script, but the fallback to the
// LCS-based algorithm results in O(N*M) time complexity in certain cases.
//
// Space complexity varies depending on the chosen algorithm variant:
//   - O(N) for the linear space algorithm.
//   - O(V) for the standard Myers algorithm (V is the length of the `v` vector).
package myers
