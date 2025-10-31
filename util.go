// Package gathuk
package gathuk

// resolveFilenames accepts a list of filenames and returns them.
// If no filenames are provided, it returns a slice containing ".env" as the fallback.
func resolveFilenames(filenames ...string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}
