package assert

// True asserts that the given boolean value is true.
func True(t testingT, got bool, msgAndArgs ...any) bool { //revive:disable-line:flag-parameter
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !got {
		reportEqualityError(t, ctx, got, true)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// False asserts that the given boolean value is false.
func False(t testingT, got bool, msgAndArgs ...any) bool { //revive:disable-line:flag-parameter
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got {
		reportEqualityError(t, ctx, got, false)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}
