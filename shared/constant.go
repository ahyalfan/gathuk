// Package shared
package shared

// scanning config for struct
var (
	name       Tag = "config" // if use name tag
	nestedName Tag = "nested" // if use nested tag
)

func GetTagName() Tag {
	return name
}

func GetTagNestedName() Tag {
	return nestedName
}
