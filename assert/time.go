package assert

import (
	"fmt"
	"strings"
	"time"
)

// TimeAfter asserts that a time is after a threshold time
func TimeAfter(t testingT, got, threshold time.Time, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !got.After(threshold) {
		var details []string
		// Add specific explanation based on the comparison result
		if got.Equal(threshold) {
			details = append(details, "Times are exactly equal")
		} else {
			details = append(details, fmt.Sprintf("Got time is %v earlier than threshold", threshold.Sub(got)))
		}

		reportTimeComparisonError(t, ctx, "Time is not after threshold", got, threshold, "after", details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// TimeBefore asserts that a time is before a threshold time
func TimeBefore(t testingT, got, threshold time.Time, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !got.Before(threshold) {
		var details []string
		//
		// Add specific explanation based on the comparison result
		if got.Equal(threshold) {
			details = append(details, "Times are exactly equal")
		} else {
			details = append(details, fmt.Sprintf("Got time is %v later than threshold", got.Sub(threshold)))
		}

		reportTimeComparisonError(t, ctx, "Time is not before threshold", got, threshold, "before", details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// TimeEqual asserts that two times represent the same instant
func TimeEqual(t testingT, got, want time.Time, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !got.Equal(want) {
		var details []string

		// Check if they're the same time but in different locations
		gotUTC, wantUTC := got.UTC(), want.UTC()
		if gotUTC.Equal(wantUTC) {
			details = append(details, "Times represent the same instant but have different time zones")
		}

		reportTimeComparisonError(t, ctx, "Times are not equal", got, want, "equal", details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// WithinDuration asserts that two times are within a certain duration of each other
func WithinDuration(t testingT, got, want time.Time, delta time.Duration, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	diff := got.Sub(want)
	if diff < 0 {
		diff = -diff
	}

	if diff > delta {
		details := []string{
			fmt.Sprintf("Maximum allowed difference: %v", delta),
			fmt.Sprintf("Actual difference: %v", diff),
			fmt.Sprintf("Difference exceeds allowed delta by: %v", diff-delta),
		}

		reportTimeComparisonError(t, ctx, "Times are not within duration of each other", got, want, "within", details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// TimeEqualWithPrecision asserts that two times are equal within a certain precision
func TimeEqualWithPrecision(t testingT, got, want time.Time, precision time.Duration, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	// Truncate both times to the specified precision
	truncatedGot := got.Truncate(precision)
	truncatedWant := want.Truncate(precision)

	if !truncatedGot.Equal(truncatedWant) {
		details := []string{
			fmt.Sprintf("Using precision: %v", precision),
			fmt.Sprintf("Truncated got: %s", truncatedGot.Format(time.RFC3339Nano)),
			fmt.Sprintf("Truncated want: %s", truncatedWant.Format(time.RFC3339Nano)),
		}

		reportTimeComparisonError(t, ctx, "Times are not equal with the specified precision", got, want, "equal", details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// WithinTime asserts that a time is within a given time window
func WithinTime(t testingT, got time.Time, start, end time.Time, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if got.Before(start) || got.After(end) {
		var details []string
		details = append(details, fmt.Sprintf("Expected time between: %s and %s",
			start.Format(time.RFC3339), end.Format(time.RFC3339)))

		windowDuration := end.Sub(start)
		details = append(details, fmt.Sprintf("Time window duration: %v", windowDuration))

		// Add details about why it failed
		if got.Before(start) {
			details = append(details,
				"Too early: time is before the start of the window",
				fmt.Sprintf("Time is %v before window start", start.Sub(got)))

			// Use start time as the reference for the report
			reportTimeComparisonError(t, ctx, "Time is outside allowed window", got, start, "within", details...)
		} else {
			details = append(details,
				"Too late: time is after the end of the window",
				fmt.Sprintf("Time is %v after window end", got.Sub(end)))

			// Use end time as the reference for the report
			reportTimeComparisonError(t, ctx, "Time is outside allowed window", got, end, "within", details...)
		}

		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	return true
}

// reportTimeComparisonError reports errors for time comparison assertions
func reportTimeComparisonError(t testingT, ctx *AssertionContext, message string, got, reference time.Time, relationship string, details ...string) { //revive:disable-line:argument-limit
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	// Start building parts of the message
	parts := []string{
		fmt.Sprintf("Got time: %s", got.Format(time.RFC3339Nano)),
	}

	// Add appropriate label based on relationship
	var referenceLabel string
	switch relationship {
	case "before":
		referenceLabel = "Threshold time"
	case "after":
		referenceLabel = "Threshold time"
	case "equal":
		referenceLabel = "Want time"
	default:
		referenceLabel = "Reference time"
	}

	parts = append(parts, fmt.Sprintf("%s: %s", referenceLabel, reference.Format(time.RFC3339Nano)))

	// Calculate and add time difference
	diff := got.Sub(reference)
	switch {
	case diff < 0:
		parts = append(parts, fmt.Sprintf("Time difference: %v earlier", -diff))
	case diff > 0:
		parts = append(parts, fmt.Sprintf("Time difference: %v later", diff))
	default:
		parts = append(parts, "Time difference: 0 (times are identical)")
	}

	// Add timezone information if they differ
	if got.Location().String() != reference.Location().String() {
		parts = append(parts,
			fmt.Sprintf("Got timezone: %s", got.Location()),
			fmt.Sprintf("%s timezone: %s", referenceLabel, reference.Location()))
	}

	// Add any additional details
	if len(details) > 0 {
		parts = append(parts, "Details:")
		for _, detail := range details {
			parts = append(parts, "  "+detail) // Add indentation to details
		}
	}

	// Format the message with consistent indentation
	messageBody := strings.Join(parts, "\n  ")
	fullMessage := fmt.Sprintf("%s:\n  %s", message, messageBody)

	// Report the error
	t.Errorf("%s%s", ctx.FileInfo(), fullMessage)
}
