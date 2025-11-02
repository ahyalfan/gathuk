// Package gathuk
package gathuk

import "reflect"

// resolveFilenames accepts a list of filenames and returns them.
// If no filenames are provided, it returns a slice containing ".env" as the fallback.
func resolveFilenames(filenames ...string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

// isZeroValue a check zero or nil use reflect packgae
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
