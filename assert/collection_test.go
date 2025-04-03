package assert_test

import (
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name       string
		collection any
		element    any
		wantPass   bool
	}{
		{"string contains substring", "hello world", "world", true},
		{"string does not contain substring", "hello world", "golang", false},
		{"slice contains element", []int{1, 2, 3}, 2, true},
		{"slice does not contain element", []int{1, 2, 3}, 4, false},
		{"map contains value", map[string]int{"a": 1, "b": 2}, 2, true},
		{"map does not contain value", map[string]int{"a": 1, "b": 2}, 3, false},
		{"invalid collection type", 123, 3, false}, // triggers the err != nil in Contains
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			pass := assert.Contains(mockT, tt.collection, tt.element)

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.Contains() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.Contains() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else {
					// Check for any of the possible error messages
					expectedSubstrings := []string{
						"Collection does not contain element",
						"String does not contain substring",
						"Map does not contain value",
						"Error checking contains",
					}
					found := false
					for _, substr := range expectedSubstrings {
						if strings.Contains(mockT.errorMessages[0], substr) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Errorf message = %q, want one of %v", mockT.errorMessages[0], expectedSubstrings)
					}
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestNotContains(t *testing.T) {
	tests := []struct {
		name       string
		collection any
		element    any
		wantPass   bool
	}{
		{"string does not contain substring", "hello world", "golang", true},
		{"string contains substring", "hello world", "world", false},
		{"slice does not contain element", []int{1, 2, 3}, 4, true},
		{"slice contains element", []int{1, 2, 3}, 2, false},
		{"map does not contain value", map[string]int{"a": 1, "b": 2}, 3, true},
		{"map contains value", map[string]int{"a": 1, "b": 2}, 2, false},
		{"invalid collection type", 123, 3, false}, // triggers the err != nil in NotContains
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			pass := assert.NotContains(mockT, tt.collection, tt.element)

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.NotContains() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.NotContains() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else {
					// Check for any of the possible error messages
					expectedSubstrings := []string{
						"Collection should not contain element",
						"String should not contain substring",
						"Map should not contain value",
						"Error checking not-contains",
					}
					found := false
					for _, substr := range expectedSubstrings {
						if strings.Contains(mockT.errorMessages[0], substr) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Errorf message = %q, want one of %v", mockT.errorMessages[0], expectedSubstrings)
					}
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestContainsKey(t *testing.T) {
	tests := []struct {
		name     string
		m        any
		key      any
		wantPass bool
	}{
		{"string key present", map[string]int{"a": 1, "b": 2}, "a", true},
		{"string key absent", map[string]int{"a": 1, "b": 2}, "c", false},
		{"int key present", map[int]bool{1: true, 2: false}, 2, true},
		{"int key absent", map[int]bool{1: true, 2: false}, 3, false},
		{"nil map", (map[string]int)(nil), "a", false},
		{"empty map", map[string]int{}, "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			var pass bool
			switch typedMap := tt.m.(type) {
			case map[string]int:
				typedKey, ok := tt.key.(string)
				if !ok {
					t.Fatalf("Key type mismatch: expected string, got %T", tt.key)
				}
				pass = assert.ContainsKey(mockT, typedMap, typedKey)
			case map[int]bool:
				typedKey, ok := tt.key.(int)
				if !ok {
					t.Fatalf("Key type mismatch: expected int, got %T", tt.key)
				}
				pass = assert.ContainsKey(mockT, typedMap, typedKey)
			default:
				if tt.m == nil {
					if _, ok := tt.key.(string); ok {
						pass = assert.ContainsKey(mockT, (map[string]int)(nil), tt.key.(string))
					} else {
						t.Logf("Skipping nil map test for key type %T", tt.key)
						pass = false
					}
				} else {
					t.Fatalf("Unsupported map type for ContainsKey test: %T", tt.m)
				}
			}

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.ContainsKey() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.ContainsKey() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], "Map does not contain key") {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], "Map does not contain key")
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestNotContainsKey(t *testing.T) {
	tests := []struct {
		name     string
		m        any
		key      any
		wantPass bool
	}{
		{"string key absent", map[string]int{"a": 1, "b": 2}, "c", true},
		{"string key present", map[string]int{"a": 1, "b": 2}, "a", false},
		{"int key absent", map[int]bool{1: true, 2: false}, 3, true},
		{"int key present", map[int]bool{1: true, 2: false}, 2, false},
		{"nil map", (map[string]int)(nil), "a", true}, // nil map doesn't contain any key
		{"empty map", map[string]int{}, "a", true},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pass bool
			switch typedMap := tt.m.(type) {
			case map[string]int:
				typedKey, ok := tt.key.(string)
				if !ok {
					t.Fatalf("Key type mismatch: expected string, got %T", tt.key)
				}
				pass = assert.NotContainsKey(mockT, typedMap, typedKey)
			case map[int]bool:
				typedKey, ok := tt.key.(int)
				if !ok {
					t.Fatalf("Key type mismatch: expected int, got %T", tt.key)
				}
				pass = assert.NotContainsKey(mockT, typedMap, typedKey)
			default:
				if tt.m == nil {
					if _, ok := tt.key.(string); ok {
						pass = assert.NotContainsKey(mockT, (map[string]int)(nil), tt.key.(string))
					} else {
						t.Logf("Skipping nil map test for key type %T", tt.key)
						pass = true
					}
				} else {
					t.Fatalf("Unsupported map type for NotContainsKey test: %T", tt.m)
				}
			}

			if pass != tt.wantPass {
				t.Errorf("NotContainsKey() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		wantPass bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"zero int", 0, true},
		{"empty slice", []int{}, true},
		{"empty map", map[string]int{}, true},
		{"non-empty string", "hello", false},
		{"non-zero int", 1, false},
		{"non-empty slice", []int{1, 2, 3}, false},
		{"non-empty map", map[string]int{"a": 1}, false},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := assert.Empty(mockT, tt.data)
			if pass != tt.wantPass {
				t.Errorf("Empty() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestNotEmpty(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		wantPass bool
	}{
		{"nil", nil, false},
		{"empty string", "", false},
		{"zero int", 0, false},
		{"empty slice", []int{}, false},
		{"empty map", map[string]int{}, false},
		{"non-empty string", "hello", true},
		{"non-zero int", 1, true},
		{"non-empty slice", []int{1, 2, 3}, true},
		{"non-empty map", map[string]int{"a": 1}, true},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := assert.NotEmpty(mockT, tt.data)
			if pass != tt.wantPass {
				t.Errorf("NotEmpty() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name        string
		data        any
		expectedLen int
		wantPass    bool
	}{
		{"slice with correct length", []int{1, 2, 3}, 3, true},
		{"slice with incorrect length", []int{1, 2, 3}, 2, false},
		{"map with correct length", map[string]int{"a": 1, "b": 2}, 2, true},
		{"map with incorrect length", map[string]int{"a": 1, "b": 2}, 3, false},
		{"string with correct length", "hello", 5, true},
		{"string with incorrect length", "hello", 4, false},
		{"unsupported type", 123, 0, false}, // triggers the default case
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := assert.Len(mockT, tt.data, tt.expectedLen)
			if pass != tt.wantPass {
				t.Errorf("Len() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}
