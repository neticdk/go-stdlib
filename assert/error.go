package assert

import (
	"errors"
	"reflect"
)

// Error asserts that the given error is not nil.
func Error(t testingT, got error, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got == nil {
		reportError(t, ctx, "Expected error, got nil")
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NoError asserts that the given error is nil.
// Opposite of Error.
func NoError(t testingT, got error, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got != nil {
		reportError(t, ctx, "Expected no error, got: %v", got)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// ErrorIs asserts that the error 'got' matches the target error 'target'
// using errors.Is. It follows the chain of wrapped errors.
func ErrorIs(t testingT, got error, target error, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got == nil {
		reportError(t, ctx, "Expected error matching %v, but got nil error", target)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	if !errors.Is(got, target) {
		reportError(t, ctx, "Error is not the target type:\n Got error: %v\n Want target: %v", got, target)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// ErrorAs asserts that the error 'got' can be assigned to the type pointed to
// by 'targetPtr'using errors.As.
// It follows the chain of wrapped errors and assigns the matching error to
// *targetPtr if found.
// 'targetPtr' must be a non-nil pointer to either an interface type or a
// concrete type that implements error.
func ErrorAs(t testingT, got error, targetPtr any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got == nil {
		// Determine target type for better error message
		targetType := "unknown type"
		if targetPtr != nil {
			ptrType := reflect.TypeOf(targetPtr)
			if ptrType.Kind() == reflect.Ptr {
				targetType = ptrType.Elem().String()
			}
		}
		reportError(t, ctx, "Expected error of type %q, but got nil error", targetType)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	if targetPtr == nil {
		reportError(t, ctx, "Target pointer for ErrorAs cannot be nil")
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	ptrType := reflect.TypeOf(targetPtr)
	if ptrType.Kind() != reflect.Ptr {
		reportError(t, ctx, "Target for ErrorAs must be a pointer, got %T", targetPtr)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	// We don't strictly need to check if targetPtr points to an error type here,
	// as errors.As will panic if it doesn't, but it could be a pre-check.
	// Let errors.As handle the check for simplicity for now.

	if !errors.As(got, targetPtr) {
		reportError(t, ctx, "Error cannot be assigned to target type %T:\n Got error: %v", targetPtr, got)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}
