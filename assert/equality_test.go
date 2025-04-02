package assert_test

import (
	"strings"
	"testing"
	"time"

	"github.com/neticdk/go-stdlib/assert"
)

func TestEqual(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		got      any
		want     any
		wantPass bool
	}{
		{"int: 5 == 5", 5, 5, true},
		{"int: 5 == 3", 5, 3, false},
		{"string: a == a", "a", "a", true},
		{"string: a == b", "a", "b", false},
		{"bool: true == true", true, true, true},
		{"bool: true == false", true, false, false},
		{"time: equal times", now, now, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.Equal(mockT, g, tt.want.(int))
			case string:
				pass = assert.Equal(mockT, g, tt.want.(string))
			case bool:
				pass = assert.Equal(mockT, g, tt.want.(bool))
			case time.Time:
				pass = assert.Equal(mockT, g, tt.want.(time.Time))
			default:
				t.Fatalf("Unsupported type for Equal test: %T", tt.got)
			}

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.Equal() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.Equal() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], "Values are not equal") {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], "Values are not equal")
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name     string
		got      any
		want     any
		wantPass bool
	}{
		{"int: 5 != 5", 5, 5, false},
		{"int: 5 != 3", 5, 3, true},
		{"string: a != a", "a", "a", false},
		{"string: a != b", "a", "b", true},
		{"bool: true != true", true, true, false},
		{"bool: true != false", true, false, true},
		{"time: not equal times", time.Now(), time.Now().Add(time.Minute), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.NotEqual(mockT, g, tt.want.(int))
			case string:
				pass = assert.NotEqual(mockT, g, tt.want.(string))
			case bool:
				pass = assert.NotEqual(mockT, g, tt.want.(bool))
			case time.Time:
				pass = assert.NotEqual(mockT, g, tt.want.(time.Time))
			default:
				t.Fatalf("Unsupported type for NotEqual test: %T", tt.got)
			}

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.NotEqual() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.NotEqual() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], "Values should not be equal") {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], "Values should not be equal")
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestInDelta(t *testing.T) {
	tests := []struct {
		name     string
		got      float64
		want     float64
		delta    float64
		wantPass bool
	}{
		{"within delta", 5.0, 5.1, 0.2, true},
		{"outside delta", 5.0, 5.3, 0.2, false},
		{"equal", 5.0, 5.0, 0.2, true},
		{"negative within delta", -5.0, -5.1, 0.2, true},
		{"negative outside delta", -5.0, -5.3, 0.2, false},
		{"negative equal", -5.0, -5.0, 0.2, true},
		{"zero delta, equal values", 5.0, 5.0, 0.0, true},
		{"zero delta, different values", 5.0, 5.1, 0.0, false},
		{"negative delta within delta", 5.0, 5.1, -0.2, true},
		{"negative delta outside delta", 5.0, 5.3, -0.2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			pass := assert.InDelta(mockT, tt.got, tt.want, tt.delta)

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.InDelta() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.InDelta() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], "Values are not within delta") {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], "Values are not within delta")
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestNotInDelta(t *testing.T) {
	tests := []struct {
		name         string
		got          float64
		want         float64
		delta        float64
		wantPass     bool
		wantErrorMsg string
	}{
		{"outside delta", 5.3, 5.0, 0.2, true, ""},
		{"barely outside delta", 5.21, 5.0, 0.2, true, ""},
		{"within delta", 5.1, 5.0, 0.2, false, "Values are within delta"},
		{"equal to delta boundary", 5.2, 5.0, 0.2, false, ""},
		{"equal values", 5.0, 5.0, 0.2, false, "Values are within delta"},
		{"negative outside delta", -5.3, -5.0, 0.2, true, ""},
		{"negative within delta", -5.1, -5.0, 0.2, false, "Values are within delta"},
		{"negative equal", -5.0, -5.0, 0.2, false, "Values are within delta"},
		{"zero delta, different values", 5.1, 5.0, 0.0, true, ""},
		{"zero delta, equal values", 5.0, 5.0, 0.0, false, "Values are within delta"},
		{"negative delta outside", 5.3, 5.0, -0.2, true, ""},
		{"negative delta within", 5.1, 5.0, -0.2, false, "Values are within delta"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			pass := assert.NotInDelta(mockT, tt.got, tt.want, tt.delta)

			// Check return value
			if pass != tt.wantPass {
				t.Errorf("assert.NotInDelta() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check failure reporting
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.NotInDelta() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// Check error message content on failure
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], tt.wantErrorMsg) {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], tt.wantErrorMsg)
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}
