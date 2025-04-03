package assert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

const maxPanicStackDepth = 8192

type testingT interface {
	Log(args ...any)
	Logf(format string, args ...any)
	Errorf(format string, args ...any)
}

type tHelper = interface {
	Helper()
}

// AssertionContext holds information about the context of an assertion
type AssertionContext struct {
	// Call site information
	File     string
	Line     int
	Function string
}

// FileInfo returns a string representation of the file and line information
func (ctx *AssertionContext) FileInfo() string {
	fileInfo := ""
	if ctx.File != "" {
		fileInfo = fmt.Sprintf("%s:%d", ctx.File, ctx.Line)
		if ctx.Function != "" {
			fileInfo += fmt.Sprintf(" (%s)", ctx.Function)
		}
		fileInfo += ": "
	}
	return fileInfo
}

// NewAssertionContext creates a new context by capturing the current call site
func NewAssertionContext(skip int) *AssertionContext {
	ctx := &AssertionContext{}

	// Capture caller information
	if pc, file, line, ok := runtime.Caller(skip); ok {
		ctx.File = file
		ctx.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			ctx.Function = fn.Name()
		}
	}

	return ctx
}

// AssertionError represents an assertion error.
//
// filename.go:line_number (function_name): <Assertion Type>: <Primary Message>
//
//	<Primary Value Label>: <Primary Value Representation>
//	<Comparison Value Label>: <Comparison Value Representation>
//	Error: <Error Message>
//	ExtraValues:
//	  <Extra Value Label>: <Extra Value Representation>
//	Diff:
//	  <Diff Line 1>
//	  <Diff Line 2>
//	Details:
//	  <Detail Line 1>
//	  <Detail Line 2>
type AssertionError struct {
	// Message is the error message.
	Message string
	// PrimaryValue is the primary value being compared.
	PrimaryValue assertionValue
	// ComparisonValue is the value being compared against.
	ComparisonValue assertionValue
	// Error is the (optional) error that occurred during the assertion.
	Error error
	// ExtraValues are values that may be useful for debugging.
	ExtraValues []assertionValue
	// Diff is the (optional) difference between the primary and comparison values.
	Diff string
	// Details are details about the assertion error.
	Details []string
	// Stack trace - only populated when tests fail
	Stack string
}

type assertionValue struct {
	Label string
	Value any
}

// Format handles all components of the message
func (ae *AssertionError) Format(ctx *AssertionContext) string {
	var sb strings.Builder

	// Location header
	sb.WriteString(ctx.FileInfo())

	// Main message
	if ae.Message == "" {
		sb.WriteString("Assertion failed")
	} else {
		sb.WriteString(ae.Message)
	}
	sb.WriteString(":\n")

	// Primary and comparison values
	if ae.PrimaryValue.Label != "" {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", ae.PrimaryValue.Label, formatValue(ae.PrimaryValue.Value)))
	}
	if ae.ComparisonValue.Label != "" {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", ae.ComparisonValue.Label, formatValue(ae.ComparisonValue.Value)))
	}

	// Error if present
	if ae.Error != nil {
		sb.WriteString(fmt.Sprintf("  Error: %v\n", ae.Error))
	}

	// Extra key-value pairs
	for _, kv := range ae.ExtraValues {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", kv.Label, formatValue(kv.Value)))
	}

	// Diff if present
	if ae.Diff != "" {
		sb.WriteString("  Diff:\n")
		for line := range strings.SplitSeq(ae.Diff, "\n") {
			if line != "" {
				sb.WriteString(fmt.Sprintf("    %s\n", line))
			}
		}
	}

	// Details if present
	if len(ae.Details) > 0 {
		sb.WriteString("  Details:\n")
		for _, detail := range ae.Details {
			sb.WriteString(fmt.Sprintf("    %s\n", detail))
		}
	}

	// Stack Trace if present
	if ae.Stack != "" {
		sb.WriteString("  Stack Trace:\n")
		// Add indentation for the stack block
		for line := range strings.SplitSeq(ae.Stack, "\n") {
			if line != "" { // Avoid adding empty lines
				sb.WriteString(fmt.Sprintf("    %s\n", line))
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

// logOptionalMessage logs the optional message and arguments if provided.
// It uses a heuristic (containsVerbs) to detect if the first argument is likely
// a format string intended for t.Logf when multiple arguments are present.
// Otherwise, it uses t.Log to print arguments space-separated.
func logOptionalMessage(t testingT, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if len(msgAndArgs) == 0 {
		return
	}

	if format, ok := msgAndArgs[0].(string); ok && len(msgAndArgs) > 1 && containsVerbs(format) {
		// Likely a format string with arguments, use Logf.
		// fmt will handle literal %% correctly.
		// fmt will produce %!verb(MISSING) or %!(EXTRA) if args mismatch verbs.
		t.Logf(format, msgAndArgs[1:]...)
	} else {
		// Either:
		// - Only one argument was provided.
		// - The first argument was not a string.
		// - The first argument was a string but contained no '%' (or had no args).
		// Treat all arguments as individual values to be logged space-separated.
		t.Log(msgAndArgs...)
	}
}

// containsVerbs checks if a string contains a '%' character
// that is likely part of a format verb (i.e., not '%%').
func containsVerbs(s string) bool {
	for i := range len(s) {
		if s[i] == '%' {
			if i+1 >= len(s) || s[i+1] != '%' {
				return true // Found unescaped '%'
			}
			// Found '%%', skip the second '%'
		}
	}
	return false // No unescaped '%' found
}

// DiffsEnabled indicates whether diffing is enabled for assertion errors.
var DiffsEnabled = true

// shouldGenerateDiff determines if diffing would be useful for these values
func shouldGenerateDiff(got, want any) bool {
	// No diff for nil values
	if got == nil || want == nil {
		return false
	}

	// Get reflect values
	gotVal := reflect.ValueOf(got)
	wantVal := reflect.ValueOf(want)

	// Skip invalid values
	if !gotVal.IsValid() || !wantVal.IsValid() {
		return false
	}

	// Get types
	gotType := gotVal.Type()
	wantType := wantVal.Type()

	// If types don't match, diffing might be misleading
	if gotType != wantType {
		return false
	}

	// Process based on kind
	switch gotVal.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
		return true
	case reflect.Ptr:
		// For pointers, check what they point to
		if gotVal.IsNil() || wantVal.IsNil() {
			return false
		}
		return shouldGenerateDiff(gotVal.Elem().Interface(), wantVal.Elem().Interface())
	default:
		// For simple types (numbers, booleans, etc.), diff isn't so useful
		return false
	}
}

// computeDiff generates a text-based diff between two values
// Uses JSON for comparison and line-by-line diff for presentation.
// It's a basic implementation that doesn't handle deep comparisons.
func computeDiff(got, want any) string {
	// Convert values to pretty-printed JSON for diffing
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error generating diff (failed to marshal 'got' value): %v", err)
	}

	wantJSON, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error generating diff (failed to marshal 'want' value): %v", err)
	}

	// Generate line-by-line diff
	gotLines := strings.Split(string(gotJSON), "\n")
	wantLines := strings.Split(string(wantJSON), "\n")

	return generateLineDiff(gotLines, wantLines)
}

// generateLineDiff creates a line-by-line diff with +/- markers
func generateLineDiff(gotLines, wantLines []string) string {
	maxLines := max(len(gotLines), len(wantLines))
	diffLines := make([]string, 0, maxLines)

	maxLineCount := max(len(gotLines), len(wantLines))

	for i := range maxLineCount {
		var gotLine, wantLine string
		if i < len(gotLines) {
			gotLine = gotLines[i]
		}
		if i < len(wantLines) {
			wantLine = wantLines[i]
		}

		if gotLine == wantLine {
			// Lines match, show without markers
			diffLines = append(diffLines, "  "+gotLine)
		} else {
			// Lines differ
			if gotLine != "" {
				diffLines = append(diffLines, "- "+gotLine)
			}
			if wantLine != "" {
				diffLines = append(diffLines, "+ "+wantLine)
			}
		}
	}

	return strings.Join(diffLines, "\n")
}

// StackTracesEnabled indicates whether stack traces are enabled for assertion errors.
var StackTracesEnabled = false

// reportAssertionError builds and reports a standardized error message
func reportAssertionError(t testingT, ctx *AssertionContext, err *AssertionError) {
	if StackTracesEnabled && err.Stack == "" {
		const size = maxPanicStackDepth
		buf := make([]byte, size)
		n := runtime.Stack(buf, false)
		err.Stack = filterStackTrace(string(buf[:n]))
	}
	t.Errorf("%s", err.Format(ctx))
}

// Prefixes for stack trace lines to be filtered out.
// These correspond to internal testing, runtime, and assert/require functions.
var stackTraceFilterPrefixes = []string{
	"\tgithub.com/neticdk/go-stdlib/assert.",  // Our own assert package
	"\tgithub.com/neticdk/go-stdlib/require.", // Our own require package
	"\ttesting.",       // Go testing framework internals
	"\truntime.goexit", // Standard exit point
	"\truntime.main",   // Main runtime entry
}

// Helper to clean up the stack trace
func filterStackTrace(fullStack string) string {
	lines := strings.Split(fullStack, "\n")
	if len(lines) == 0 {
		return ""
	}

	var filteredLines []string

	// Keep the first line (e.g., "goroutine 1 [running]:")
	if len(lines) > 0 {
		filteredLines = append(filteredLines, lines[0])
	}

	// Process the rest of the lines in pairs (file:line + function)
	for i := 1; i < len(lines)-1; i += 2 {
		fileLine := lines[i]
		funcLine := lines[i+1] // The indented line with the function call

		isInternal := false
		// Check only the function line for internal prefixes
		if strings.HasPrefix(funcLine, "\t") { // Ensure it's an indented function line
			for _, prefix := range stackTraceFilterPrefixes {
				if strings.HasPrefix(funcLine, prefix) {
					isInternal = true
					break
				}
			}
		}

		// If it's not an internal function, keep both lines
		if !isInternal {
			filteredLines = append(filteredLines, fileLine)
			filteredLines = append(filteredLines, funcLine)
		}
	}

	// Handle potential trailing empty line if the input ended with \n
	if len(filteredLines) > 0 && filteredLines[len(filteredLines)-1] == "" {
		filteredLines = filteredLines[:len(filteredLines)-1]
	}

	return strings.Join(filteredLines, "\n")
}

// formatValue formats a value for display in the assertion error message.
// It handles nil values, strings, and uses fmt.Sprintf("%v") for other types.
func formatValue(value any) string {
	if value == nil {
		return "<nil>"
	}

	// Handle string explicitly for cleaner output (no quotes)
	if str, ok := value.(string); ok {
		return str
	}

	// Use reflection to check for zero values for common types
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return "<invalid value>" // Handle invalid reflect.Value
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() == 0 {
			return "0"
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Uint() == 0 {
			return "0"
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() == 0.0 {
			return "0.0"
		}
	case reflect.Bool:
		if !val.Bool() {
			return "false"
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if val.Len() == 0 {
			return "<empty>"
		}
	}

	return fmt.Sprintf("%v", value)
}
