package xstructs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/neticdk/go-stdlib/xslices"
)

// tagCategories is a list of tag categories to check for.
var tagCategories = []string{"json", "yaml"}

// ToMap converts a struct or map to a map[string]any.
// It handles nested structs, maps, and slices.
// It uses the "json" and "yaml" tags to determine the key names.
// It respects the `omitempty` tag for fields.
// It respects the `inline` tag for nested structs.
// It respects the `-` tag to omit fields.
//
// If the input is nil, it returns nil.
// If the input is not a struct or map, it returns an error.
func ToMap(obj any) (map[string]any, error) {
	if obj == nil {
		return nil, nil
	}
	res := handle(obj)
	if v, ok := res.(map[string]any); ok {
		return v, nil
	}
	return nil, fmt.Errorf("converting %T to map: not supported", obj)
}

// handle is a helper function that recursively handles
// the conversion of structs, maps, and slices to a map[string]any.
func handle(obj any) any {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		return handleMap(obj)
	case reflect.Struct:
		return handleStruct(obj)
	case reflect.Slice:
		return handleSlice(obj)
	default:
		return obj
	}
}

// handleStruct handles the conversion of a struct to a map[string]any.
// It uses the "json" and "yaml" tags to determine the key names.
func handleStruct(obj any) any {
	res := map[string]any{}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := range val.NumField() {
		field := val.Type().Field(i)

		// Cannot get value from unexported field
		if !field.IsExported() {
			continue
		}

		name := field.Name
		value := val.Field(i)
		tagName, tagOpts := getTag(field)
		if tagName != "" {
			name = tagName
		}

		// Omit struct tag "-"
		if _, ok := xslices.FindFunc(tagOpts, func(s string) bool {
			return s == "-"
		}); ok || (name == "-" && len(tagOpts) == 0) {
			continue
		}

		if _, ok := xslices.FindFunc(tagOpts, func(s string) bool {
			return s == "omitempty"
		}); ok {
			if reflect.DeepEqual(value.Interface(), reflect.Zero(val.Field(i).Type()).Interface()) {
				continue
			}
		}

		// Handle nested structs with inline tag
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}
		if value.Kind() == reflect.Struct || value.Kind() == reflect.Map {
			if _, ok := xslices.FindFunc(tagOpts, func(s string) bool {
				return s == "inline"
			}); ok {
				if nestedValues, ok := handle(value.Interface()).(map[string]any); ok {
					for k, v := range nestedValues {
						if _, ok := res[k]; !ok {
							res[k] = v
						}
					}
					continue
				}
			}
		}

		res[name] = handle(value.Interface())
	}

	return res
}

// handleMap handles the conversion of a map to a map[string]any,
// recursively converting nested maps, slices and structs.
func handleMap(obj any) any {
	m := map[string]any{}
	val := reflect.ValueOf(obj)
	for _, key := range val.MapKeys() {
		k := key.Interface()
		v := val.MapIndex(key).Interface()
		if k == nil {
			continue
		}
		if v == nil {
			continue
		}
		m[fmt.Sprintf("%v", k)] = handle(v)
	}
	return m
}

// handleSlice handles the conversion of a slice to a slice of any,
// recursively converting nested maps, slices and structs.
func handleSlice(obj any) any {
	s := []any{}
	val := reflect.ValueOf(obj)
	for i := range val.Len() {
		s = append(s, handle(val.Index(i).Interface()))
	}
	return s
}

// getTag retrieves the tag name and options from a struct field.
// It checks for the "json" and "yaml" tags in that order.
// If one tag is empty, it will return the other tag.
// If both tags are empty, it returns an empty string and an empty slice.
func getTag(field reflect.StructField) (string, []string) {
	for _, category := range tagCategories {
		if tag := field.Tag.Get(category); tag != "" {
			splitTag := strings.Split(tag, ",")
			return splitTag[0], splitTag[1:]
		}
	}
	return "", []string{}
}
