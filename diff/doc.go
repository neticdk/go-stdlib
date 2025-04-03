// Package diff provides implementations for computing the differences between
// strings or string slices.
//
// The package offers multiple diff algorithms, including a simple longest common
// subsequence (LCS) based approach and Myers' efficient diff algorithm. It also
// provides flexible options for customizing the output format.
//
// # Subpackages
//
// The package is organized into the following subpackages:
//
//   - simple: Provides a straightforward implementation of a diff algorithm
//     using a longest common subsequence (LCS) based approach. The `Diff` and
//     `DiffStrings` functions can be used directly, or a `Differ` instance can
//     be created for more complex scenarios.
//
//   - myers: Implements the Myers' diff algorithm, which is an efficient
//     method for computing the shortest edit script between two sequences. The
//     `Diff` and `DiffStrings` functions can be used directly, or a `Differ`
//     instance can be created.
//
// # Usage
//
// To use the diff package, you can choose either the simple or myers subpackage
// depending on your performance requirements. Both subpackages provide a
// `Differ` interface with `Diff` and `DiffStrings` methods for computing
// differences between strings and string slices, respectively. Alternatively,
// the `Diff` and `DiffStrings` functions in each subpackage can be used directly
// for simpler use cases.
//
// Example using the simple subpackage:
//
//		import (
//			"fmt"
//			"github.com/neticdk/go-stdlib/diff/simple"
//		)
//
//		differ := simple.NewDiffer()
//		diff := differ.Diff("hello\nworld", "hello\nthere\nworld")
//		fmt.Println(diff)
//
//	 	// Alternatively, use the Diff function directly
//	 	simpleDiff := simple.Diff("hello\nworld", "hello\nthere\nworld")
//
// Example using the myers subpackage:
//
//		import (
//			"fmt"
//			"github.com/neticdk/go-stdlib/diff/myers"
//		)
//
//		differ := myers.NewDiffer()
//		diff := differ.Diff("hello\nworld", "hello\nthere\nworld")
//		fmt.Println(diff)
//
//	 	// Alternatively, use the Diff function directly
//	 	myersDiff := myers.Diff("hello\nworld", "hello\nthere\nworld")
package diff
