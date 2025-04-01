// Package assert provides a collection of assertion helpers designed to
// integrate seamlessly with the Go standard testing package. It aims to improve
// the readability and expressiveness of test code.
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
// It has been written as an alternative to the testify package but it is not a
// drop-in replacement:
//   - it supports most basic assertions
//   - it is type-safe
//   - it has far less features than testify
//   - it does not support advanced diffs - use github.com/google/go-cmp for that
//
// It draws inspiration from:
//   - testify - https://github.com/stretchr/testify
//   - is - https://github.com/matryer/is
//
// Assertions accept a testing.T interface, the value(s) being tested, and
// optional message arguments. They report failures using t.Errorf.
//
// The package also includes detailed error reporting with file and line information
// for easy debugging when assertions fail.
//
// A companion package 'require' provides the same assertions, but calls t.FailNow()
// to stop test execution immediately on failure.
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
package assert

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o README.md
