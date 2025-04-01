package require

import (
	"cmp"
	"time"

	"github.com/neticdk/go-stdlib/assert"
)

type testingT interface {
	Log(args ...any)
	Logf(format string, args ...any)
	Errorf(format string, args ...any)
	FailNow()
}

type tHelper = interface {
	Helper()
}

// --- bool ---

// True requires that the specified value is true.
func True(t testingT, value bool, msgAndArgs ...any) { //revive:disable-line:flag-parameter
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.True(t, value, msgAndArgs...) {
		t.FailNow()
	}
}

// False requires that the specified value is false.
func False(t testingT, value bool, msgAndArgs ...any) { //revive:disable-line:flag-parameter
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.False(t, value, msgAndArgs...) {
		t.FailNow()
	}
}

// --- collection ---

// Contains requires that the specified list(array, slice...) contains the specified element.
func Contains(t testingT, collection any, element any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Contains(t, collection, element, msgAndArgs...) {
		t.FailNow()
	}
}

// NotContains requires that the specified list(array, slice...) does not contain the specified element.
func NotContains(t testingT, collection any, element any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NotContains(t, collection, element, msgAndArgs...) {
		t.FailNow()
	}
}

// ContainsKey requires that the specified map contains the specified key.
func ContainsKey[K comparable, V any](t testingT, m map[K]V, key K, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.ContainsKey(t, m, key, msgAndArgs...) {
		t.FailNow()
	}
}

// NotContainsKey requires that the specified map does not contain the specified key.
func NotContainsKey[K comparable, V any](t testingT, m map[K]V, key K, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NotContainsKey(t, m, key, msgAndArgs...) {
		t.FailNow()
	}
}

// Empty requires that the specified object is empty. Objects considered empty are
// nil, "", false, 0, or a channel, slice, array or map with len == 0.
func Empty(t testingT, data any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Empty(t, data, msgAndArgs...) {
		t.FailNow()
	}
}

// NotEmpty requires that the specified object is not empty. Objects considered empty are
// nil, "", false, 0, or a channel, slice, array or map with len == 0.
func NotEmpty(t testingT, data any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NotEmpty(t, data, msgAndArgs...) {
		t.FailNow()
	}
}

// Len requires that the specified object has the specified length. Len also fails if the object has a type that len() not accept.
func Len(t testingT, data any, expectedLen int, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Len(t, data, expectedLen, msgAndArgs...) {
		t.FailNow()
	}
}

// --- comparison ---

// Zero requires that the specified value is the zero value for its type.
func Zero(t testingT, data any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Zero(t, data, msgAndArgs...) {
		t.FailNow()
	}
}

// NotZero requires that the specified value is not the zero value for its type.
func NotZero(t testingT, data any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NotZero(t, data, msgAndArgs...) {
		t.FailNow()
	}
}

// Greater requires that the first element is greater than the second.
func Greater[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Greater(t, got, threshold, msgAndArgs...) {
		t.FailNow()
	}
}

// GreaterOrEqual requires that the first element is greater than or equal to the second.
func GreaterOrEqual[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.GreaterOrEqual(t, got, threshold, msgAndArgs...) {
		t.FailNow()
	}
}

// Less requires that the first element is less than the second.
func Less[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Less(t, got, threshold, msgAndArgs...) {
		t.FailNow()
	}
}

// LessOrEqual requires that the first element is less than or equal to the second.
func LessOrEqual[T cmp.Ordered](t testingT, got T, threshold T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.LessOrEqual(t, got, threshold, msgAndArgs...) {
		t.FailNow()
	}
}

// Positive requires that the specified value is positive.
func Positive[T cmp.Ordered](t testingT, got T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Positive(t, got, msgAndArgs...) {
		t.FailNow()
	}
}

// Negative requires that the specified value is negative.
func Negative[T cmp.Ordered](t testingT, got T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Negative(t, got, msgAndArgs...) {
		t.FailNow()
	}
}

// ElementsMatch requires that the specified lists or maps contain the same elements in any order.
func ElementsMatch[T comparable](t testingT, listA []T, listB []T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.ElementsMatch(t, listA, listB, msgAndArgs...) {
		t.FailNow()
	}
}

// --- equality ---

// Equal requires that the two objects are equal.
func Equal[T any](t testingT, got T, want T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Equal(t, got, want, msgAndArgs...) {
		t.FailNow()
	}
}

// NotEqual requires that the specified values are not equal.
func NotEqual[T any](t testingT, got T, want T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NotEqual(t, got, want, msgAndArgs...) {
		t.FailNow()
	}
}

// InDelta requires that the specified values are within the delta of each other.
func InDelta[T ~float32 | ~float64](t testingT, got T, want T, delta T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.InDelta(t, got, want, delta, msgAndArgs...) {
		t.FailNow()
	}
}

// NotInDelta requires that the specified values are not within the delta of each other.
// It fails if |got - want| <= |delta|.
func NotInDelta[T ~float32 | ~float64](t testingT, got T, want T, delta T, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	// Call the corresponding 'is' function
	if !assert.NotInDelta(t, got, want, delta, msgAndArgs...) {
		t.FailNow() // Fail immediately if assert.NotInDelta returns false
	}
}

// --- error ---

// Error requires that the specified err is not nil.
func Error(t testingT, got error, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Error(t, got, msgAndArgs...) {
		t.FailNow()
	}
}

// NoError requires that the specified err is nil.
func NoError(t testingT, got error, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NoError(t, got, msgAndArgs...) {
		t.FailNow()
	}
}

// ErrorIs requires that the specified error is assignable to the target.
func ErrorIs(t testingT, got error, target error, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.ErrorIs(t, got, target, msgAndArgs...) {
		t.FailNow()
	}
}

// ErrorAs requires that the specified error is assignable to the target.
func ErrorAs(t testingT, got error, target any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.ErrorAs(t, got, target, msgAndArgs...) {
		t.FailNow()
	}
}

// --- nil ---

// Nil requires that the specified object is nil.
func Nil(t testingT, data any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.Nil(t, data, msgAndArgs...) {
		t.FailNow()
	}
}

// NotNil requires that the specified object is not nil.
func NotNil(t testingT, data any, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NotNil(t, data, msgAndArgs...) {
		t.FailNow()
	}
}

// --- panic ---

// Panics requires that the code inside the specified PanicTestFunc panics.
func Panics(t testingT, f func(), msgAndArgs ...any) (didPanic bool, panicValue any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	didPanic, panicValue = assert.Panics(t, f, msgAndArgs...)
	if !didPanic {
		t.FailNow()
	}
	return
}

// NotPanics requires that the code inside the specified PanicTestFunc does not panic.
func NotPanics(t testingT, f func(), msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.NotPanics(t, f, msgAndArgs...) {
		t.FailNow()
	}
}

// TimeAfter requires that a time is after a threshold time
func TimeAfter(t testingT, got, threshold time.Time, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.TimeAfter(t, got, threshold, msgAndArgs...) {
		t.FailNow()
	}
}

// TimeBefore requires that a time is before a threshold time
func TimeBefore(t testingT, got, threshold time.Time, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.TimeBefore(t, got, threshold, msgAndArgs...) {
		t.FailNow()
	}
}

// TimeEqual requires that two times represent the same instant
func TimeEqual(t testingT, got, want time.Time, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.TimeEqual(t, got, want, msgAndArgs...) {
		t.FailNow()
	}
}

// WithinDuration requires that two times are within a certain duration of each other
func WithinDuration(t testingT, got, want time.Time, delta time.Duration, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.WithinDuration(t, got, want, delta, msgAndArgs...) {
		t.FailNow()
	}
}

// TimeEqualWithPrecision requires that two times are equal within a certain precision
func TimeEqualWithPrecision(t testingT, got, want time.Time, precision time.Duration, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.TimeEqualWithPrecision(t, got, want, precision, msgAndArgs...) {
		t.FailNow()
	}
}

// WithinTime requires that a time is within a given time window
func WithinTime(t testingT, got time.Time, start, end time.Time, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if !assert.WithinTime(t, got, start, end, msgAndArgs...) {
		t.FailNow()
	}
}
