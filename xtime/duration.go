package xtime

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Day  time.Duration = 24 * time.Hour
	Week time.Duration = 7 * Day

	// YearApprox is the average duration of a year in the Gregorian calendar (365.2425 days).
	// Use this for approximations only; it does not account for calendar specifics.
	YearApprox time.Duration = time.Duration(365.2425 * float64(Day))

	// MonthApprox is the average duration of a month in the Gregorian calendar (30.436875 days).
	// Use this for approximations only; it does not account for calendar specifics.
	MonthApprox time.Duration = time.Duration(float64(YearApprox) / 12)
)

// timeUnitDef holds the definition for a single time unit used in formatting.
type timeUnitDef struct {
	// Unit is the duration value corresponding to one of this unit (e.g., xtime.Hour).
	Unit time.Duration
	// NameSingular is the full singular name (e.g., "hour").
	NameSingular string
	// NamePlural is the full plural name (e.g., "hours").
	NamePlural string
	// Symbol is the short symbol (e.g., "h").
	Symbol string
}

var definedUnits = []timeUnitDef{
	{Week, "week", "weeks", "w"},
	{Day, "day", "days", "d"},
	{time.Hour, "hour", "hours", "h"},
	{time.Minute, "minute", "minutes", "m"},
	{time.Second, "second", "seconds", "s"},
	{time.Millisecond, "millisecond", "milliseconds", "ms"},
	{time.Microsecond, "microsecond", "microseconds", "µs"},
	{time.Nanosecond, "nanosecond", "nanoseconds", "ns"},
}

// FormatStyle defines the output style for formatted durations.
type FormatStyle string

const (
	// FormatStyleCompact uses abbreviated units without spaces ("1h5m30s").
	FormatStyleCompact FormatStyle = "compact"
	// FormatStyleShort uses abbreviated units ("1h 5m 30s").
	FormatStyleShort FormatStyle = "short"
	// FormatStyleLong uses full unit names ("1 hour, 5 minutes, 30 seconds").
	FormatStyleLong FormatStyle = "long"
	// FormatStyleLongAnd uses full names with "and" before the last component
	// ("1 hour, 5 minutes and 30 seconds").
	FormatStyleLongAnd FormatStyle = "long-and"
)

// FormatOptions provides configuration for Format.
type FormatOptions struct {
	// MaxUnit is the largest time unit to display (e.g., time.Hour,
	// xtime.Day). Components larger than this will be represented in terms of
	// this unit.
	// Default: xtime.Day.
	MaxUnit time.Duration
	// MinUnit is the smallest time unit to display (e.g., time.Second,
	// time.Millisecond). Any remaining duration smaller than this will be
	// truncated or rounded depending on Rounding.
	// Default: time.Second.
	MinUnit time.Duration
	// Rounding enables rounding of the smallest displayed unit based on the
	// remainder.
	// If false (default), the remainder is truncated.
	Rounding bool
	// MaxComponents limits the maximum number of components displayed (e.g.,
	// 2 // might yield "1h 5m").
	// Set to 0 or negative for unlimited components (down to MinUnit).
	// Default: 0 (unlimited).
	MaxComponents int
	// Style determines the format of unit names (compact, short, long,
	// long-and).
	// Default: FormatStyleShort.
	Style FormatStyle
	// Separator is the string used between components (ignored if only one
	// component).
	// Default: ", " for long styles, " " for short style, "" for compact style.
	Separator string
	// Conjunction is the string used before the last component in "long-and"
	// style.
	// Default: " and ".
	Conjunction string
}

// DefaultFormatOptions creates options with default values.
func DefaultFormatOptions() FormatOptions {
	return FormatOptions{
		MaxUnit:       Day,
		MinUnit:       time.Second,
		Rounding:      false,
		MaxComponents: 0, // Unlimited
		Style:         FormatStyleShort,
		Separator:     " ", // Default for short
		Conjunction:   " and ",
	}
}

// FormatOption is a function type for setting format options.
type FormatOption func(*FormatOptions)

// WithMaxUnit sets the largest unit to display.
func WithMaxUnit(unit time.Duration) FormatOption {
	return func(o *FormatOptions) {
		if unit >= time.Nanosecond {
			o.MaxUnit = unit
		}
	}
}

// WithMinUnit sets the smallest unit to display.
func WithMinUnit(unit time.Duration) FormatOption {
	return func(o *FormatOptions) {
		if unit >= time.Nanosecond {
			o.MinUnit = unit
		}
	}
}

// WithRounding enables rounding of the MinUnit.
func WithRounding() FormatOption {
	return func(o *FormatOptions) {
		o.Rounding = true
	}
}

// WithMaxComponents sets the max number of components. 0 means unlimited.
func WithMaxComponents(p int) FormatOption {
	return func(o *FormatOptions) {
		o.MaxComponents = p
	}
}

// WithStyle sets the output style (short, long, long-and).
func WithStyle(style FormatStyle) FormatOption {
	return func(o *FormatOptions) {
		o.Style = style
		// Adjust default separator based on style if not explicitly set later
		switch style {
		case FormatStyleShort:
			o.Separator = " "
		case FormatStyleCompact:
			o.Separator = ""
		default:
			o.Separator = ", "
		}
	}
}

// WithSeparator sets the separator string.
func WithSeparator(sep string) FormatOption {
	return func(o *FormatOptions) {
		o.Separator = sep
	}
}

// WithConjunction sets the conjunction string for "long-and" style.
func WithConjunction(conj string) FormatOption {
	return func(o *FormatOptions) {
		o.Conjunction = conj
	}
}

// componentResult holds the calculated quantity and unit definition for one part.
type componentResult struct {
	Quantity int64
	UnitDef  timeUnitDef
}

// FormatDuration formats a duration into a human-readable string
// composed of multiple time units (e.g., "1h 5m 30s", "2 days, 3 hours").
//
// Use FormatOptions functions (WithMaxUnit, WithMinUnit, WithRounding,
// WithMaxComponents, WithStyle, WithSeparator, WithConjunction) to customize output.
//
// Example:
//
//	FormatDuration(90*time.Minute, WithStyle(FormatStyleLong)) // "1 hour, 30 minutes"
//	FormatDuration(3723*time.Second, WithMinUnit(time.Second), WithMaxComponents(2)) // "1h 2m"
//	FormatDuration(time.Second + 600*time.Millisecond, WithRounding()) // "2s"
func FormatDuration(d time.Duration, opts ...FormatOption) string {
	options := DefaultFormatOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if d == 0 {
		minUnitDef, foundMin := findUnitDef(options.MinUnit)
		if !foundMin {
			minUnitDef, _ = findUnitDef(time.Second)
		}
		return formatComponent(0, minUnitDef, options.Style)
	}

	isNegative := d < 0
	if isNegative {
		d = -d
	}

	// Simple rounding: Add half of the minimum unit before processing
	if options.Rounding && options.MinUnit > 0 {
		d += options.MinUnit / 2
	}

	// Calculate components
	results := []componentResult{}
	remaining := d

	minUnitDef, foundMin := findUnitDef(options.MinUnit)
	if !foundMin {
		minUnitDef, _ = findUnitDef(time.Second) // Fallback
	}

	for _, unitDef := range definedUnits {
		if unitDef.Unit > options.MaxUnit {
			continue
		}
		// Stop if we hit the component limit *before* processing this unit
		if options.MaxComponents > 0 && len(results) >= options.MaxComponents {
			break
		}
		if unitDef.Unit < options.MinUnit {
			break
		}

		if remaining >= unitDef.Unit {
			quantity := remaining / unitDef.Unit
			// Note: 'remaining' here still holds the full remainder including
			// sub-MinUnit parts because we modified 'd' upfront.
			currentUnitRemainder := remaining % unitDef.Unit

			if quantity > 0 {
				results = append(results, componentResult{
					Quantity: int64(quantity),
					UnitDef:  unitDef,
				})
				// Update remaining for the *next* iteration
				remaining = currentUnitRemainder
			}
		}
	}

	// Handle edge case where rounding resulted in zero components, but
	// original > 0 or where the original duration was less than MinUnit but
	// rounded up.
	if len(results) == 0 {
		// If the (potentially rounded) duration is >= MinUnit, display 1 MinUnit
		if d >= options.MinUnit {
			return formatComponent(1, minUnitDef, options.Style)
		}
		// Otherwise, it was < MinUnit and didn't round up, display 0 MinUnit
		return formatComponent(0, minUnitDef, options.Style)
	}

	// Format components
	components := make([]string, 0, len(results))
	for _, res := range results {
		components = append(
			components,
			formatComponent(res.Quantity, res.UnitDef, options.Style))
	}

	// Join components
	result := joinComponents(
		components,
		options.Style,
		options.Separator,
		options.Conjunction)

	if isNegative {
		return "-" + result
	}

	return result
}

// formatComponent formats a single quantity and unit definition based on style.
func formatComponent(quantity int64, unitDef timeUnitDef, style FormatStyle) string {
	s := strconv.FormatInt(quantity, 10)

	switch style {
	case FormatStyleLong, FormatStyleLongAnd:
		unitName := unitDef.NameSingular
		if quantity != 1 {
			unitName = unitDef.NamePlural
		}
		return s + " " + unitName
	default: // FormatStyleCompact, FormatStyleShort
		return s + unitDef.Symbol
	}
}

// joinComponents joins the formatted string components based on style and separators.
func joinComponents(components []string, style FormatStyle, separator, conjunction string) string {
	count := len(components)
	if count == 0 {
		return ""
	}
	if count == 1 {
		return components[0]
	}

	if style == FormatStyleLongAnd && count > 1 {
		allButLast := strings.Join(components[:count-1], separator)
		return allButLast + conjunction + components[count-1]
	}

	return strings.Join(components, separator)
}

// findUnitDef searches definedUnits for a specific duration value.
func findUnitDef(unit time.Duration) (timeUnitDef, bool) {
	for _, def := range definedUnits {
		if def.Unit == unit {
			return def, true
		}
	}
	return timeUnitDef{}, false
}

// YearsFromDuration converts a duration into an approximate number of years.
// It calculates this based on the average length of a year in the
// Gregorian calendar (365.2425 days).
//
// WARNING: This function provides an estimation based on duration only.
// It does not account for specific calendar start/end dates, leap year
// occurrences within a specific period, or time zones. For calendar-accurate
// differences involving years and months, use functions operating on time.Time
// values.
func YearsFromDuration(d time.Duration) float64 {
	return float64(d) / float64(YearApprox)
}

// init ensures definedUnits is sorted correctly on package load
func init() {
	sort.SliceStable(definedUnits, func(i, j int) bool {
		return definedUnits[i].Unit > definedUnits[j].Unit
	})
}
