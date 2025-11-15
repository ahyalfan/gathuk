// Package utility
package utility

import (
	"strings"
	"unicode"
)

// PascalToUpperSnakeCase converts a string from PascalCase to UPPER_SNAKE_CASE.
// It transforms the input string by:
// 1. Inserting an underscore (_) between each capital letter and its preceding character (if necessary).
// 2. Converting all letters to uppercase.
//
// Example:
//
//	Input: "PascalCaseExample"
//	Output: "PASCAL_CASE_EXAMPLE"
//
// If the input string is already in UPPER_SNAKE_CASE, it will return the string unchanged.
func PascalToUpperSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 { // Add underscore before uppercase letters, except for the first character
				result.WriteRune('_')
			}
			result.WriteRune(r)
		} else {
			result.WriteRune(unicode.ToUpper(r))
		}
	}
	return result.String()
}

// PascalToLowerSnakeCase converts a string from PascalCase to lower_snake_case.
func PascalToLowerSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 { // Add underscore before uppercase letters, except for the first character
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
