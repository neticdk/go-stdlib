// Package simple provides a straightforward implementation of a diff algorithm
// for comparing strings or string slices.
//
// This package uses a simple longest common subsequence (LCS) based approach to
// identify differences between inputs. It's designed for clarity and ease of use
// rather than maximum performance on large inputs.
//
// # Basic Usage
//
// To compare two strings:
//
//     diff, err := simple.Diff("hello\nworld", "hello\nthere\nworld")
//
// To compare slices of strings directly:
//
//     diff, err := simple.DiffStrings([]string{"hello", "world"},
//                                   []string{"hello", "there", "world"})
//
// # Output Format
//
// The output is a formatted string showing the differences between inputs:
//   - Lines prefixed with spaces are unchanged (present in both inputs)
//   - Lines prefixed with "+" were added (present only in the second input)
//   - Lines prefixed with "-" were deleted (present only in the first input)
//
// By default, line numbers from both inputs are shown in the output.
//
// # Configuration Options
//
// The package supports the following option:
//
//     // Hide line numbers in the output
//     diff, err := simple.Diff(a, b, simple.WithShowLineNumbers(false))
//
// # Algorithm Characteristics
//
// The simple diff algorithm uses a greedy longest common subsequence approach.
// It's suitable for small to medium-sized inputs and situations where clarity
// is more important than performance optimization.

package simple
