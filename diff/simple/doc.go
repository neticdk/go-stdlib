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
// Alternatively, you can use the Diff or DiffStrings functions directly:
//
//	diff, err := simple.Diff("hello\nworld", "hello\nthere\nworld")
//	diff, err := simple.DiffStrings([]string{"hello", "world"}, []string{"hello", "there", "world"})
//
// # Output Format
//
// The output is a formatted string showing the differences between inputs:
//   - Lines prefixed with spaces are unchanged (present in both inputs)
//   - Lines prefixed with "+" were added (present only in the second input)
//   - Lines prefixed with "-" were deleted (present only in the first input)
//
// By default, line numbers are shown in the output.
//
// # Configuration Options
//
// The package supports the following option:
//
//	// Create a custom differ that hides line numbers in the output
//	differ := simple.NewCustomDiffer(simple.WithShowLineNumbers(false))
//
// Alternatively, you can use the Diff or DiffStrings functions directly with options:
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
