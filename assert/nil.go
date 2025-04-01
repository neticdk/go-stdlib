package assert

import (
	"reflect"
)

// Nil asserts that the given value is nil.
func Nil(t testingT, data any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !isNilInternal(data) {
		reportError(t, ctx, "Expected nil, got: %#v", data)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NotNil asserts that the given value is not nil.
// Opposite of Nil.
func NotNil(t testingT, data any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if isNilInternal(data) {
		reportError(t, ctx, "Expected not nil, got nil")
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

func isNilInternal(data any) bool {
	if data == nil {
		return true
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}
