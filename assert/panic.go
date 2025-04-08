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
		err := &AssertionError{
			Message: "Expected panic, but code did not panic",
		}
		reportAssertionError(t, ctx, err)
		logOptionalMessage(t, msgAndArgs...)
	}

	return
}

// NotPanics asserts that the code inside the specified function doesn't panic.
func NotPanics(t testingT, f func(), msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	didPanic := false

	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			// Format panic value for clear reporting
			panicValueStr := fmt.Sprintf("[%T]: %v", r, r)

			err := &AssertionError{
				Message: "Unexpected panic",
				ExtraValues: []assertionValue{
					{Label: "Panic Value", Value: panicValueStr},
				},
			}

			// Always capture stack trace on unexpected panic
			const size = maxPanicStackDepth
			buf := make([]byte, size)
			n := runtime.Stack(buf, false)
			err.Stack = filterStackTrace(string(buf[:n]))

			reportAssertionError(t, ctx, err)

			logOptionalMessage(t, msgAndArgs...)
		}
	}()

	f()

	return !didPanic
}
