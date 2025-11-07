// Package gathuk
package gathuk

import "reflect"

// resolveFilenames accepts a list of filenames and returns them as-is.
// If no filenames are provided, it returns a slice containing ".env" as the fallback.
//
// This function ensures that there is always at least one configuration file to load,
// defaulting to the common ".env" file name if none is specified.
//
// Parameters:
//   - filenames: Variable number of file paths
//
// Returns a slice of filenames (or [".env"] if empty).
//
// Example:
//
//	files := resolveFilenames()                    // Returns: [".env"]
//	files := resolveFilenames("config.env")        // Returns: ["config.env"]
//	files := resolveFilenames("a.env", "b.env")    // Returns: ["a.env", "b.env"]
func resolveFilenames(filenames ...string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

// isZeroValue checks if a reflect.Value represents a zero value or nil.
//
// This function uses reflect.DeepEqual to compare the value with the zero value
// of its type. It's used during configuration merging to determine whether a
// value should override an existing value.
//
// Parameters:
//   - v: reflect.Value to check
//
// Returns true if the value is zero/nil, false otherwise.
//
// Example:
//
//	var s string
//	v := reflect.ValueOf(s)
//	isZero := isZeroValue(v)  // Returns: true
//
//	s = "hello"
//	v = reflect.ValueOf(s)
//	isZero := isZeroValue(v)  // Returns: false
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
