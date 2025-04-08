// Package simple provides a straightforward implementation of a diff algorithm
// for comparing strings or string slices.
//
// This package uses a simple longest common subsequence (LCS) based approach to
// identify differences between inputs. It prioritizes clarity and ease of use
// over maximum performance, making it suitable for small to medium-sized inputs.
//
// # Basic Usage
//
// To compare two strings using the default settings:
//
//	differ := simple.NewDiffer()
//	diff, err := differ.Diff("hello\nworld", "hello\nthere\nworld")
//
// To compare slices of strings directly:
//
//	differ := simple.NewDiffer()
//	diff, err := differ.DiffStrings([]string{"hello", "world"}, []string{"hello", "there", "world"})
//
// You can also use the Diff or DiffStrings functions directly:
//
//	diff, err := simple.Diff("hello\nworld", "hello\nthere\nworld")
//	diff, err := simple.DiffStrings([]string{"hello", "world"}, []string{"hello", "there", "world"})
//
// # Output Format
//
// The output is a formatted string showing the differences between inputs.
// By default, it uses a context diff format with line numbers:
//   - Lines prefixed with spaces are unchanged (present in both inputs)
//   - Lines prefixed with "+" were added (present only in the second input)
//   - Lines prefixed with "-" were deleted (present only in the first input)
//
// It's possible to configure the output to use a unified diff, or to remove the line numbers.
// The format of the diff output is controlled by a `Formatter`.
//
// # Configuration Options
//
// The package supports the following options:
//
//   - `WithContextFormatter`:  Selects the context diff format (default).
//   - `WithUnifiedFormatter`:  Selects the unified diff format.
//   - `WithFormatter`: Allows specifying a custom `Formatter` implementation.
//
// Other options include:
//
// Create a custom differ that hides line numbers in the output:
//
//	differ := simple.NewCustomDiffer(simple.WithShowLineNumbers(false))
//
// Create a custom differ that uses the unified diff format:
//
//	differ := simple.NewCustomDiffer(simple.WithUnifiedFormatter())
//
// Create a custom formatter:
//
//	customFormatter := &MyCustomFormatter{}  // Replace with your custom implementation
//	differ := simple.NewCustomDiffer(simple.WithFormatter(customFormatter))
//
// You can also use the Diff or DiffStrings functions directly with options:
//
//	diff, err := simple.Diff("hello\nworld", "hello\nthere\nworld",
//	       simple.WithShowLineNumbers(false))
//	diff, err := simple.DiffStrings([]string{"hello", "world"}, []string{"hello", "there", "world"},
//	       simple.WithShowLineNumbers(false))
//
// # Algorithm Characteristics
//
// The simple diff algorithm uses a greedy longest common subsequence (LCS) approach
// implemented in the internal/diffcore package. It provides reasonable performance
// for small to medium-sized inputs where clarity is more important than
// optimal performance.
package simple
