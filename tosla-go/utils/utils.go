package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

// LuhnCheck validates a credit card number using the LUHN algorithm
func LuhnCheck(cardNumber string) bool {
	// Remove all spaces (just in case the input contains spaces)
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")

	// Validate if the card number only contains digits
	if len(cardNumber) < 13 || len(cardNumber) > 19 { // Credit card numbers are usually between 13 and 19 digits
		return false
	}

	// Reverse the digits to start from the rightmost
	var sum int
	for i := len(cardNumber) - 1; i >= 0; i-- {
		// Convert the character to an integer
		digit, err := strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return false // Return false if the character is not a valid digit
		}

		// Double every second digit, starting from the rightmost digit
		if (len(cardNumber)-i)%2 == 0 {
			digit *= 2
			// If doubling the digit results in a number greater than 9, subtract 9
			if digit > 9 {
				digit -= 9
			}
		}

		// Add the digit to the sum
		sum += digit
	}

	// If the total sum is divisible by 10, it's a valid card number
	return sum%10 == 0
}
