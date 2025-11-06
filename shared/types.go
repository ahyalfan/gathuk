// Package shared provides utility types and functions for handling custom tags used in structs.
package shared

// Tag is a custom type alias for string that is used to represent tags in struct field annotations.
// This allows for better type safety and clearer intent when working with struct tags, especially in
// scenarios where specific tags like "config" or "nested" are used to define struct field properties
type Tag string

func (t *Tag) Set(v Tag) {
	*t = v
}

func (t Tag) Get() Tag {
	return t
}
