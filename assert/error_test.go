package assert_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

var (
	errBase    = errors.New("base error")
	errWrapped = fmt.Errorf("wrapped: %w", errBase)
)

var (
	errCustom        = &customError{msg: "custom error"}
	errCustomWrapped = fmt.Errorf("wrapped custom: %w", errCustom)
)

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantPass bool
	}{
		{"nil error", nil, false},
		{"non-nil error", errors.New("test error"), true},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := assert.Error(mockT, tt.err)
			if pass != tt.wantPass {
				t.Errorf("Error() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestNoError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantPass bool
	}{
		{"nil error", nil, true},
		{"non-nil error", errors.New("test error"), false},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := assert.NoError(mockT, tt.err)
			if pass != tt.wantPass {
				t.Errorf("NoError() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		target   error
		wantPass bool
	}{
		{"nil error, non-nil target", nil, errBase, false},
		{"non-nil error, nil target", errBase, nil, false},
		{"nil error, nil target", nil, nil, false},
		{"direct match", errBase, errBase, true},
		{"wrapped match", errWrapped, errBase, true},
		{"no match", errBase, io.EOF, false},
		{"wrapped no match", errWrapped, io.EOF, false},
		{"different wrapped", errWrapped, errCustom, false},
		{"stdlib error match", io.EOF, io.EOF, true},
		{"stdlib wrapped match", fmt.Errorf("wrapped EOF: %w", io.EOF), io.EOF, true},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := assert.ErrorIs(mockT, tt.err, tt.target)
			if pass != tt.wantPass {
				t.Errorf("ErrorIs() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestErrorAs(t *testing.T) {
	// Helper function for creating target pointer
	createCustomErrorTarget := func() any {
		var target *customError // Note: pointer to concrete type
		return &target
	}

	createErrorTarget := func() any {
		var target error
		return &target
	}

	tests := []struct {
		name          string
		err           error
		targetFactory func() any
		wantPass      bool
		wantTargetVal *customError
	}{
		{"nil error", nil, createCustomErrorTarget, false, nil},
		{"non-nil error, nil target ptr", errCustom, func() any { return nil }, false, nil},
		{"direct match", errCustom, createCustomErrorTarget, true, errCustom},
		{"wrapped match", errCustomWrapped, createCustomErrorTarget, true, errCustom},
		{"no match type", errBase, createCustomErrorTarget, false, nil},
		{"wrapped no match type", errWrapped, createCustomErrorTarget, false, nil},
		// This case expects true because any error satisfies error interface
		{"interface target match", errCustom, createErrorTarget, true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			targetPtr := tt.targetFactory()

			pass := assert.ErrorAs(mockT, tt.err, targetPtr)

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.ErrorAs() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.ErrorAs() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else {
					expectedMsgs := []string{
						"Values are not equa",
						"Expected error of type",
						"Target pointer for ErrorAs cannot be nil",
						"Error cannot be assigned to target type",
					}
					found := false
					for _, msg := range expectedMsgs {
						if strings.Contains(mockT.errorMessages[0], msg) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Errorf message = %q, want one of %v", mockT.errorMessages[0], expectedMsgs)
					}
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}

			// Check target value if successful match expected
			if tt.wantPass && tt.wantTargetVal != nil {
				if targetPtr, ok := targetPtr.(**customError); ok {
					if *targetPtr == nil {
						t.Errorf("Target pointer is nil after successful match")
					} else if (*targetPtr).msg != tt.wantTargetVal.msg {
						t.Errorf("Target value mismatch: got %v, want %v", (*targetPtr).msg, tt.wantTargetVal.msg)
					}
				}
			}
		})
	}
}
