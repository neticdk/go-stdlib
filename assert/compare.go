package assert

import (
	"cmp"
	"fmt"
	"reflect"
)

// Zero asserts that a value is the zero value for its type.
// Uses reflect.Value.IsZero(). Note this differs slightly from Empty for
// collections.e.g., Zero is false for an empty non-nil slice, but Empty is
// true.
func Zero(t testingT, data any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	isZero := false
	if data == nil {
		isZero = true // Interface containing nil is considered zero here
	} else {
		v := reflect.ValueOf(data)
		if v.IsValid() {
			isZero = v.IsZero()
		}
	}

	if !isZero {
		details := []string{fmt.Sprintf("Expected zero value, got zero for type %T: %#v", data, data)}
		reportEqualityError(t, ctx, data, "<zero value>", details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NotZero asserts that a value is not the zero value for its type.
// Opposite of Zero.
func NotZero(t testingT, data any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	isZero := false
	if data == nil {
		isZero = true
	} else {
		v := reflect.ValueOf(data)
		if v.IsValid() {
			isZero = v.IsZero()
		}
	}

	if isZero {
		details := []string{fmt.Sprintf("Expected non-zero value, got zero for type %T: %#v", data, data)}
		reportInequalityError(t, ctx, data, "<non-zero value>", details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// Greater asserts that 'got' is strictly greater than 'threshold'.
// Requires types compatible with cmp.Ordered (numbers, strings).
func Greater[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got <= threshold {
		reportComparisonError(t, ctx, "Value is not greater than threshold", got, threshold, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// GreaterOrEqual asserts that 'got' is greater than or equal to 'threshold'.
// Requires types compatible with cmp.Ordered.
func GreaterOrEqual[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got < threshold {
		reportComparisonError(t, ctx, "Value is not greater than or equal to threshold", got, threshold, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// Less asserts that 'got' is strictly less than 'threshold'.
// Requires types compatible with cmp.Ordered.
func Less[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got >= threshold {
		reportComparisonError(t, ctx, "Value is not less than threshold", got, threshold, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// LessOrEqual asserts that 'got' is less than or equal to 'threshold'.
// Requires types compatible with cmp.Ordered.
func LessOrEqual[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got > threshold {
		reportComparisonError(t, ctx, "Value is not less than or equal to threshold", got, threshold, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// Positive asserts that 'got' is greater than zero.
// Requires types compatible with cmp.Ordered.
func Positive[T cmp.Ordered](t testingT, got T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	var zero T // Get the zero value for the type T
	if got <= zero {
		err := &AssertionError{
			Message: "Value not positive",
			PrimaryValue: assertionValue{
				Label: "Got",
				Value: got,
			},
		}
		reportAssertionError(t, ctx, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// Negative asserts that 'got' is less than zero.
// Requires types compatible with cmp.Ordered.
func Negative[T cmp.Ordered](t testingT, got T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	var zero T       // Get the zero value for the type T
	if got >= zero { // Check if got >= 0
		err := &AssertionError{
			Message: "Value not negative",
			PrimaryValue: assertionValue{
				Label: "Got",
				Value: got,
			},
		}
		reportAssertionError(t, ctx, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// ElementsMatch asserts that the specified listA(array, slice...) is equal to
// specified listB(array, slice...) ignoring the order of the elements.
// If there are duplicate elements, the number of appearances of each of them in
// both lists should match.
//
// ElementsMatch(t, [1, 3, 2, 3], [1, 3, 3, 2])
func ElementsMatch[T comparable](t testingT, listA []T, listB []T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	listALen, listBLen := len(listA), len(listB)
	if listALen != listBLen {
		details := []string{fmt.Sprintf("length of listA (%d) does not match length of listB (%d)", listALen, listBLen)}
		reportEqualityError(t, ctx, listALen, listBLen, details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	countsA := make(map[T]int)
	countsB := make(map[T]int)

	for _, elem := range listA {
		countsA[elem]++
	}

	for _, elem := range listB {
		countsB[elem]++
	}

	if !reflect.DeepEqual(countsA, countsB) {
		details := []string{fmt.Sprintf("listA (%v) does not match listB (%v)", listA, listB)}
		reportEqualityError(t, ctx, listA, listB, details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	return true
}

// reportComparisonError reports errors for comparison assertions (Greater, Less, etc.)
// nolint: unparam
func reportComparisonError(t testingT, ctx *AssertionContext, message string, got, threshold any, err error, details ...string) { //revive:disable-line:argument-limit
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	assertErr := &AssertionError{
		Message: message,
		PrimaryValue: assertionValue{
			Label: "Got",
			Value: got,
		},
		ComparisonValue: assertionValue{
			Label: "Threshold",
			Value: threshold,
		},
	}

	if err != nil {
		assertErr.Error = err
	}

	if len(details) > 0 {
		assertErr.Details = details
	}

	reportAssertionError(t, ctx, assertErr)
}
