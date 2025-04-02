package assert

import (
	"fmt"
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

		reportTimeComparisonError(t, ctx, "Time is not after threshold", got, threshold, "after", nil, details...)
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

		reportTimeComparisonError(t, ctx, "Time is not before threshold", got, threshold, "before", nil, details...)
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

		reportTimeComparisonError(t, ctx, "Times are not equal", got, want, "equal", nil, details...)
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

		reportTimeComparisonError(t, ctx, "Times are not within duration of each other", got, want, "within", nil, details...)
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

		reportTimeComparisonError(t, ctx, "Times are not equal with the specified precision", got, want, "equal", nil, details...)
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
			reportTimeComparisonError(t, ctx, "Time is outside allowed window", got, start, "within", nil, details...)
		} else {
			details = append(details,
				"Too late: time is after the end of the window",
				fmt.Sprintf("Time is %v after window end", got.Sub(end)))

			// Use end time as the reference for the report
			reportTimeComparisonError(t, ctx, "Time is outside allowed window", got, end, "within", nil, details...)
		}

		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	return true
}

// reportTimeComparisonError reports errors for time comparison assertions
// nolint: unparam
func reportTimeComparisonError(t testingT, ctx *AssertionContext, message string, got, reference time.Time, relationship string, err error, details ...string) { //revive:disable-line:argument-limit
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	// Determine the label for the reference value based on relationship
	referenceLabel := "Reference time"
	switch relationship {
	case "before":
		referenceLabel = "Threshold time"
	case "after":
		referenceLabel = "Threshold time"
	case "equal":
		referenceLabel = "Want time"
	}

	// Create the assertion error
	assertErr := &AssertionError{
		Message: message,
		PrimaryValue: assertionValue{
			Label: "Got time",
			Value: got,
		},
		ComparisonValue: assertionValue{
			Label: referenceLabel,
			Value: reference,
		},
	}

	// Add time difference as an extra value
	diff := got.Sub(reference)
	var diffDesc string
	switch {
	case diff < 0:
		diffDesc = fmt.Sprintf("%v earlier", -diff)
	case diff > 0:
		diffDesc = fmt.Sprintf("%v later", diff)
	default:
		diffDesc = "0 (times are identical)"
	}
	assertErr.ExtraValues = append(assertErr.ExtraValues, assertionValue{
		Label: "Time difference",
		Value: diffDesc,
	})

	// Add timezone information if they differ
	if got.Location().String() != reference.Location().String() {
		assertErr.ExtraValues = append(assertErr.ExtraValues,
			assertionValue{
				Label: "Got timezone",
				Value: got.Location().String(),
			},
			assertionValue{
				Label: fmt.Sprintf("%s timezone", referenceLabel),
				Value: reference.Location().String(),
			},
		)
	}

	// Add error if present
	if err != nil {
		assertErr.Error = err
	}

	// Add details
	assertErr.Details = details

	// Report the error
	reportAssertionError(t, ctx, assertErr)
}
