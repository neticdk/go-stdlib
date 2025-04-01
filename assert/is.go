package assert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

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

	// Stack trace (can be lazy-loaded if needed)
	Stack []byte
}

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

// logOptionalMessage logs the optional message and arguments if provided.
// func logOptionalMessage(t *testing.T, msgAndArgs ...any) {
func logOptionalMessage(t testingT, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if len(msgAndArgs) > 0 {
		format, ok := msgAndArgs[0].(string)
		if ok {
			// Check if the format string contains format verbs
			if containsVerbs(format) && len(msgAndArgs) > 1 {
				t.Logf(format, msgAndArgs[1:]...)
			} else {
				t.Log(msgAndArgs...)
			}
		} else {
			t.Log(msgAndArgs...)
		}
	}
}

// containsVerbs checks if a format string contains format verbs.
// Rudementary implementation.
func containsVerbs(format string) bool {
	return strings.ContainsRune(format, '%')
}

// errorMessage formats an assertion error message with context information
func errorMessage(ctx *AssertionContext, message string, args ...any) string {
	fileInfo := ctx.FileInfo()
	formattedMessage := fmt.Sprintf(message, args...)
	return fileInfo + formattedMessage
}

// reportError reports an assertion error with context
func reportError(t testingT, ctx *AssertionContext, message string, args ...any) {
	t.Errorf("%s", errorMessage(ctx, message, args...))
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

	// Get types, handling the case where one or both might be nil
	var gotType, wantType reflect.Type
	if gotVal.IsValid() {
		gotType = gotVal.Type()
	}
	if wantVal.IsValid() {
		wantType = wantVal.Type()
	}

	// If types don't match, diffing might be misleading
	if gotType != wantType {
		return false
	}

	// Only generate diffs for complex types where it's useful
	if gotVal.IsValid() {
		switch gotVal.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			return true
		case reflect.Ptr:
			// For pointers, check what they point to
			if gotVal.IsNil() || wantVal.IsNil() {
				return false
			}
			return shouldGenerateDiff(gotVal.Elem().Interface(), wantVal.Elem().Interface())
		}
	}

	// For simple types (numbers, booleans, etc.), diff isn't very useful
	return false
}

// computeDiff generates a text-based diff between two values
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
// It's very simple
func generateLineDiff(gotLines, wantLines []string) string {
	var diffLines []string
	maxLines := max(len(gotLines), len(wantLines))

	for i := range maxLines {
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
