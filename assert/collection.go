package assert

import (
	"fmt"
	"reflect"
	"strings"
)

// Contains asserts that a collection contains a specific element/substring.
// Supports:
// - string: checks for substring presence.
// - slice/array: checks for element presence using the internal 'equal' comparison.
// - map: checks for value presence using the internal 'equal' comparison (unlike ContainsKey which checks for keys).
func Contains(t testingT, collection any, element any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	found, err := containsElementInternal(collection, element)
	if err != nil {
		reportCollectionError(t, ctx, "Error checking contains", collection, element, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	if !found {
		colKind := reflect.ValueOf(collection).Kind()
		errMsg := "Collection does not contain element"
		switch colKind {
		case reflect.String:
			errMsg = "String does not contain substring"
		case reflect.Map:
			errMsg = "Map does not contain value"
		}
		reportCollectionError(t, ctx, errMsg, collection, element, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NotContains asserts that a collection does not contain a specific element/substring.
// Opposite of Contains.
// Supports the same types as Contains.
func NotContains(t testingT, collection any, element any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	found, err := containsElementInternal(collection, element)
	if err != nil {
		reportCollectionError(t, ctx, "Error checking not-contains", collection, element, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	if found {
		colKind := reflect.ValueOf(collection).Kind()
		errMsg := "Collection should not contain element"
		switch colKind {
		case reflect.String:
			errMsg = "String should not contain substring"
		case reflect.Map:
			errMsg = "Map should not contain value"
		}
		reportCollectionError(t, ctx, errMsg, collection, element, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// ContainsKey asserts that a map contains a specific key.
func ContainsKey[K comparable, V any](t testingT, m map[K]V, key K, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	_, ok := m[key]
	if !ok {
		reportCollectionError(t, ctx, "Map does not contain key", m, key, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NotContainsKey asserts that a map does not contain a specific key.
// Opposite of ContainsKey.
func NotContainsKey[K comparable, V any](t testingT, m map[K]V, key K, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	_, ok := m[key]
	if ok {
		reportCollectionError(t, ctx, "Map should not contain key", m, key, nil)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// Empty asserts that a value is considered "empty".
// True for: nil pointers/interfaces/slices/maps/channels/funcs,
// zero values (0, "", false, zero structs), and zero-length
// slices/maps/arrays/strings/channels.
func Empty(t testingT, data any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if !isEmptyInternal(data) {
		reportEqualityError(t, ctx, data, "<empty/zero value>")
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// NotEmpty asserts that a value is not considered "empty".
// Opposite of Empty.
func NotEmpty(t testingT, data any, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	if isEmptyInternal(data) {
		reportEqualityError(t, ctx, data, "<non-empty/zero value>")
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// Len asserts that a collection (slice, map, array, string, channel) has a specific length.
func Len(t testingT, data any, expectedLen int, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	ctx := NewAssertionContext(1)

	v := reflect.ValueOf(data)
	var actualLen int

	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.String, reflect.Chan, reflect.Array:
		actualLen = v.Len()
	default:
		err := &AssertionError{
			Message: fmt.Sprintf("Cannot get length of type %T, value: %#v", data, data),
		}
		reportAssertionError(t, ctx, err)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}

	if actualLen != expectedLen {
		details := []string{"Length mismatch"}
		reportEqualityError(t, ctx, actualLen, expectedLen, details...)
		logOptionalMessage(t, msgAndArgs...)
		return false
	}
	return true
}

// isEmptyInternal checks if a value is nil, zero, or empty length. Used by
// Empty/NotEmpty.
func isEmptyInternal(data any) bool {
	if data == nil {
		return true
	}
	v := reflect.ValueOf(data)
	// Use IsZero first as it covers many cases (nil pointers, zero
	// numbers/strings/bools/structs, nil maps/slices/chans/funcs)
	// Need to check IsValid because IsZero panics on zero Value (e.g. untyped
	// nil handled above)
	if !v.IsValid() {
		return true // Should be caught by data == nil, but defensive check
	}
	// Check IsZero BEFORE length check, because IsZero is true for nil
	// maps/slices etc.
	if v.IsZero() {
		return true
	}

	// If not zero (e.g., non-nil but potentially empty
	// map/slice/channel/string/array), check length.
	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.String, reflect.Array:
		return v.Len() == 0
	}

	return false
}

// containsElementInternal is used by Contains/NotContains.
func containsElementInternal(collection any, element any) (bool, error) {
	colVal := reflect.ValueOf(collection)
	if !colVal.IsValid() {
		return false, fmt.Errorf("collection is invalid (e.g. untyped nil)")
	}

	switch colVal.Kind() {
	case reflect.Slice, reflect.Array:
		for i := range colVal.Len() {
			if equalInternal(colVal.Index(i).Interface(), element) {
				return true, nil
			}
		}
		return false, nil
	case reflect.Map:
		// Checks if element is present as a *value* in the map
		mapIter := colVal.MapRange()
		for mapIter.Next() {
			if equalInternal(mapIter.Value().Interface(), element) {
				return true, nil
			}
		}
		return false, nil
	case reflect.String:
		// For string, element must also be a string or convertible for
		// substring check
		elemVal := reflect.ValueOf(element)
		if elemVal.Kind() != reflect.String {
			return false, fmt.Errorf("cannot check for non-string element (%T) in a string", element)
		}
		return strings.Contains(colVal.String(), elemVal.String()), nil
	default:
		return false, fmt.Errorf("type %T is not searchable for elements (only slice, array, map, string)", collection)
	}
}

// reportCollectionError reports errors for collection assertions (Contains, NotContains, etc.)
// nolint: unparam
func reportCollectionError(t testingT, ctx *AssertionContext, errType string, collection, element any, err error, details ...string) { //revive:disable-line:argument-limit
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	assertErr := &AssertionError{
		Message: errType,
		PrimaryValue: assertionValue{
			Label: "Collection",
			Value: collection,
		},
		ComparisonValue: assertionValue{
			Label: "Element",
			Value: element,
		},
	}

	if err != nil {
		assertErr.Error = err
	}

	if len(details) > 0 {
		assertErr.Details = details
	}

	reportAssertionError(t, ctx, assertErr)
}
