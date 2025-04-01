package assert

import "strings"

type testingT interface {
	Log(args ...any)
	Logf(format string, args ...any)
	Errorf(format string, args ...any)
}

type tHelper = interface {
	Helper()
}

// logOptionalMessage logs the optional message and arguments if provided.
// func logOptionalMessage(t *testing.T, msgAndArgs ...any) {
func logOptionalMessage(t testingT, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if len(msgAndArgs) > 0 {
		format, ok := msgAndArgs[0].(string)
		if ok {
			// Check if the format string contains format verbs
			if containsVerbs(format) && len(msgAndArgs) > 1 {
				t.Logf(format, msgAndArgs[1:]...)
			} else {
				t.Log(msgAndArgs...)
			}
		} else {
			t.Log(msgAndArgs...)
		}
	}
}

// containsVerbs checks if a format string contains format verbs.
// Rudementary implementation.
func containsVerbs(format string) bool {
	return strings.ContainsRune(format, '%')
}
