// Package shared provides utility types and functions for handling custom tags used in structs.
package shared

// scanning config for struct
var (
	name       Tag = "config" // if use name tag
	nestedName Tag = "nested" // if use nested tag
)

// GetTagName returns the tag used for the name field (config).
func GetTagName() Tag {
	return name
}

// GetTagNestedName returns the tag used for nested fields.
func GetTagNestedName() Tag {
	return nestedName
}

func SetTagName(tagName string) {
	name = Tag(tagName)
}

func SetTagNestedName(tagNestedName string) {
	nestedName = Tag(tagNestedName)
}
