package assert_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/neticdk/go-stdlib/assert"
)

func TestGreater(t *testing.T) {
	tests := []struct {
		name      string
		got       any
		threshold any
		wantPass  bool
	}{
		// Integers
		{"int: 5 > 3", 5, 3, true},
		{"int: 3 > 5", 3, 5, false},
		{"int: 5 > 5", 5, 5, false},
		// Floats
		{"float: 5.1 > 5.0", 5.1, 5.0, true},
		{"float: 5.0 > 5.1", 5.0, 5.1, false},
		{"float: 5.0 > 5.0", 5.0, 5.0, false},
		// Strings
		{"string: b > a", "b", "a", true},
		{"string: a > b", "a", "b", false},
		{"string: a > a", "a", "a", false},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.Greater(mockT, g, tt.threshold.(int))
			case float64:
				pass = assert.Greater(mockT, g, tt.threshold.(float64))
			case string:
				pass = assert.Greater(mockT, g, tt.threshold.(string))
			default:
				t.Fatalf("Unsupported type for Greater test: %T", tt.got)
			}

			if pass != tt.wantPass {
				t.Errorf("Greater() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestGreaterOrEqual(t *testing.T) {
	tests := []struct {
		name      string
		got       any
		threshold any
		wantPass  bool
	}{
		// Integers
		{"int: 5 >= 3", 5, 3, true},
		{"int: 3 >= 5", 3, 5, false},
		{"int: 5 >= 5", 5, 5, true},
		// Floats
		{"float: 5.1 >= 5.0", 5.1, 5.0, true},
		{"float: 5.0 >= 5.1", 5.0, 5.1, false},
		{"float: 5.0 >= 5.0", 5.0, 5.0, true},
		// Strings
		{"string: b >= a", "b", "a", true},
		{"string: a >= b", "a", "b", false},
		{"string: a >= a", "a", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.GreaterOrEqual(mockT, g, tt.threshold.(int))
			case float64:
				pass = assert.GreaterOrEqual(mockT, g, tt.threshold.(float64))
			case string:
				pass = assert.GreaterOrEqual(mockT, g, tt.threshold.(string))
			default:
				t.Fatalf("Unsupported type for Greater test: %T", tt.got)
			}

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.Greater() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.GreaterOrEqual() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], "Value is not greater than or equal to threshold") {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], "Value is not greater than or equal to threshold")
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		name      string
		got       any
		threshold any
		wantPass  bool
	}{
		// Integers
		{"int: 3 < 5", 3, 5, true},
		{"int: 5 < 3", 5, 3, false},
		{"int: 5 < 5", 5, 5, false},
		// Floats
		{"float: 3.0 < 5.0", 3.0, 5.0, true},
		{"float: 5.1 < 5.0", 5.1, 5.0, false},
		{"float: 5.0 < 5.0", 5.0, 5.0, false},
		// Strings
		{"string: a < b", "a", "b", true},
		{"string: b < a", "b", "a", false},
		{"string: a < a", "a", "a", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.Less(mockT, g, tt.threshold.(int))
			case float64:
				pass = assert.Less(mockT, g, tt.threshold.(float64))
			case string:
				pass = assert.Less(mockT, g, tt.threshold.(string))
			default:
				t.Fatalf("Unsupported type for Less test: %T", tt.got)
			}

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.Less() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.Less() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], "Value is not less than threshold") {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], "Value is not less than threshold")
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestLessOrEqual(t *testing.T) {
	tests := []struct {
		name      string
		got       any
		threshold any
		wantPass  bool
	}{
		// Integers
		{"int: 3 <= 5", 3, 5, true},
		{"int: 5 <= 3", 5, 3, false},
		{"int: 5 <= 5", 5, 5, true},
		// Floats
		{"float: 5.0 <= 5.1", 5.0, 5.1, true},
		{"float: 5.1 <= 5.0", 5.1, 5.0, false},
		{"float: 5.0 <= 5.0", 5.0, 5.0, true},
		// Strings
		{"string: a <= b", "a", "b", true},
		{"string: b <= a", "b", "a", false},
		{"string: a <= a", "a", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &mockTestingT{}

			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.LessOrEqual(mockT, g, tt.threshold.(int))
			case float64:
				pass = assert.LessOrEqual(mockT, g, tt.threshold.(float64))
			case string:
				pass = assert.LessOrEqual(mockT, g, tt.threshold.(string))
			default:
				t.Fatalf("Unsupported type for LessOrEqual test: %T", tt.got)
			}

			// Check the boolean return value
			if pass != tt.wantPass {
				t.Errorf("assert.LessOrEqual() returned = %v, wantPass %v", pass, tt.wantPass)
			}

			// Check if `Errorf` was called (or not) as expected
			reportedFailure := mockT.Failed()
			expectedFailure := !tt.wantPass

			if reportedFailure != expectedFailure {
				t.Errorf("assert.LessOrEqual() called Errorf? = %v, expected? %v", reportedFailure, expectedFailure)
			}

			// If failure was expected, check error messages
			if expectedFailure {
				if len(mockT.errorMessages) == 0 {
					t.Errorf("Expected Errorf call, but no error messages recorded")
				} else if !strings.Contains(mockT.errorMessages[0], "Value is not less than or equal to threshold") {
					t.Errorf("Errorf message = %q, want substring %q", mockT.errorMessages[0], "Value is not less than or equal to threshold")
				}
			} else if len(mockT.errorMessages) > 0 {
				t.Errorf("Expected no Errorf call, but got messages: %v", mockT.errorMessages)
			}
		})
	}
}

func TestPositive(t *testing.T) {
	tests := []struct {
		name     string
		got      any
		wantPass bool
	}{
		// Integers
		{"int: 5", 5, true},
		{"int: 0", 0, false},
		{"int: -5", -5, false},
		// Floats
		{"float: 5.1", 5.1, true},
		{"float: 0.0", 0.0, false},
		{"float: -5.1", -5.1, false},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.Positive(mockT, g)
			case float64:
				pass = assert.Positive(mockT, g)
			default:
				t.Fatalf("Unsupported type for Positive test: %T", tt.got)
			}

			if pass != tt.wantPass {
				t.Errorf("Positive() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestNegative(t *testing.T) {
	tests := []struct {
		name     string
		got      any
		wantPass bool
	}{
		// Integers
		{"int: -5", -5, true},
		{"int: 0", 0, false},
		{"int: 5", 5, false},
		// Floats
		{"float: -5.1", -5.1, true},
		{"float: 0.0", 0.0, false},
		{"float: 5.1", 5.1, false},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pass bool
			switch g := tt.got.(type) {
			case int:
				pass = assert.Negative(mockT, g)
			case float64:
				pass = assert.Negative(mockT, g)
			default:
				t.Fatalf("Unsupported type for Negative test: %T", tt.got)
			}

			if pass != tt.wantPass {
				t.Errorf("Negative() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}

func TestZero(t *testing.T) {
	mockT := new(testing.T)

	// Use a helper that explicitly passes *testing.T
	runZeroTest := func(t *testing.T, name string, data any, wantPass bool) {
		t.Run(name, func(t *testing.T) {
			pass := assert.Zero(mockT, data)
			if pass != wantPass {
				t.Errorf("Zero(%#v) assertion result = %v, wantPass %v", data, pass, wantPass)
			}
		})
	}

	// Test cases
	var nilIntPtr *int
	var zeroInt int = 0
	var nonZeroInt int = 1
	var nilMap map[string]int
	emptyMap := map[string]int{}
	var nilSlice []int
	emptySlice := []int{}
	var zeroStruct struct {
		i int
		s string
	}
	nonZeroStruct := struct{ i int }{i: 1}
	var nilChan chan int
	emptyChan := make(chan int) // Non-nil but zero struct
	var nilFunc func()
	var zeroString string = ""
	var nonZeroString string = "a"
	var zeroBool bool = false
	var nonZeroBool bool = true

	runZeroTest(t, "untyped nil", nil, true)
	runZeroTest(t, "nil int pointer", nilIntPtr, true)
	runZeroTest(t, "zero int", zeroInt, true)
	runZeroTest(t, "non-zero int", nonZeroInt, false)
	runZeroTest(t, "nil map", nilMap, true)
	runZeroTest(t, "empty map", emptyMap, false) // empty map is not the zero value
	runZeroTest(t, "nil slice", nilSlice, true)
	runZeroTest(t, "empty slice", emptySlice, false) // empty slice is not the zero value
	runZeroTest(t, "zero struct", zeroStruct, true)
	runZeroTest(t, "non-zero struct", nonZeroStruct, false)
	runZeroTest(t, "nil channel", nilChan, true)
	runZeroTest(t, "made channel", emptyChan, false) // made channel is not the zero value
	runZeroTest(t, "nil func", nilFunc, true)
	runZeroTest(t, "zero string", zeroString, true)
	runZeroTest(t, "non-zero string", nonZeroString, false)
	runZeroTest(t, "zero bool", zeroBool, true)
	runZeroTest(t, "non-zero bool", nonZeroBool, false)

	// Test reflection specific edge cases
	runZeroTest(t, "reflect.Value of zero", reflect.ValueOf(0), false) // A reflect.Value itself is not zero unless uninitialized
	var zeroVal reflect.Value
	runZeroTest(t, "uninitialized reflect.Value", zeroVal, true) // The zero reflect.Value IS zero
}

func TestNotZero(t *testing.T) {
	mockT := new(testing.T)

	runNotZeroTest := func(t *testing.T, name string, data any, wantPass bool) {
		t.Run(name, func(t *testing.T) {
			pass := assert.NotZero(mockT, data)
			if pass != wantPass {
				t.Errorf("NotZero(%#v) assertion result = %v, wantPass %v", data, pass, wantPass)
			}
		})
	}

	var nilIntPtr *int
	var zeroInt int = 0
	var nonZeroInt int = 1
	var nilMap map[string]int
	emptyMap := map[string]int{}
	var nilSlice []int
	emptySlice := []int{}
	var zeroStruct struct {
		i int
		s string
	}
	nonZeroStruct := struct{ i int }{i: 1}
	var nilChan chan int
	emptyChan := make(chan int)
	var nilFunc func()
	var zeroString string = ""
	var nonZeroString string = "a"
	var zeroBool bool = false
	var nonZeroBool bool = true

	runNotZeroTest(t, "untyped nil", nil, false)
	runNotZeroTest(t, "nil int pointer", nilIntPtr, false)
	runNotZeroTest(t, "zero int", zeroInt, false)
	runNotZeroTest(t, "non-zero int", nonZeroInt, true)
	runNotZeroTest(t, "nil map", nilMap, false)
	runNotZeroTest(t, "empty map", emptyMap, true)
	runNotZeroTest(t, "nil slice", nilSlice, false)
	runNotZeroTest(t, "empty slice", emptySlice, true)
	runNotZeroTest(t, "zero struct", zeroStruct, false)
	runNotZeroTest(t, "non-zero struct", nonZeroStruct, true)
	runNotZeroTest(t, "nil channel", nilChan, false)
	runNotZeroTest(t, "made channel", emptyChan, true)
	runNotZeroTest(t, "nil func", nilFunc, false)
	runNotZeroTest(t, "zero string", zeroString, false)
	runNotZeroTest(t, "non-zero string", nonZeroString, true)
	runNotZeroTest(t, "zero bool", zeroBool, false)
	runNotZeroTest(t, "non-zero bool", nonZeroBool, true)

	var zeroVal reflect.Value
	runNotZeroTest(t, "uninitialized reflect.Value", zeroVal, false)
	runNotZeroTest(t, "reflect.Value of zero", reflect.ValueOf(0), true)
}

func TestElementsMatch(t *testing.T) {
	tests := []struct {
		name     string
		listA    any
		listB    any
		wantPass bool
	}{
		{
			name:     "int: matching slices different order",
			listA:    []int{1, 2, 3},
			listB:    []int{3, 1, 2},
			wantPass: true,
		},
		{
			name:     "int: matching slices with duplicates",
			listA:    []int{1, 2, 2, 3},
			listB:    []int{3, 1, 2, 2},
			wantPass: true,
		},
		{
			name:     "int: different lengths",
			listA:    []int{1, 2},
			listB:    []int{1, 2, 3},
			wantPass: false,
		},
		{
			name:     "int: different element counts",
			listA:    []int{1, 2, 2, 3},
			listB:    []int{1, 2, 3, 3},
			wantPass: false,
		},
		{
			name:     "int: empty slices",
			listA:    []int{},
			listB:    []int{},
			wantPass: true,
		},
		{
			name:     "string: matching slices different order",
			listA:    []string{"a", "b", "c"},
			listB:    []string{"c", "a", "b"},
			wantPass: true,
		},
		{
			name:     "string: matching slices with duplicates",
			listA:    []string{"a", "b", "b", "c"},
			listB:    []string{"c", "a", "b", "b"},
			wantPass: true,
		},
		{
			name:     "string: different lengths",
			listA:    []string{"a", "b"},
			listB:    []string{"a", "b", "c"},
			wantPass: false,
		},
		{
			name:     "string: different element counts",
			listA:    []string{"a", "b", "b", "c"},
			listB:    []string{"a", "b", "c", "c"},
			wantPass: false,
		},
		{
			name:     "string: empty slices",
			listA:    []string{},
			listB:    []string{},
			wantPass: true,
		},
	}

	mockT := new(testing.T)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pass bool
			switch a := tt.listA.(type) {
			case []int:
				b, ok := tt.listB.([]int)
				if !ok {
					t.Fatalf("Type mismatch for listB: expected []int, got %T", tt.listB)
				}
				pass = assert.ElementsMatch(mockT, a, b)
			case []string:
				b, ok := tt.listB.([]string)
				if !ok {
					t.Fatalf("Type mismatch for listB: expected []string, got %T", tt.listB)
				}
				pass = assert.ElementsMatch(mockT, a, b)
			default:
				t.Fatalf("Unsupported type for ElementsMatch test: %T", tt.listA)
			}

			if pass != tt.wantPass {
				t.Errorf("ElementsMatch() assertion result = %v, wantPass %v", pass, tt.wantPass)
			}
		})
	}
}
