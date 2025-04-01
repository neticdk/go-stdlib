package assert

import (
	"reflect"
	"time"
)

// Equal asserts that the given values are equal.
// It does not not create advanced diff output for complex types. Use something
// like github.com/google/go-cmp/cmp for more advanced diff output.
func Equal[T any](t testingT, got T, want T, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if !equalInternal(got, want) {
		// Note to self: this is where we might want to add diff output
		t.Errorf("Values are not equal:\n Got: %#v\nWant: %#v", got, want)
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

	if equalInternal(got, want) {
		// Note to self: this is where we might want to add diff output
		t.Errorf("Values should not be equal:\n Got: %#v", got)
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

	diff, absDelta := deltaInternal(got, want, delta)
	if diff > absDelta {
		t.Errorf("Values are not within delta (%.6g):\n Got: %g\nWant: %g", delta, got, want)
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

	diff, absDelta := deltaInternal(got, want, delta)
	if diff <= absDelta {
		t.Errorf("Values are within delta (%.6g):\n Got: %g\nWant: %g", delta, got, want)
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
