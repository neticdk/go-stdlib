package xtime_test

import (
	"testing"
	"time"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/xtime"
)

func TestParseDuration(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		// --- Valid Cases ---
		{"single number seconds int", "10", 10 * time.Second, false},
		{"single number seconds float", "15.5", time.Duration(15.5 * float64(time.Second)), false},
		{"single number zero", "0", 0, false},
		{"simple hour", "1h", time.Hour, false},
		{"simple minute", "5m", 5 * time.Minute, false},
		{"simple second", "30s", 30 * time.Second, false},
		{"simple day", "2d", 2 * xtime.Day, false},
		{"simple week", "3w", 3 * xtime.Week, false},
		{"simple millisecond", "500ms", 500 * time.Millisecond, false},
		{"simple microsecond", "250us", 250 * time.Microsecond, false},
		{"simple microsecond mu", "250μs", 250 * time.Microsecond, false},
		{"simple nanosecond", "100ns", 100 * time.Nanosecond, false},
		{"float value", "1.5h", time.Duration(1.5 * float64(time.Hour)), false},
		{"combined no space", "1h30m", time.Hour + 30*time.Minute, false},
		{"combined with space", "1h 30m", time.Hour + 30*time.Minute, false},
		{"multiple components", "1h 30m 15s", time.Hour + 30*time.Minute + 15*time.Second, false},
		{"multiple components no spaces", "1h30m15s", time.Hour + 30*time.Minute + 15*time.Second, false},
		{"mixed units", "1d 12h 30m 5s", xtime.Day + 12*time.Hour + 30*time.Minute + 5*time.Second, false},
		{"abbreviations sec", "10sec", 10 * time.Second, false},
		{"abbreviations secs", "15secs", 15 * time.Second, false},
		{"abbreviations min", "5min", 5 * time.Minute, false},
		{"abbreviations mins", "2mins", 2 * time.Minute, false},
		{"abbreviations hr", "3hr", 3 * time.Hour, false},
		{"long names singular", "1 day 2 hour 3 minute 4 second", xtime.Day + 2*time.Hour + 3*time.Minute + 4*time.Second, false},
		{"long names plural", "2 days 3 hours 4 minutes 5 seconds", 2*xtime.Day + 3*time.Hour + 4*time.Minute + 5*time.Second, false},
		{"leading space", " 1h", time.Hour, false},
		{"trailing space", "1h ", time.Hour, false},
		{"multiple spaces", "1h  30m", time.Hour + 30*time.Minute, false},
		{"comma separator", "1h,30m", time.Hour + 30*time.Minute, false},
		{"comma separator with space", "1 day, 12 hours", xtime.Day + 12*time.Hour, false},
		{"'and ' separator", "1h and 30m", time.Hour + 30*time.Minute, false}, // Relies on clean removing "and "
		{"approximate month", "1mo", xtime.MonthApprox, false},
		{"approximate year", "2y", 2 * xtime.YearApprox, false},
		{"mixed approx and fixed", "1y 6mo 2w 3d", xtime.YearApprox + 6*xtime.MonthApprox + 2*xtime.Week + 3*xtime.Day, false},
		{"leading decimal use", ".h", time.Duration(0 * float64(time.Second)), false},
		{"leading decimal number", ".5s", time.Duration(0.5 * float64(time.Second)), false},
		{"leading decimal combined", ".5h30m", time.Duration(0.5*float64(time.Hour)) + 30*time.Minute, false},
		{"trailing decimal number", "1.s", time.Second, false}, // strconv.ParseFloat handles "1." as 1.0
		{"trailing decimal combined", "1.h30m", time.Hour + 30*time.Minute, false},
		{"zero value units", "0h 0m 0s", 0, false},
		{"large value", "1000d", 1000 * xtime.Day, false},

		// --- Invalid Cases (Expect Error) ---
		{"empty string", "", 0, true},
		{"whitespace only", "   ", 0, true},
		{"just comma", ",", 0, true},
		{"just and", "and ", 0, true},
		{"single non-number", "h", 0, true},
		{"single invalid token", "abc", 0, true},
		{"odd number of tokens", "1h 30", 0, true},
		{"non-number in value pos 1", "h 1", 0, true},
		{"non-number in value pos 2", "1h m 30s", 0, true},
		{"unknown unit", "1h 30foo", 0, true},
		{"unknown unit short", "10 quark", 0, true},
		{"invalid character embedded", "1h$30m", 0, true},
		{"multiple decimals", "1.2.3h", 0, true},
		{"unit value reversed", "h1", 0, true},
		{"value value", "1 2", 0, true},
		{"unit unit", "h m", 0, true},
		{"number touching invalid", "1h$", 0, true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xtime.ParseDuration(tt.input)

			if tt.wantErr {
				assert.Error(t, err, "ParseDuration/%s", tt.name)
				var parseErr *xtime.DurationParseError
				assert.ErrorAs(t, err, &parseErr, "ParseDuration/%s", tt.name)
			} else {
				assert.NoError(t, err, "ParseDuration/%s", tt.name)
				assert.Equal(t, got, tt.want, "ParseDuration/%s", tt.name)
			}
		})
	}
}
