// Package diff provides implementations for computing the differences between
// strings or string slices.
//
// The package offers two diff algorithms: a simple longest common subsequence
// (LCS) based approach and Myers' efficient diff algorithm. It also provides
// flexible options for customizing the output format, including context and
// unified diff formats, through the use of formatters.
//
// # Subpackages
//
// The package is organized into the following subpackages:
//
//   - simple: Provides a straightforward implementation of a diff algorithm
//     using a longest common subsequence (LCS) based approach. The `Diff` and
//     `DiffStrings` functions can be used directly, or a `Differ` instance can
//     be created for more complex scenarios. This implementation prioritizes simplicity
//     and readability, suitable for small to medium-sized inputs.
//
//   - myers: Implements the Myers' diff algorithm, which is an efficient
//     method for computing the shortest edit script between two sequences. The
//     `Diff` and `DiffStrings` functions can be used directly, or a `Differ`
//     instance can be created. The Myers implementation includes a linear space
//     optimization for large inputs and a fallback to the LCS algorithm for
//     large (> 10000 lines) inputs or when recursion limits are reached.
//
// # Usage
//
// To use the diff package, you can choose either the simple or myers subpackage
// depending on your performance requirements. Both subpackages provide a
// `Differ` interface with `Diff` and `DiffStrings` methods for computing
// differences between strings and string slices, respectively. The `Diff` and
// `DiffStrings` functions in each subpackage can also be used directly for
// simpler use cases. The format of the diff output is controlled by a
// `Formatter`. The default formatter produces a context diff with line numbers.
//
// Example using the simple subpackage:
//
//	import (
//	    "fmt"
//	    "github.com/neticdk/go-stdlib/diff/simple"
//	)
//
//	differ := simple.NewDiffer()
//	diff, err := differ.Diff("hello\nworld", "hello\nthere\nworld")
//	fmt.Println(diff)
//
// Use the Diff function directly:
//
//	simpleDiff, err := simple.Diff("hello\nworld", "hello\nthere\nworld")
//
// Example using the myers subpackage:
//
//	import (
//	      "fmt"
//	      "github.com/neticdk/go-stdlib/diff/myers"
//	)
//
//	differ := myers.NewDiffer()
//	diff, err := differ.Diff("hello\nworld", "hello\nthere\nworld")
//	fmt.Println(diff)
//
// Use the Diff function directly:
//
//	myersDiff, err := myers.Diff("hello\nworld", "hello\nthere\nworld")
package diff
