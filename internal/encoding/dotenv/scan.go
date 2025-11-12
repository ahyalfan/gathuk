// Package dotenv
package dotenv

import (
	"os"
	"reflect"
	"strings"

	utility "github.com/ahyalfan/gathuk/internal/utils"
	"github.com/ahyalfan/gathuk/shared"
)

// scanWithNestedPrefix initiates the recursive scanning process to populate
// a struct from the parsed key-value pairs.
//
// This method validates that the input is a pointer and starts the recursive
// scanning process.
//
// Parameters:
//   - v: Pointer to the configuration struct to populate
//
// Returns:
//   - error: An error if scanning fails
//
// Panics if v is not a pointer.
func (c *Codec[T]) scanWithNestedPrefix(v *T) error {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("value is not pointer")
	}

	vt := reflect.ValueOf(v).Elem()
	parent := reflect.TypeOf(v)
	c.scanNestedWithNestedPrefix(parent, vt, "")

	return nil
}

// scanNestedWithNestedPrefix recursively scans a struct and populates its fields
// from the parsed key-value pairs, handling nested structures with prefixes.
//
// This method processes each field of the struct:
//   - For nested structs: Recursively processes with the appropriate prefix
//   - For basic types: Maps configuration keys to field values
//   - Respects `config` and `nested` struct tags
//   - Handles environment variable fallback based on DecodeOption
//
// Parameters:
//   - parent: The parent type (used to prevent infinite recursion)
//   - v: The reflect.Value of the struct to populate
//   - nestedPrefix: The prefix to prepend to field names (e.g., "DB_" for nested database config)
func (c *Codec[T]) scanNestedWithNestedPrefix(
	parent reflect.Type, v reflect.Value, nestedPrefix string,
) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := v.Type().Field(i)

		if structField.Type.Kind() == reflect.Struct && structField.Type != parent {
			nestedName := structField.Tag.Get(string(shared.GetTagNestedName()))
			if nestedName == "-" {
				continue
			}
			if nestedName == "" {
				nestedName = structField.Tag.Get(string(shared.GetTagName()))
			}
			if nestedName == "" {
				nestedName = utility.PascalToUpperSnakeCase(structField.Name)
			}
			if nestedPrefix != "" {
				nestedName = nestedPrefix + "_" + nestedName
			}
			c.scanNestedWithNestedPrefix(parent, field, nestedName)
			continue
		}

		var name string
		name = structField.Tag.Get(string(shared.GetTagName()))
		if name == "-" {
			continue
		}
		if name == "" {
			name = utility.PascalToUpperSnakeCase(structField.Name)
		}

		if nestedPrefix != "" {
			sub := nestedPrefix + "_"
			name = sub + name
		}
		name = strings.ToUpper(name)

		var (
			val []byte
			ok  bool
		)

		if c.do.AutomaticEnv {
			if c.do.PreferFileOverEnv {

				val, ok = c.temp[name]
				if !ok {
					r := os.Getenv(name)
					if r != "" {
						val, ok = []byte(r), true
					}
				}
			} else {
				r := os.Getenv(name)
				if r == "" {
					val, ok = c.temp[name]
				} else {
					val, ok = []byte(r), true
				}

			}
		} else {
			val, ok = c.temp[name]
		}

		if !ok || !field.CanSet() {
			continue
		}

		setValue(field, string(val))
	}
}
