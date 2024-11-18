package utils

import (
	"fmt"
	"reflect"
)

// CombineMaps combines multiple maps into a single map.
// If there are duplicate keys, the value from the last map will overwrite previous values.
func CombineMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Iterate over each map in the slice
	for _, m := range maps {
		// Iterate over each key-value pair in the current map
		for key, value := range m {
			result[key] = value
		}
	}

	return result
}

// StructToMap converts a struct to a map using its JSON tags as keys.
// Fields with json:"-" are skipped, and fields with json:"omitempty" are excluded if they are empty.
func StructToMap(s interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(s)

	// Ensure the input is a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a pointer to a struct")
	}

	// Dereference the pointer to access the struct's value
	val = val.Elem()

	// Iterate over the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Get the JSON tag from the struct field
		tag := field.Tag.Get("json")

		// Handle the "-" tag which means the field should be skipped
		if tag == "-" {
			continue
		}

		// If "omitempty" is set, and the field is empty, skip it
		if tag == "omitempty" && isEmptyValue(fieldValue) {
			continue
		}

		// If there's a JSON tag, use it as the key, otherwise use the field name
		if tag != "" {
			// Split tag to handle both key and omitempty
			tagParts := splitJsonTag(tag)
			fieldKey := tagParts[0]

			// Add to map
			result[fieldKey] = fieldValue.Interface()
		}
	}

	return result, nil
}

// isEmptyValue checks if a reflect.Value is considered empty (zero or nil)
func isEmptyValue(val reflect.Value) bool {
	// A value is empty if it is the zero value for its type
	return val.IsZero()
}

// splitJsonTag splits a JSON tag like `json:"key,omitempty"` into its parts.
func splitJsonTag(tag string) []string {
	// JSON tag can contain both the key and "omitempty", we split it by comma
	tagParts := []string{tag}
	if len(tag) > 0 && tag[0] == '"' {
		tagParts = append(tagParts, tag[1:])
	}
	return tagParts
}
