package xstructs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/neticdk/go-stdlib/xslices"
)

// tagCategories is a list of tag categories to check for.
var tagCategories = []string{"json", "yaml"}

// WithTags allows you to specify custom tag categories to check for.
// It can be used to override the default "json" and "yaml" tags.
// The tags are checked in the order they are provided.
func WithTags(tags ...string) ToMapOptions {
	return func(h *handler) {
		h.tags = tags
	}
}

// WithAllowNoTags allows you to specify whether to allow fields without tags.
// If used, fields without tags will be included in the output map.
func WithAllowNoTags() ToMapOptions {
	return func(h *handler) {
		h.allowNoTags = true
	}
}

// ToMapOptions is a function that modifies the handler.
type ToMapOptions func(*handler)

// handler is a struct that contains the options for the ToMap function.
// It contains a list of tags to check for and a flag to allow fields
// without tags.
type handler struct {
	tags        []string
	allowNoTags bool
}

// tagWrapper is a struct that contains the name and options of a tag.
// It is used to store the tag information for a field.
// The name is the key name to use in the output map.
// The options are the options specified in the tag.
type tagWrapper struct {
	Name    string
	Options []string
}

// newHandler creates a new handler with the default options.
// It initializes the tags to the default "json" and "yaml" tags.
// It also initializes the allowNoTags flag to false.
// It can be modified using the ToMapOptions functions.
// It returns a pointer to the handler.
func newHandler(opts ...ToMapOptions) *handler {
	h := &handler{
		tags:        tagCategories,
		allowNoTags: false,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// ToMap converts a struct or map to a map[string]any.
// It handles nested structs, maps, and slices.
// It uses the "json" and "yaml" tags to determine the key names.
// It respects the `omitempty` tag for fields.
// It respects the `inline` tag for nested structs.
// It respects the `-` tag to omit fields.
//
// If the input is nil, it returns nil.
// If the input is not a struct or map, it returns an error.
func ToMap(obj any, opts ...ToMapOptions) (map[string]any, error) {
	handler := newHandler(opts...)

	if obj == nil {
		return nil, nil
	}
	res := handler.handle(obj)
	if v, ok := res.(map[string]any); ok {
		return v, nil
	}
	return nil, fmt.Errorf("converting %T to map: not supported", obj)
}

// handle is a helper function that recursively handles
// the conversion of structs, maps, and slices to a map[string]any.
func (h *handler) handle(obj any) any {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		return h.handleMap(obj)
	case reflect.Struct:
		return h.handleStruct(obj)
	case reflect.Slice:
		return h.handleSlice(obj)
	default:
		return obj
	}
}

// handleStruct handles the conversion of a struct to a map[string]any.
// It uses the "json" and "yaml" tags to determine the key names.
func (h *handler) handleStruct(obj any) any {
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
		tagInfo, err := h.getTag(field)
		if err != nil && !h.allowNoTags {
			continue
		}

		if h.allowNoTags && tagInfo == nil {
			tagInfo = &tagWrapper{
				Name:    "",
				Options: []string{},
			}
		}

		if tagInfo.Name != "" {
			name = tagInfo.Name
		}

		// Omit struct tag "-"
		if _, ok := xslices.FindFunc(tagInfo.Options, func(s string) bool {
			return s == "-"
		}); ok || (name == "-" && len(tagInfo.Options) == 0) {
			continue
		}

		if _, ok := xslices.FindFunc(tagInfo.Options, func(s string) bool {
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
			if _, ok := xslices.FindFunc(tagInfo.Options, func(s string) bool {
				return s == "inline"
			}); ok {
				if nestedValues, ok := h.handle(value.Interface()).(map[string]any); ok {
					for k, v := range nestedValues {
						if _, ok := res[k]; !ok {
							res[k] = v
						}
					}
					continue
				}
			}
		}

		res[name] = h.handle(value.Interface())
	}

	return res
}

// handleMap handles the conversion of a map to a map[string]any,
// recursively converting nested maps, slices and structs.
func (h *handler) handleMap(obj any) any {
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
		m[fmt.Sprintf("%v", k)] = h.handle(v)
	}
	return m
}

// handleSlice handles the conversion of a slice to a slice of any,
// recursively converting nested maps, slices and structs.
func (h *handler) handleSlice(obj any) any {
	s := []any{}
	val := reflect.ValueOf(obj)
	for i := range val.Len() {
		s = append(s, h.handle(val.Index(i).Interface()))
	}
	return s
}

// getTag retrieves the tag name and options from a struct field.
// It checks for the "json" and "yaml" tags in that order.
// If one tag is empty, it will return the other tag.
// If both tags are empty, it returns an empty string and an empty slice.
func (h *handler) getTag(field reflect.StructField) (*tagWrapper, error) {
	for _, category := range h.tags {
		if tag := field.Tag.Get(category); tag != "" {
			splitTag := strings.Split(tag, ",")
			return &tagWrapper{
				Name:    splitTag[0],
				Options: splitTag[1:],
			}, nil
		}
	}
	return nil, fmt.Errorf("no tag of %s found for field %s", strings.Join(h.tags, ", "), field.Name)
}
