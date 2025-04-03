package assert_test

import (
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestTrue(t *testing.T) {
	tests := []struct {
		name            string
		value           bool
		wantPass        bool   // Expected return value from `assert.True`
		wantErrorMsg    string // Substring expected in `Errorf` message if `wantPass` is false
		optionalMsg     []any  // Optional message to pass to `assert.True`
		wantOptionalLog string // Expected substring in `Logf`/`Log` if `optionalMsg` provided on failure
	}{
		{
			name:         "true passes",
			value:        true,
			wantPass:     true,
			wantErrorMsg: "",
		},
		{
			name:         "false fails",
			value:        false,
			wantPass:     false,
			wantErrorMsg: "Values are not equal",
		},
		{
			name:            "false fails with custom message",
			value:           false,
			wantPass:        false,
			wantErrorMsg:    "Values are not equal",
			optionalMsg:     []any{"custom reason %d", 123},
			wantOptionalLog: "custom reason 123",
		},
		{
			name:            "false fails with unformatted custom message",
			value:           false,
			wantPass:        false,
			wantErrorMsg:    "Values are not equal",
			optionalMsg:     []any{"another", "reason"},
			wantOptionalLog: "another reason",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			pass := assert.True(mockT, tt.value, tt.optionalMsg...)

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.True() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass // Failure expected if `wantPass` is false

			if reportedFailure != expectedFailure {
				t.Errorf("assert.True() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check the error message content
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], tt.wantErrorMsg) {
					// Check only the first error message for the primary failure reason
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], tt.wantErrorMsg)
				}

				// If failure and custom message expected, check logs
				if tt.wantOptionalLog != "" {
					foundLog := false
					for _, logMsg := range mockT.logMessages { // Check all log messages
						if strings.Contains(logMsg, tt.wantOptionalLog) {
							foundLog = true
							break
						}
					}
					if !foundLog {
						t.Errorf("Expected log message containing %q, got logs: %v", tt.wantOptionalLog, mockT.logMessages)
					}
				}
			} else {
				// If success was expected, ensure no error messages were logged
				if len(mockT.errorMessages) > 0 {
					t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
				}
			}
		})
	}
}

func TestFalse(t *testing.T) {
	tests := []struct {
		name            string
		value           bool
		wantPass        bool   // Expected return value from `assert.False`
		wantErrorMsg    string // Substring expected in `Errorf` message if `wantPass` is false
		optionalMsg     []any  // Optional message to pass to `assert.False`
		wantOptionalLog string // Expected substring in Logf/Log if `optionalMsg` provided on failure
	}{
		{
			name:         "false passes",
			value:        false,
			wantPass:     true,
			wantErrorMsg: "",
		},
		{
			name:         "true fails",
			value:        true,
			wantPass:     false,
			wantErrorMsg: "Values are not equal",
		},
		{
			name:            "true fails with custom message",
			value:           true,
			wantPass:        false,
			wantErrorMsg:    "Values are not equal",
			optionalMsg:     []any{"custom reason for true %s", "value"},
			wantOptionalLog: "custom reason for true value",
		},
		{
			name:            "true fails with unformatted custom message",
			value:           true,
			wantPass:        false,
			wantErrorMsg:    "Values are not equal",
			optionalMsg:     []any{"another reason"},
			wantOptionalLog: "another reason",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			// Call the function under test
			pass := assert.False(mockT, tt.value, tt.optionalMsg...)

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.False() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass // Failure expected if `wantPass` is false

			if reportedFailure != expectedFailure {
				t.Errorf("assert.False() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check the error message content
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], tt.wantErrorMsg) {
					// Check only the first error message for the primary failure reason
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], tt.wantErrorMsg)
				}

				// If failure and custom message expected, check logs
				if tt.wantOptionalLog != "" {
					foundLog := false
					for _, logMsg := range mockT.logMessages { // Check all log messages
						if strings.Contains(logMsg, tt.wantOptionalLog) {
							foundLog = true
							break
						}
					}
					if !foundLog {
						t.Errorf("Expected log message containing %q, got logs: %v", tt.wantOptionalLog, mockT.logMessages)
					}
				}
			} else {
				// If success was expected, ensure no error messages were logged
				if len(mockT.errorMessages) > 0 {
					t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
				}
			}
		})
	}
}
