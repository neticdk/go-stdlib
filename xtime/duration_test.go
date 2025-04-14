package xtime_test

import (
	"testing"
	"time"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/xtime"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		options  []xtime.FormatOption
		expected string
	}{
		// Basic Tests (Default Options: Max=Day, Min=Sec, Style=Short, No Rounding)
		{
			name:     "zero",
			duration: 0,
			options:  nil,
			expected: "0s",
		},
		{
			name:     "simple seconds",
			duration: 45 * time.Second,
			options:  nil,
			expected: "45s",
		},
		{
			name:     "simple minutes seconds",
			duration: 90 * time.Second,
			options:  nil,
			expected: "1m 30s",
		},
		{
			name:     "simple hours minutes seconds",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  nil,
			expected: "1h 2m 3s",
		},
		{
			name:     "simple days hours",
			duration: 2*xtime.Day + 5*time.Hour,
			options:  nil,
			expected: "2d 5h",
		},
		{
			name:     "simple weeks days (requires WithMaxUnit)",
			duration: 2*xtime.Week + 3*xtime.Day + 1*time.Hour,
			options:  []xtime.FormatOption{xtime.WithMaxUnit(xtime.Week)},
			expected: "2w 3d 1h",
		},
		{
			name:     "negative duration",
			duration: -(90 * time.Second),
			options:  nil,
			expected: "-1m 30s",
		},
		{
			name:     "sub-second default truncation",
			duration: 1*time.Second + 600*time.Millisecond,
			options:  nil,
			expected: "1s", // 600ms truncated
		},

		// Style Tests
		{
			name:     "compact style",
			duration: 90 * time.Second,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleCompact)},
			expected: "1m30s",
		},
		{
			name:     "compact style single component",
			duration: 2 * time.Hour,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleCompact)},
			expected: "2h",
		},
		{
			name:     "long style",
			duration: 90 * time.Second,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleLong)},
			expected: "1 minute, 30 seconds",
		},
		{
			name:     "long style single component (plural)",
			duration: 2 * time.Hour,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleLong)},
			expected: "2 hours",
		},
		{
			name:     "long style single component (singular)",
			duration: 1 * time.Minute,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleLong)},
			expected: "1 minute",
		},
		{
			name:     "long-and style",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleLongAnd)},
			expected: "1 hour, 2 minutes and 3 seconds",
		},
		{
			name:     "long-and style two components",
			duration: 90 * time.Second,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleLongAnd)},
			expected: "1 minute and 30 seconds",
		},

		// MaxComponents Tests
		{
			name:     "max components 1",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  []xtime.FormatOption{xtime.WithMaxComponents(1)},
			expected: "1h",
		},
		{
			name:     "max components 2",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  []xtime.FormatOption{xtime.WithMaxComponents(2)},
			expected: "1h 2m",
		},
		{
			name:     "max components 3",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  []xtime.FormatOption{xtime.WithMaxComponents(3)},
			expected: "1h 2m 3s",
		},
		{
			name:     "max components unlimited (0)",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  []xtime.FormatOption{xtime.WithMaxComponents(0)},
			expected: "1h 2m 3s",
		},
		{
			name:     "max components cuts off before min unit",
			duration: 1*time.Minute + 30*time.Second + 500*time.Millisecond,
			options:  []xtime.FormatOption{xtime.WithMaxComponents(1), xtime.WithMinUnit(time.Millisecond)},
			expected: "1m",
		},

		// Min/Max Unit Tests
		{
			name:     "min unit minutes",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  []xtime.FormatOption{xtime.WithMinUnit(time.Minute)},
			expected: "1h 2m", // 3s is truncated
		},
		{
			name:     "max unit minutes",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second, // 62m 3s
			options:  []xtime.FormatOption{xtime.WithMaxUnit(time.Minute)},
			expected: "62m 3s", // 1h becomes 60m
		},
		{
			name:     "min unit ms",
			duration: 1*time.Second + 500*time.Millisecond + 500*time.Microsecond + 100*time.Nanosecond,
			options:  []xtime.FormatOption{xtime.WithMinUnit(time.Millisecond)},
			expected: "1s 500ms", // µs and ns truncated
		},
		{
			name:     "min unit ns",
			duration: 1*time.Second + 5*time.Nanosecond,
			options:  []xtime.FormatOption{xtime.WithMinUnit(time.Nanosecond)},
			expected: "1s 5ns",
		},
		{
			name:     "duration less than min unit (default sec)",
			duration: 500 * time.Millisecond,
			options:  nil,
			expected: "0s", // Duration is non-zero but less than MinUnit, show 0 of min unit
		},
		{
			name:     "duration less than min unit (explicit ms)",
			duration: 500 * time.Microsecond,
			options:  []xtime.FormatOption{xtime.WithMinUnit(time.Millisecond)},
			expected: "0ms",
		},

		// Rounding Tests
		{
			name:     "rounding disabled (default)",
			duration: 1*time.Second + 600*time.Millisecond,
			options:  nil, // MinUnit=Second
			expected: "1s",
		},
		{
			name:     "rounding enabled, rounds up",
			duration: 1*time.Second + 600*time.Millisecond,
			options:  []xtime.FormatOption{xtime.WithRounding()}, // MinUnit=Second
			expected: "2s",
		},
		{
			name:     "rounding enabled, rounds down",
			duration: 1*time.Second + 400*time.Millisecond,
			options:  []xtime.FormatOption{xtime.WithRounding()}, // MinUnit=Second
			expected: "1s",
		},
		{
			name:     "rounding enabled, exactly half",
			duration: 1*time.Second + 500*time.Millisecond,
			options:  []xtime.FormatOption{xtime.WithRounding()}, // MinUnit=Second
			expected: "2s",                                       // Ties round up
		},
		{
			name:     "rounding with min unit ms",
			duration: 1*time.Millisecond + 600*time.Microsecond,
			options:  []xtime.FormatOption{xtime.WithMinUnit(time.Millisecond), xtime.WithRounding()},
			expected: "2ms",
		},
		{
			name:     "rounding with carry-over seconds to minutes",
			duration: 59*time.Second + 700*time.Millisecond,
			options:  []xtime.FormatOption{xtime.WithRounding()}, // MinUnit=Second
			expected: "1m",
		},
		{
			name:     "rounding with carry-over minutes to hours",
			duration: 59*time.Minute + 45*time.Second, // rounds to 60m
			options:  []xtime.FormatOption{xtime.WithRounding(), xtime.WithMinUnit(time.Minute)},
			expected: "1h",
		},
		{
			name:     "rounding with carry-over multiple levels",
			duration: 1*time.Hour + 59*time.Minute + 59*time.Second + 800*time.Millisecond,
			options:  []xtime.FormatOption{xtime.WithRounding()}, // MinUnit=Second, rounds to 1h 60m 0s -> 2h
			expected: "2h",
		},
		{
			name:     "rounding with max components",
			duration: 1*time.Hour + 59*time.Minute + 30*time.Second, // rounds to 1h 60m -> 2h
			options:  []xtime.FormatOption{xtime.WithRounding(), xtime.WithMinUnit(time.Minute), xtime.WithMaxComponents(1)},
			expected: "2h", // Rounds up and only shows the first component
		},
		{
			name:     "rounding zero component",
			duration: 1*time.Hour + 30*time.Second, // Add 1h 0m 30s
			options:  []xtime.FormatOption{xtime.WithRounding()},
			expected: "1h 30s",
		},

		// Separator/Conjunction Tests
		{
			name:     "custom separator",
			duration: 90 * time.Second,
			options:  []xtime.FormatOption{xtime.WithSeparator(":")},
			expected: "1m:30s",
		},
		{
			name:     "custom conjunction",
			duration: 90 * time.Second,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleLongAnd), xtime.WithConjunction(" plus ")},
			expected: "1 minute plus 30 seconds",
		},
		{
			name:     "long style custom separator",
			duration: 1*time.Hour + 2*time.Minute + 3*time.Second,
			options:  []xtime.FormatOption{xtime.WithStyle(xtime.FormatStyleLong), xtime.WithSeparator(" - ")},
			expected: "1 hour - 2 minutes - 3 seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function under test
			got := xtime.FormatDuration(tt.duration, tt.options...)

			// Assert equality using the assert package
			assert.Equal(t, got, tt.expected)
		})
	}
}

func TestYearsFromDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected float64
	}{
		{"zero", 0, 0.0},
		{"one day", xtime.Day, 1.0 / 365.2425},
		{"365 days", 365 * xtime.Day, 365.0 / 365.2425},
		{"one approx year", xtime.YearApprox, 1.0},
		{"two approx years", 2 * xtime.YearApprox, 2.0},
		{"negative approx year", -1 * xtime.YearApprox, -1.0},
	}

	// Define a small delta for float comparison
	delta := 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := xtime.YearsFromDuration(tt.duration)
			assert.InDelta(t, got, tt.expected, delta)
		})
	}
}

// Test helper to ensure rounding carry-over rebuilds components correctly
// (Needed special test case because the main table rebuild logic had a subtle bug)
func TestFormat_RoundingCarryOverRebuild(t *testing.T) {
	// 59 seconds + 700 ms -> rounds up to 60 seconds -> carries over to 1 minute
	duration := 59*time.Second + 700*time.Millisecond
	opts := []xtime.FormatOption{xtime.WithRounding(), xtime.WithMinUnit(time.Second)}
	expected := "1m" // Should not be "1m 0s"
	got := xtime.FormatDuration(duration, opts...)
	assert.Equal(t, got, expected, "Rounding 59.7s to 1m")

	// 59m 59s + 700ms -> rounds to 59m 60s -> carries to 60m 0s -> carries to 1h 0m 0s
	duration = 59*time.Minute + 59*time.Second + 700*time.Millisecond
	opts = []xtime.FormatOption{xtime.WithRounding(), xtime.WithMinUnit(time.Second)}
	expected = "1h" // Should not be "1h 0m 0s"
	got = xtime.FormatDuration(duration, opts...)
	assert.Equal(t, got, expected, "Rounding 59m59.7s to 1h")
}
