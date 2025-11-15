// Package utility
package utility

import (
	"strings"
	"unicode"
)

// PascalToUpperSnakeCase converts a string from PascalCase to UPPER_SNAKE_CASE.
//
// This function is used to automatically generate configuration key names from
// struct field names. It handles the conversion by:
//  1. Inserting an underscore (_) before each uppercase letter (except the first)
//  2. Converting all letters to uppercase
//
// The conversion is useful for mapping Go struct field names (which follow PascalCase
// convention) to environment variable names (which follow UPPER_SNAKE_CASE convention).
//
// Note: This function treats each uppercase letter as the start of a new word,
// which means acronyms like "APIKey" will become "A_P_I_KEY" rather than "API_KEY".
// Use the `config` struct tag for custom naming when this behavior is not desired.
//
// Parameters:
//   - s: The input string in PascalCase format
//
// Returns:
//   - string: The converted string in UPPER_SNAKE_CASE format
//
// Examples:
//
//	PascalToUpperSnakeCase("PascalCase")      // Returns: "PASCAL_CASE"
//	PascalToUpperSnakeCase("DatabaseURL")     // Returns: "DATABASE_U_R_L"
//	PascalToUpperSnakeCase("APIKey")          // Returns: "A_P_I_KEY"
//	PascalToUpperSnakeCase("UserID")          // Returns: "USER_I_D"
//	PascalToUpperSnakeCase("Port")            // Returns: "PORT"
//	PascalToUpperSnakeCase("MaxConnections")  // Returns: "MAX_CONNECTIONS"
//	PascalToUpperSnakeCase("isEnabled")       // Returns: "IS_ENABLED"
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
//
// This function is used for JSON field name generation where lowercase with
// underscores is the common convention. It handles the conversion by:
//  1. Inserting an underscore (_) before each uppercase letter (except the first)
//  2. Converting all uppercase letters to lowercase
//
// This is particularly useful for JSON configuration where field names typically
// use lower_snake_case convention rather than PascalCase or camelCase.
//
// Note: Similar to PascalToUpperSnakeCase, this function treats each uppercase
// letter as the start of a new word. For acronyms like "APIKey", use the
// `config` struct tag to specify custom field names.
//
// Parameters:
//   - s: The input string in PascalCase format
//
// Returns:
//   - string: The converted string in lower_snake_case format
//
// Examples:
//
//	PascalToLowerSnakeCase("PascalCase")      // Returns: "pascal_case"
//	PascalToLowerSnakeCase("DatabaseURL")     // Returns: "database_u_r_l"
//	PascalToLowerSnakeCase("APIKey")          // Returns: "a_p_i_key"
//	PascalToLowerSnakeCase("UserID")          // Returns: "user_i_d"
//	PascalToLowerSnakeCase("Port")            // Returns: "port"
//	PascalToLowerSnakeCase("MaxConnections")  // Returns: "max_connections"
func PascalToLowerSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
