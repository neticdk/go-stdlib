package assert

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

// Define epsilon constants based on precision requirements
const (
	float32Epsilon = 1e-6  // For float32
	float64Epsilon = 1e-15 // For float64
)

// Equal asserts that the given values are equal.
// It doesn't create advanced diff output for complex types. Use something
// like github.com/google/go-cmp/cmp for more advanced diff output.
func Equal[T any](t testingT, got T, want T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !equalInternal(got, want) {
		err := &AssertionError{
			Message: "Values are not equal",
			PrimaryValue: assertionValue{
				Label: "Got",
				Value: got,
			},
			ComparisonValue: assertionValue{
				Label: "Want",
				Value: want,
			},
		}

		if DiffsEnabled && shouldGenerateDiff(got, want) {
			diff := computeDiff(got, want)
			err.Diff = diff
		}
		reportAssertionError(t, ctx, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NotEqual asserts that the given values are not equal.
// It doesn't create advanced diff output for complex types. Use something
// like github.com/google/go-cmp/cmp for more advanced diff output.
// Opposite of Equal.
func NotEqual[T any](t testingT, got T, want T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if equalInternal(got, want) {
		reportInequalityError(t, ctx, got, want)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// InDelta asserts that the given values are within delta.
func InDelta[T ~float32 | ~float64](t testingT, got T, want T, delta T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	diff, absDelta := deltaInternal(got, want, delta)

	// Calculate an appropriate epsilon based on the scale of the values
	epsilon := calculateEpsilon(got, want, diff, absDelta)

	if diff > (absDelta + epsilon) {
		err := &AssertionError{
			Message: "Values are not within delta",
			PrimaryValue: assertionValue{
				Label: "Got",
				Value: got,
			},
			ComparisonValue: assertionValue{
				Label: "Expected",
				Value: want,
			},
			ExtraValues: []assertionValue{
				{
					Label: "Delta",
					Value: delta,
				},
				{
					Label: "Difference",
					Value: diff,
				},
				{
					Label: "Allowed difference",
					Value: absDelta,
				},
			},
		}
		reportAssertionError(t, ctx, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	return true
}

// NotInDelta asserts that the given values are not within delta.
// Opposite of InDelta.
func NotInDelta[T ~float32 | ~float64](t testingT, got T, want T, delta T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	diff, absDelta := deltaInternal(got, want, delta)

	// Calculate an appropriate epsilon based on the scale of the values
	epsilon := calculateEpsilon(got, want, diff, absDelta)

	if diff <= (absDelta + epsilon) {
		err := &AssertionError{
			Message: "Values are within delta",
			PrimaryValue: assertionValue{
				Label: "Got",
				Value: got,
			},
			ComparisonValue: assertionValue{
				Label: "Expected",
				Value: want,
			},
			ExtraValues: []assertionValue{
				{
					Label: "Delta",
					Value: delta,
				},
				{
					Label: "Difference",
					Value: diff,
				},
				{
					Label: "Allowed difference",
					Value: absDelta,
				},
			},
		}
		reportAssertionError(t, ctx, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// deltaInternal returns the absolute difference between got and want, and the
// absolute value of delta
func deltaInternal[T ~float32 | ~float64](got, want, delta T) (diff, absDelta T) {
	diff = got - want
	if diff < 0 {
		diff = -diff
	}
	absDelta = delta
	if absDelta < 0 {
		absDelta = -absDelta
	}
	return
}

func equalInternal[T any](got T, want T) bool {
	// Use reflect.DeepEqual as the base, it handles many cases including nil.
	if reflect.DeepEqual(got, want) {
		return true
	}

	// Special handling for `time.Time` using its Equal method, as DeepEqual
	// compares unexported fields like `loc` which can differ spuriously.
	// We need to check AFTER DeepEqual fails, as two nil times are DeepEqual.
	vGot := reflect.ValueOf(got)
	vWant := reflect.ValueOf(want)

	if vGot.IsValid() && vWant.IsValid() && vGot.Type().AssignableTo(reflect.TypeOf(time.Time{})) && vWant.Type().AssignableTo(reflect.TypeOf(time.Time{})) {
		// Use Interface() to get the underlying value and assert type
		gotTime, gotOk := vGot.Interface().(time.Time)
		wantTime, wantOk := vWant.Interface().(time.Time)

		if gotOk && wantOk {
			return gotTime.Equal(wantTime)
		}
	}

	return false
}

// calculateEpsilon determines an appropriate epsilon based on the magnitude of values
// Epsilon is scaled to handle varying magnitudes of floats.
func calculateEpsilon[T ~float32 | ~float64](got, want, diff, absDelta T) T {
	// Default epsilon based on type
	var baseEpsilon T
	if reflect.TypeOf(got).Kind() == reflect.Float32 {
		baseEpsilon = T(float32Epsilon)
	} else {
		baseEpsilon = T(float64Epsilon)
	}

	// Scale epsilon based on magnitude of values and delta
	maxMagnitude := T(math.Max(math.Abs(float64(got)), math.Abs(float64(want))))
	if maxMagnitude > 1.0 {
		// For large values, scale epsilon relative to their magnitude
		return baseEpsilon * maxMagnitude
	}

	// For values near delta boundary, use a more conservative epsilon
	if math.Abs(float64(diff-absDelta)) < float64(baseEpsilon*100) {
		return baseEpsilon * 10
	}

	return baseEpsilon
}

// Helper function for equality assertions
func reportEqualityError(t testingT, ctx *AssertionContext, got, want any, details ...string) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	err := &AssertionError{
		Message: "Values are not equal",
		PrimaryValue: assertionValue{
			Label: "Got",
			Value: fmt.Sprintf("%#v", got),
		},
		ComparisonValue: assertionValue{
			Label: "Want",
			Value: fmt.Sprintf("%#v", want),
		},
		Details: details,
	}

	reportAssertionError(t, ctx, err)
}

// Helper function for inequality assertions
func reportInequalityError(t testingT, ctx *AssertionContext, got, want any, details ...string) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	err := &AssertionError{
		Message: "Values should not be equal",
		PrimaryValue: assertionValue{
			Label: "Got",
			Value: got,
		},
		ComparisonValue: assertionValue{
			Label: "Want",
			Value: want,
		},
		Details: details,
	}

	reportAssertionError(t, ctx, err)
}
