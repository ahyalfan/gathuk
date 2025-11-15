// Package shared provides utility types and functions for handling custom tags used in structs.
package shared

// Package-level variables defining the standard struct tags used by Gathuk.
//
// These tags control how struct fields are mapped to configuration keys
// and how nested structures are handled.
var (
	// name is the struct tag used to map a field to a specific configuration key.
	//
	// Usage: `config:"custom_key_name"`
	//
	// When this tag is present, the field will be mapped to the specified
	// configuration key instead of the auto-generated key based on the field name.
	//
	// Special value:
	//  - Use `config:"-"` to exclude a field from configuration parsing
	//
	// Example:
	//  type Config struct {
	//      Port int `config:"server_port"`  // Maps to SERVER_PORT in env, server_port in JSON
	//  }
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

// SetTagName sets a custom tag name for field mapping.
//
// This function allows changing the default "config" tag to a different name.
// This is useful when integrating with other libraries or when you want to use
// a different naming convention.
//
// Warning: Changing the tag name affects all codecs globally. Make sure to call
// this function before creating any Gathuk instances or codecs.
//
// Parameters:
//   - tagName: The new tag name to use for field mapping
//
// Example:
//
//	// Use "json" tag instead of "config" tag
//	shared.SetTagName("json")
//
//	type Config struct {
//	    Port int `json:"port"`  // Now uses json tag instead of config tag
//	}
//
//	// Later, in codec:
//	field := structType.Field(i)
//	customName := field.Tag.Get(string(shared.GetTagName()))  // Gets value from "json" tag
func SetTagName(tagName string) {
	name = Tag(tagName)
}

// SetTagNestedName sets a custom tag name for nested structure prefixes.
//
// This function allows changing the default "nested" tag to a different name.
// This is useful when you want to use a different naming convention or avoid
// conflicts with other libraries.
//
// Warning: Changing the nested tag name affects all codecs globally. Make sure
// to call this function before creating any Gathuk instances or codecs.
//
// Parameters:
//   - tagNestedName: The new tag name to use for nested structure prefixes
//
// Example:
//
//	// Use "prefix" tag instead of "nested" tag
//	shared.SetTagNestedName("prefix")
//
//	type Config struct {
//	    Database struct {
//	        Host string  // Will use PREFIX from "prefix" tag
//	    } `prefix:"db"`  // Now uses prefix tag instead of nested tag
//	}
//
//	// Later, in codec:
//	field := structType.Field(i)
//	prefix := field.Tag.Get(string(shared.GetTagNestedName()))  // Gets value from "prefix" tag
func SetTagNestedName(tagNestedName string) {
	nestedName = Tag(tagNestedName)
}
