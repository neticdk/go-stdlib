package assert_test

import (
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestPanics(t *testing.T) {
	t.Run("panics", func(t *testing.T) {
		mockT := &mockTestingT{}
		assert.Panics(mockT, func() {
			panic("test panic")
		})
		if mockT.Failed() {
			t.Errorf("Test should have passed because a panic occurred, but reported failure")
		}
		if len(mockT.errorMessages) > 0 {
			t.Errorf("Expected no error messages, got: %v", mockT.errorMessages)
		}
	})

	t.Run("doesn't panic", func(t *testing.T) {
		mockT := &mockTestingT{}
		assert.Panics(mockT, func() {})
		if mockT.Failed() {
			if len(mockT.errorMessages) == 0 {
				t.Errorf("Expected error message, got none")
			} else if !strings.Contains(mockT.errorMessages[0], "Expected panic") {
				t.Errorf("Expected 'Expected panic, but code did not panic', got: %v", mockT.errorMessages[0])
			}
		} else {
			t.Errorf("Test should have failed because no panic occurred, but it passed")
		}
	})
}

func TestNotPanics(t *testing.T) {
	t.Run("doesn't panic", func(t *testing.T) {
		mockT := &mockTestingT{}
		assert.NotPanics(mockT, func() {})
		if mockT.Failed() {
			t.Errorf("Test should have passed because no panic occurred, but reported failure")
		}
		if len(mockT.errorMessages) > 0 {
			t.Errorf("Expected no error messages, got: %v", mockT.errorMessages)
		}
	})

	t.Run("panics", func(t *testing.T) {
		mockT := &mockTestingT{}
		assert.NotPanics(mockT, func() {
			panic("test panic")
		})

		if !mockT.Failed() {
			t.Errorf("Test should have failed because a panic occurred, but it passed")
		} else {
			if len(mockT.errorMessages) == 0 {
				t.Errorf("Expected error message, got none")
			} else {
				message := mockT.errorMessages[0]

				// Check the first line (message part)
				expectedMessagePart := "Unexpected panic"
				lines := strings.Split(message, "\n")
				if len(lines) < 1 || !strings.Contains(lines[0], expectedMessagePart) {
					t.Errorf("Expected message starting with %q, got: %q",
						expectedMessagePart, lines[0])
				}

				// Check for stack trace presence
				hasStackTrace := false
				for _, line := range lines {
					if strings.Contains(line, "Panic Value") {
						hasStackTrace = true
						break
					}
					if strings.Contains(line, "Stack Trace") {
						hasStackTrace = true
						break
					}
				}
				if !hasStackTrace {
					t.Errorf("Error message doesn't contain stack trace")
				}
			}
		}
	})
}
