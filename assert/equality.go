package assert

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

// Define epsilon constants based on precision requirements
const (
	float32Epsilon = 1e-6  // For float32
	float64Epsilon = 1e-15 // For float64
)

// Equal asserts that the given values are equal.
// It does not not create advanced diff output for complex types. Use something
// like github.com/google/go-cmp/cmp for more advanced diff output.
func Equal[T any](t testingT, got T, want T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !equalInternal(got, want) {
		var details []string
		if DiffsEnabled && shouldGenerateDiff(got, want) {
			diff := computeDiff(got, want)
			if diff != "" {
				// Add the diff with proper indentation
				diffLines := strings.Split(diff, "\n")
				if len(diffLines) > 0 {
					details = append(details, "Diff:")
					for _, line := range diffLines {
						details = append(details, "  "+line)
					}
				}
			}
		}
		reportEqualityError(t, ctx, "equal", got, want, details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NotEqual asserts that the given values are not equal.
// It does not not create advanced diff output for complex types. Use something
// like github.com/google/go-cmp/cmp for more advanced diff output.
// Opposite of Equal.
func NotEqual[T any](t testingT, got T, want T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if equalInternal(got, want) {
		reportEqualityError(t, ctx, "not equal", got, want)
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
		reportEqualityError(t, ctx, "within delta", got, want,
			fmt.Sprintf("Delta: %g", delta),
			fmt.Sprintf("Difference: %g", diff),
			fmt.Sprintf("Allowed difference: %g", absDelta))
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
		reportEqualityError(t, ctx, "not within delta", got, want,
			fmt.Sprintf("Delta: %g", delta),
			fmt.Sprintf("Difference: %.6g", diff),
			fmt.Sprintf("Allowed difference: %g", absDelta))
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

	// Special handling for time.Time using its Equal method, as DeepEqual
	// compares unexported fields like loc which can differ spuriously.
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

// reportEqualityError reports errors for equality assertions (Equal, NotEqual, etc.)
func reportEqualityError(t testingT, ctx *AssertionContext, assertion string, got, want any, details ...string) { //revive:disable-line:argument-limit
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	// Start building parts of the message
	parts := []string{
		fmt.Sprintf("Got: %#v", got),
		fmt.Sprintf("Want: %#v", want),
	}

	// Add any additional details
	if len(details) > 0 {
		parts = append(parts, "Details:")
		for _, detail := range details {
			parts = append(parts, "  "+detail) // Add indentation to details
		}
	}

	// Create a message based on the assertion type
	var messageHeading string
	switch assertion {
	case "equal":
		messageHeading = "Values are not equal"
	case "not equal":
		messageHeading = "Values should not be equal"
	case "within delta":
		messageHeading = "Values are not within delta"
	case "not within delta":
		messageHeading = "Values are within delta"
	case "same":
		messageHeading = "Values do not reference the same object"
	default:
		messageHeading = fmt.Sprintf("Equality assertion '%s' failed", assertion)
	}

	messageBody := strings.Join(parts, "\n  ") // Consistent indentation
	message := fmt.Sprintf("%s:\n  %s", messageHeading, messageBody)

	reportError(t, ctx, "%s", message)
}

// calculateEpsilon determines an appropriate epsilon based on the magnitude of values
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
