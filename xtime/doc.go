// Package xtime provides functionality for working with time.
//
// # Duration Formatting (FormatDuration)
//
// The FormatDuration function converts a time.Duration into a string composed
// of multiple time units (e.g., "1h 5m 30s", "2 days, 3 hours"). This offers
// more detailed breakdowns than the standard time.Duration.String() method,
// which typically scales to the largest appropriate single unit (e.g.,
// "1m30.5s").
//
//	duration := 2*xtime.Day + 3*time.Hour + 15*time.Minute + 30*time.Second +
//				500*time.Millisecond
//
// Default formatting (short style, max unit Day, min unit Second, no
// rounding):
//
//	fmt.Println(xtime.FormatDuration(duration)) // Output: 2d 3h 15m 30s
//
// ## Custom formatting examples
//
// Long style, max 2 components, rounding enabled:
//
//	fmt.Println(xtime.FormatDuration(duration,
//		xtime.WithStyle(xtime.FormatStyleLong),
//		xtime.WithMaxComponents(2),
//		xtime.WithRounding(),
//	)) // Output: 2 days, 3 hours
//
// Compact style:
//
//	fmt.Println(xtime.FormatDuration(duration,
//		xtime.WithStyle(xtime.FormatStyleCompact),
//	)) // Output: 2d3h
//
// Short style, rounding enabled, only seconds and smaller displayed:
//
//	fmt.Println(xtime.FormatDuration(time.Second+600*time.Millisecond,
//		xtime.WithRounding(),
//		xtime.WithMinUnit(time.Millisecond),
//	)) // Output: 2s (Original: 1s 600ms, rounded up)
//
// Formatting Options (Functional Options Pattern):
//   - WithMaxUnit(unit time.Duration): Sets the largest unit for decomposition
//     (default: Day).
//   - WithMinUnit(unit time.Duration): Sets the smallest unit to display
//     (default: Second). Remainder is truncated or rounded.
//   - WithRounding(): Enables rounding of the MinUnit based on the remainder.
//     The duration is adjusted by adding half of MinUnit before decomposition.
//   - WithoutRounding(): Disables rounding (default). Remainder is truncated.
//   - WithMaxComponents(n int): Limits output to at most 'n' components
//     (default: 0 = unlimited).
//   - WithStyle(style FormatStyle): Sets output style (short, long, long-and).
//   - WithSeparator(sep string): Custom separator between components (default
//     depends on style).
//   - WithConjunction(conj string): Custom conjunction (" and " by default)
//     used before the last component in "long-and" style.
//
// # Duration Parsing (ParseDuration)
//
// The ParseDuration function converts a human-readable string representation
// into a time.Duration value. It accepts various formats, including combined
// units and common abbreviations.
//
//	d1, err := xtime.ParseDuration("1h 30m 15s")
//	d2, err := xtime.ParseDuration("1.5hours 10sec") // Combined number/unit and spaces work
//	d3, err := xtime.ParseDuration("10 years, 2 months, 5 days") // Approximate units allowed
//	d4, err := xtime.ParseDuration("3d12h")
//
// Input String Processing:
//  1. Cleaning: Leading/trailing whitespace is trimmed, "and " sequences are
//     removed, and commas are replaced with spaces.
//  2. Tokenization: The cleaned string is split by spaces. Tokens containing
//     both numbers and letters (e.g., "10years", "1h30m", "h1") are further
//     split into number and unit parts (e.g., "10", "years", "1", "h", "30",
//     "m", "h", "1").
//  3. Parsing:
//     - If only one token results and it's a valid number, it's interpreted as
//     seconds.
//     - Otherwise, tokens are processed in pairs (value, unit). The value must
//     be a number (integer or float), and the unit must be one of the
//     recognized unit strings (e.g., "h", "hour", "hours", "d", "day", "days",
//     "mo", "month", "y", "year").
//     - Parsing fails if the token sequence is invalid (e.g., odd number of
//     tokens, non-number where value is expected, unknown unit).
//
// Error Handling:
//   - Returns a specific error type (*DurationParseError) containing details
//     about the failure, including the original input, problematic token, and index.
//
// # Units and Approximations
//
// The package defines standard fixed-duration units (Week, Day) and also
// provides approximate average durations for Month (MonthApprox) and Year
// (YearApprox) based on the Gregorian calendar average (365.2425 days/year).
//
// The YearsFromDuration(d time.Duration) function converts a duration to an
// approximate number of years using YearApprox. Useful for rough estimations
// only.
//
// Note on Units Discrepancy: ParseDuration can parse approximate units like
// "month" (mo) and "year" (y) based on the average durations (MonthApprox,
// YearApprox). However, FormatDuration does not format durations using these
// approximate units; it will decompose them into weeks, days, etc., for more
// precise representation based on the fixed time.Duration value.
//
// Note: For calendar-accurate calculations involving months and years (which
// vary in length), always operate on time.Time values using functions like
// time.AddDate and time.Sub, rather than relying solely on time.Duration
// arithmetic.
package xtime
