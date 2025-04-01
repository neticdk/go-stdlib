package assert

import (
	"fmt"
	"runtime"
)

// Panics asserts that the code inside the specified function panics.
func Panics(t testingT, f func(), msgAndArgs ...any) (didPanic bool, panicValue any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			panicValue = r
		}
	}()

	f()

	if !didPanic {
		reportError(t, ctx, "Expected panic, but code did not panic")
		logOptionalMessage(t, msgAndArgs...)
	}

	return
}

// NotPanics asserts that the code inside the specified function does NOT panic.
func NotPanics(t testingT, f func(), msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	defer func() {
		if r := recover(); r != nil {
			// Get stack trace
			stackTrace := make([]byte, 4096)
			n := runtime.Stack(stackTrace, false)
			stackString := string(stackTrace[:n])
			//
			// Format error message with panic value and type
			panicType := fmt.Sprintf("%T", r)

			reportError(t, ctx, "Unexpected panic [%s]: %v\n\nStack trace:\n%s", panicType, r, stackString)
			logOptionalMessage(t, msgAndArgs...)
		}
	}()

	f()

	return true
}
