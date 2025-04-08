package require_test

import (
	"errors"
	"fmt" // Import fmt for Errorf implementation
	"testing"
	"time"

	"github.com/neticdk/go-stdlib/require"
)

// --- bool ---

func TestTrue(t *testing.T) {
	require.True(t, true, "True should pass for true")
}

func TestFalse(t *testing.T) {
	require.False(t, false, "False should pass for false")
}

// --- collection ---

func TestContains(t *testing.T) {
	require.Contains(t, []int{1, 2, 3}, 2, "Contains should find existing element")
	require.Contains(t, map[string]string{"a": "a"}, "a", "Contains should find existing key in map")
}

func TestNotContains(t *testing.T) {
	require.NotContains(t, []int{1, 2, 3}, 4, "NotContains should not find non-existing element")
	require.NotContains(t, map[string]int{"a": 1}, "b", "NotContains should not find non-existing key in map")
}

func TestContainsKey(t *testing.T) {
	require.ContainsKey(t, map[string]int{"a": 1}, "a", "ContainsKey should find existing key in map")
}

func TestNotContainsKey(t *testing.T) {
	require.NotContainsKey(t, map[string]int{"a": 1}, "b", "NotContainsKey should not find non-existing key in map")
}

func TestEmpty(t *testing.T) {
	require.Empty(t, nil, "Empty should pass for nil")
	require.Empty(t, "", "Empty should pass for empty string")
	require.Empty(t, 0, "Empty should pass for 0")
	require.Empty(t, false, "Empty should pass for false")
	require.Empty(t, []int{}, "Empty should pass for empty slice")
	require.Empty(t, map[string]string{}, "Empty should pass for empty map")
	var ch chan int
	require.Empty(t, ch, "Empty should pass for nil channel")
}

func TestNotEmpty(t *testing.T) {
	require.NotEmpty(t, "a", "NotEmpty should pass for non-empty string")
	require.NotEmpty(t, 1, "NotEmpty should pass for non-zero int")
	require.NotEmpty(t, true, "NotEmpty should pass for true")
	require.NotEmpty(t, []int{1}, "NotEmpty should pass for non-empty slice")
	require.NotEmpty(t, map[string]string{"a": "b"}, "NotEmpty should pass for non-empty map")
	ch := make(chan int, 1)
	ch <- 1 // Send a value
	require.NotEmpty(t, ch, "NotEmpty should pass for non-nil channel")
}

func TestLen(t *testing.T) {
	require.Len(t, []int{1, 2, 3}, 3, "Len should pass for correct slice length")
	require.Len(t, "abc", 3, "Len should pass for correct string length")
	require.Len(t, map[string]int{"a": 1, "b": 2}, 2, "Len should pass for correct map length")
}

// --- comparison ---

func TestZero(t *testing.T) {
	require.Zero(t, 0, "Zero should pass for 0 int")
	require.Zero(t, 0.0, "Zero should pass for 0.0 float")
	require.Zero(t, "", "Zero should pass for empty string")
	require.Zero(t, false, "Zero should pass for false")
	var p *int
	require.Zero(t, p, "Zero should pass for nil pointer")
}

func TestNotZero(t *testing.T) {
	require.NotZero(t, 1, "NotZero should pass for 1")
	require.NotZero(t, "a", "NotZero should pass for non-empty string")
	require.NotZero(t, true, "NotZero should pass for true")
	p := 5
	require.NotZero(t, &p, "NotZero should pass for non-nil pointer")
}

func TestGreater(t *testing.T) {
	require.Greater(t, 5, 4, "Greater should pass when first is greater")
	require.Greater(t, 5.1, 5.0, "Greater should pass for floats")
}

func TestGreaterOrEqual(t *testing.T) {
	require.GreaterOrEqual(t, 5, 4, "GreaterOrEqual should pass when first is greater")
	require.GreaterOrEqual(t, 5, 5, "GreaterOrEqual should pass when equal")
	require.GreaterOrEqual(t, 5.0, 5.0, "GreaterOrEqual should pass for floats")
}

func TestLess(t *testing.T) {
	require.Less(t, 4, 5, "Less should pass when first is less")
	require.Less(t, 5.0, 5.1, "Less should pass for floats")
}

func TestLessOrEqual(t *testing.T) {
	require.LessOrEqual(t, 4, 5, "LessOrEqual should pass when first is less")
	require.LessOrEqual(t, 5, 5, "LessOrEqual should pass when equal")
	require.LessOrEqual(t, 5.0, 5.0, "LessOrEqual should pass for floats")
}

func TestPositive(t *testing.T) {
	require.Positive(t, 1, "Positive should pass for positive int")
	require.Positive(t, 0.1, "Positive should pass for positive float")
}

func TestNegative(t *testing.T) {
	require.Negative(t, -1, "Negative should pass for negative int")
	require.Negative(t, -0.1, "Negative should pass for negative float")
}

func TestElementsMatch(t *testing.T) {
	require.ElementsMatch(t, []int{1, 2, 3}, []int{3, 1, 2}, "ElementsMatch should pass for same elements in different order")
	require.ElementsMatch(t, []int{1, 1, 2}, []int{1, 2, 1}, "ElementsMatch should pass for duplicates")
	require.ElementsMatch(t, []string{"a", "b"}, []string{"b", "a"}, "ElementsMatch should pass for strings")
}

// --- equality ---

func TestEqual(t *testing.T) {
	require.Equal(t, 5, 5, "Equal should pass for equal ints")
	require.Equal(t, "hello", "hello", "Equal should pass for equal strings")
	require.Equal(t, []int{1, 2}, []int{1, 2}, "Equal should pass for equal slices") // Note: deep equality handled by assert.Equal
}

func TestNotEqual(t *testing.T) {
	require.NotEqual(t, 5, 6, "NotEqual should pass for different ints")
	require.NotEqual(t, "hello", "world", "NotEqual should pass for different strings")
	require.NotEqual(t, []int{1, 2}, []int{1, 3}, "NotEqual should pass for different slices")
}

func TestInDelta(t *testing.T) {
	require.InDelta(t, 10.1, 10.0, 0.2, "InDelta should pass when within delta")
	require.InDelta(t, 10.0, 10.0, 0.1, "InDelta should pass when equal")
}

func TestNotInDelta(t *testing.T) {
	// This should pass and not call FailNow
	require.NotInDelta(t, 10.3, 10.0, 0.2, "NotInDelta should pass when outside delta")
	require.NotInDelta(t, 10.0, 10.3, 0.2, "NotInDelta should pass when outside delta (reversed)")
	require.NotInDelta(t, 10.1, 10.0, 0.0, "NotInDelta should pass with zero delta and different values")
}

// --- error ---

func TestError(t *testing.T) {
	err := errors.New("test error")
	require.Error(t, err, "Error should pass for non-nil error")
}

func TestNoError(t *testing.T) {
	var err error
	require.NoError(t, err, "NoError should pass for nil error")
}

func TestErrorIs(t *testing.T) {
	err := errors.New("test error")
	require.ErrorIs(t, err, err, "ErrorIs should pass when error matches target value")

	baseErr := errors.New("base error for wrapping")
	wrappedErr := fmt.Errorf("some context: %w", baseErr)
	require.ErrorIs(t, wrappedErr, baseErr, "ErrorIs should pass when error wraps target")
}

func TestErrorAs(t *testing.T) {
	err := errors.New("test error")
	var target error
	require.ErrorAs(t, err, &target, "ErrorAs should pass when error is assignable to target")
}

// --- nil ---

func TestNil(t *testing.T) {
	var p *int
	require.Nil(t, p, "Nil should pass for nil pointer")
	var i any
	require.Nil(t, i, "Nil should pass for nil interface")
}

func TestNotNil(t *testing.T) {
	p := 5
	require.NotNil(t, &p, "NotNil should pass for non-nil pointer")
	var i any = "hello"
	require.NotNil(t, i, "NotNil should pass for non-nil interface")
}

// --- panic ---

func TestPanics(t *testing.T) {
	require.Panics(t, func() { panic("oh no") }, "Panics should pass when function panics")
}

func TestNotPanics(t *testing.T) {
	require.NotPanics(t, func() {}, "NotPanics should pass when function does not panic")
}

// --- time ---

func TestTimeAfter(t *testing.T) {
	now := time.Now()
	after := now.Add(1 * time.Hour)

	// Test passing case
	require.TimeAfter(t, after, now, "Time should be after now")
}

func TestTimeBefore(t *testing.T) {
	now := time.Now()
	before := now.Add(-1 * time.Hour)

	// Test passing case
	require.TimeBefore(t, before, now, "Time should be before now")
}

func TestTimeEqual(t *testing.T) {
	now := time.Now()
	sameTimeOtherZone := now.In(time.FixedZone("UTC+2", 2*60*60))

	// Test passing case
	require.TimeEqual(t, now, now, "Identical times should be equal")
	require.TimeEqual(t, now, sameTimeOtherZone, "Same time in different zone should be equal")
}

func TestWithinDuration(t *testing.T) {
	now := time.Now()
	nearTime := now.Add(30 * time.Second)

	// Test passing case
	require.WithinDuration(t, nearTime, now, 1*time.Minute, "Times should be within 1 minute")
}

func TestTimeEqualWithPrecision(t *testing.T) {
	baseTime := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	sameMinute := time.Date(2023, 5, 15, 10, 30, 45, 0, time.UTC)

	// Test passing case
	require.TimeEqualWithPrecision(t, baseTime, sameMinute, time.Minute, "Times should be equal when truncated to minutes")
}

func TestWithinTime(t *testing.T) {
	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)

	// Test passing case
	require.WithinTime(t, now, start, end, "Time should be within the window")
}

// --- Failure Path Tests ---

// mockTestingT is a mock implementation of require.TestingT and tHelper
// for testing the require functions themselves.
type mockTestingT struct {
	failedNow bool     // Flag to record if FailNow was called
	errors    []string // Store error messages
}

// Implement require.tHelper
func (m *mockTestingT) Helper() {} // No-op for the mock

// Implement require.TestingT
func (m *mockTestingT) Log(args ...any)                 {} // No-op
func (m *mockTestingT) Logf(format string, args ...any) {} // No-op

func (m *mockTestingT) Errorf(format string, args ...any) {
	// Store the error message to potentially assert its content later
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockTestingT) FailNow() {
	m.failedNow = true
	// IMPORTANT: Do not actually call panic or runtime.Goexit() here.
	// Just record that it *would* have been called.
}

// Helper assertion for the tests below
func assertFailedNow(t *testing.T, mock *mockTestingT, funcName string) {
	t.Helper()
	if !mock.failedNow {
		t.Errorf("%s did not call FailNow() when it should have", funcName)
	}
	// Optional: Check if Errorf was also called (most require funcs should call it before FailNow)
	if len(mock.errors) == 0 && mock.failedNow { // Added check for mock.failedNow to avoid error when FailNow wasn't called
		t.Logf("Note: %s called FailNow() but did not log an error via Errorf()", funcName) // Use t.Logf for informational messages
	}
}

func TestTrue_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.True(mockT, false)
	assertFailedNow(t, mockT, "require.True")
}

func TestFalse_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.False(mockT, true)
	assertFailedNow(t, mockT, "require.False")
}

func TestContains_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Contains(mockT, []int{1, 3}, 2)
	assertFailedNow(t, mockT, "require.Contains")
}

func TestNotContains_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.NotContains(mockT, []int{1, 2, 3}, 2)
	assertFailedNow(t, mockT, "require.NotContains")
}

func TestContainsKey_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.ContainsKey(mockT, map[string]int{"a": 1}, "b")
	assertFailedNow(t, mockT, "require.ContainsKey")
}

func TestNotContainsKey_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.NotContainsKey(mockT, map[string]int{"a": 1}, "a")
	assertFailedNow(t, mockT, "require.NotContainsKey")
}

func TestEqual_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Equal(mockT, 1, 2)
	assertFailedNow(t, mockT, "require.Equal")
}

func TestNotEqual_Failure(t *testing.T) {
	mockT := &mockTestingT{}

	require.NotEqual(mockT, 1, 1)
	assertFailedNow(t, mockT, "require.NotEqual")
}

func TestInDelta_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.InDelta(mockT, 10.3, 10.0, 0.2) // 10.3 is outside the 0.2 delta from 10.0
	assertFailedNow(t, mockT, "require.InDelta")
}

func TestNotInDelta_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	// 10.1 is within 0.2 delta of 10.0, so assert.NotInDelta fails, require.NotInDelta should FailNow
	require.NotInDelta(mockT, 10.1, 10.0, 0.2)
	assertFailedNow(t, mockT, "require.NotInDelta (within delta)")

	mockT = &mockTestingT{}
	// 10.2 is exactly at the delta boundary, so assert.NotInDelta fails (diff is not > delta)
	require.NotInDelta(mockT, 10.2, 10.0, 0.2)
	assertFailedNow(t, mockT, "require.NotInDelta (at delta boundary)")

	mockT = &mockTestingT{}
	// Equal values are within any non-negative delta
	require.NotInDelta(mockT, 10.0, 10.0, 0.1)
	assertFailedNow(t, mockT, "require.NotInDelta (equal values)")

	mockT = &mockTestingT{}
	// Equal values are within zero delta
	require.NotInDelta(mockT, 10.0, 10.0, 0.0)
	assertFailedNow(t, mockT, "require.NotInDelta (equal values, zero delta)")
}

func TestError_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	var err error
	require.Error(mockT, err)
	assertFailedNow(t, mockT, "require.Error")
}

func TestNoError_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.NoError(mockT, errors.New("boom"))
	assertFailedNow(t, mockT, "require.NoError")
}

func TestErrorIs_Failure(t *testing.T) {
	// mockT := &mockTestingT{}
	// err := errors.New("test error")
	// var target error
	// require.ErrorIs(mockT, err, target)
	// assertFailedNow(t, mockT, "require.ErrorIs")
	mockT := &mockTestingT{}
	err := errors.New("test error")
	var target error // target is nil
	require.ErrorIs(mockT, err, target)
	assertFailedNow(t, mockT, "require.ErrorIs")
}

type myCustomErrorForTest struct{ msg string }

func (e *myCustomErrorForTest) Error() string { return e.msg }

func TestErrorAs_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	err := errors.New("standard error") // Use a standard error
	var target *myCustomErrorForTest
	require.ErrorAs(mockT, err, &target)
	assertFailedNow(t, mockT, "require.ErrorAs")
}

func TestNil_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	x := 5
	require.Nil(mockT, &x)
	assertFailedNow(t, mockT, "require.Nil")
}

func TestEmpty_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Empty(mockT, "not empty")
	assertFailedNow(t, mockT, "require.Empty")
}

func TestNotEmpty_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.NotEmpty(mockT, "")
	assertFailedNow(t, mockT, "require.NotEmpty")
}

func TestLen_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Len(mockT, []int{1, 2}, 3) // Actual length is 2, expected 3
	assertFailedNow(t, mockT, "require.Len")
}

func TestZero_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Zero(mockT, 1)
	assertFailedNow(t, mockT, "require.Zero")
}

func TestNotZero_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.NotZero(mockT, 0)
	assertFailedNow(t, mockT, "require.NotZero")
}

func TestGreater_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Greater(mockT, 5, 5) // Not greater
	assertFailedNow(t, mockT, "require.Greater (equal)")

	mockT = &mockTestingT{}
	require.Greater(mockT, 4, 5) // Less
	assertFailedNow(t, mockT, "require.Greater (less)")
}

func TestGreaterOrEqual_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.GreaterOrEqual(mockT, 4, 5) // Less
	assertFailedNow(t, mockT, "require.GreaterOrEqual")
}

func TestLess_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Less(mockT, 5, 5) // Equal
	assertFailedNow(t, mockT, "require.Less (equal)")

	mockT = &mockTestingT{}
	require.Less(mockT, 6, 5) // Greater
	assertFailedNow(t, mockT, "require.Less (greater)")
}

func TestLessOrEqual_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.LessOrEqual(mockT, 6, 5) // Greater
	assertFailedNow(t, mockT, "require.LessOrEqual")
}

func TestPositive_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Positive(mockT, 0) // Zero
	assertFailedNow(t, mockT, "require.Positive (zero)")

	mockT = &mockTestingT{}
	require.Positive(mockT, -1) // Negative
	assertFailedNow(t, mockT, "require.Positive (negative)")
}

func TestNegative_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Negative(mockT, 0) // Zero
	assertFailedNow(t, mockT, "require.Negative (zero)")

	mockT = &mockTestingT{}
	require.Negative(mockT, 1) // Positive
	assertFailedNow(t, mockT, "require.Negative (positive)")
}

func TestElementsMatch_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.ElementsMatch(mockT, []int{1, 2}, []int{1, 3}) // Different elements
	assertFailedNow(t, mockT, "require.ElementsMatch (different elements)")

	mockT = &mockTestingT{}
	require.ElementsMatch(mockT, []int{1, 2, 2}, []int{1, 1, 2}) // Different counts
	assertFailedNow(t, mockT, "require.ElementsMatch (different counts)")
}

func TestNotNil_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	var p *int
	require.NotNil(mockT, p)
	assertFailedNow(t, mockT, "require.NotNil")
}

func TestPanics_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.Panics(mockT, func() {})
	assertFailedNow(t, mockT, "require.Panics")
}

func TestNotPanics_Failure(t *testing.T) {
	mockT := &mockTestingT{}
	require.NotPanics(mockT, func() { panic("oh no") })
	assertFailedNow(t, mockT, "require.NotPanics")
}

func TestTimeAfterFailure(t *testing.T) {
	mockT := &mockTestingT{}
	now := time.Now()
	before := now.Add(-1 * time.Hour)
	require.TimeAfter(mockT, before, now)
	assertFailedNow(t, mockT, "require.TimeAfter")
}

func TestTimeBeforeFailure(t *testing.T) {
	mockT := &mockTestingT{}
	now := time.Now()
	after := now.Add(1 * time.Hour)
	require.TimeBefore(mockT, after, now)
	assertFailedNow(t, mockT, "require.TimeBefore")
}

func TestTimeEqualFailure(t *testing.T) {
	mockT := &mockTestingT{}
	now := time.Now()
	require.TimeEqual(mockT, now, now.Add(1*time.Second))
	assertFailedNow(t, mockT, "require.TimeEqual")
}

func TestWithinTimeFailure(t *testing.T) {
	mockT := &mockTestingT{}
	now := time.Now()
	start := now.Add(1 * time.Hour)
	end := now.Add(2 * time.Hour)
	require.WithinTime(mockT, now, start, end)
	assertFailedNow(t, mockT, "require.WithinTime")
}

func TestWithinDurationFailure(t *testing.T) {
	mockT := &mockTestingT{}
	now := time.Now()
	duration := 1 * time.Hour
	require.WithinDuration(mockT, now, now.Add(2*time.Hour), duration)
	assertFailedNow(t, mockT, "require.WithinDuration")
}

func TestTimeEqualWithPrecisionFailure(t *testing.T) {
	mockT := &mockTestingT{}
	now := time.Now()
	require.TimeEqualWithPrecision(mockT, now, now.Add(1*time.Second), time.Millisecond)
	assertFailedNow(t, mockT, "require.TimeEqualWithPrecision")
}
