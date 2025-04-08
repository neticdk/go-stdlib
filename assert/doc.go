// Package assert provides a collection of assertion helpers designed to
// integrate seamlessly with the Go standard testing package. It aims to improve
// the readability and expressiveness of test code.
//
// It has been written as an alternative to the testify package but it's not a
// drop-in replacement:
//   - it supports most basic assertions
//   - it's type-safe
//   - it has far less features than testify
//   - it does not support advanced diffs - use github.com/google/go-cmp for that
//
// Key features include:
//   - Type-safe comparisons using generics (e.g., assert.Greater, assert.Equal).
//   - Checks for nil, errors, boolean conditions.
//   - Collection assertions (Contains, Empty, Len, ElementsMatch).
//   - Panic detection (Panics, NotPanics).
//   - Time assertions (TimeEqual, TimeBefore, TimeAfter, WithinDuration).
//   - Floating-point comparisons with delta for handling precision issues.
//   - Support for custom failure messages with optional formatting.
//   - Integration with t.Helper() for cleaner test failure reporting.
//   - Consistent error formatting with file and line information.
//   - Optional stack traces on failure for debugging complex call paths.
//
// It draws inspiration from:
//   - testify - https://github.com/stretchr/testify
//   - is - https://github.com/matryer/is
//
// # Examples
//
// Basic usage:
//
//	func TestSomething(t *testing.T) {
//	    result := Calculate()
//	    assert.Equal(t, result, 42)
//	    assert.NoError(t, err)
//	}
//
// Custom Error Messages:
//
//	assert.Equal(t, got, want, "calculation with input %d failed", input)
//
// Collection Assertions:
//
//	values := []int{1, 2, 3}
//	assert.Contains(t, values, 2)
//	assert.Len(t, values, 3)
//
// Error Handling:
//
//	err := doSomething()
//	if assert.Error(t, err) {
//		assert.ErrorIs(t, err, ErrNotFound)
//	}
//
// Time Assertions:
//
//	now := time.Now()
//	later := now.Add(time.Hour)
//	assert.TimeAfter(t, later, now)
//	assert.WithinDuration(t, eventTime, expectedTime, 5 * time.Second)
//
// Floating-Point Comparisons:
//
//	assert.InDelta(t, calculatedValue, expectedValue, 0.001)
//
// # Diffs
//
// The `Equal` function supports diffs when the `DiffsEnabled` global variable
// is set to true (default). This feature can be useful for debugging and
// understanding the differences between expected and actual values. That said,
// the implementation is simple and does not support complex types or nested
// structures. It's based on JSON serialization and comparison and thus might
// fail for unexported fields, channels or functions. Diff output uses the Myers
// diff algorithm.
//
// To disable diffs, set the `DiffsEnabled` global variable to false.
//
// # Stack Traces
//
// On assertion failure, an optional stack trace can be included in the error
// output. This is controlled by the `StackTracesEnabled` global variable, which
// defaults to `false`.
//
// Capturing stack traces incurs a performance cost due to the use of `runtime.Stack`.
// It's recommended to enable this feature primarily during debugging sessions
// when the context of the call path is needed to understand a complex failure.
//
// Example: Enable stack traces (e.g., in TestMain or a specific test):
//
//	assert.StackTracesEnabled = true
//	assert.Equal(t, someComplexResult(), expectedValue) // Stack trace printed if this fails
//	assert.StackTracesEnabled = false // Optionally disable afterwards
//
// **Special Case: `NotPanics`**
// Because unexpected panics are critical failures where the stack trace is essential
// for debugging, `assert.NotPanics` (and `require.NotPanics`) will **always**
// include a filtered stack trace in the error output when a panic occurs,
// regardless of the `StackTracesEnabled` setting.
//
// The generated stack trace is filtered to remove internal frames from the Go
// runtime, testing framework, and the assertion library itself, focusing on
// the user's test code.
//
// # Assertions
//
// All assertions accept a `testing.T` interface, the value(s) being tested, and
// optional message arguments. They report failures using t.Errorf.
//
// # Error reporting
//
// The package leverages `t.Helper()` and includes detailed error reporting
// with file and line information for debugging when assertions fail. Failed
// assertions will include details about the values compared, specific error
// messages, optional diffs, and optional stack traces.
//
// # Halting on failure
//
// A companion package 'require' provides the same assertions, but calls t.FailNow()
// to stop test execution immediately on failure.
package assert
