// Package dotenv
package dotenv

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
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
	err := c.scanNestedWithNestedPrefix(parent, vt, "")

	return err
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
) error {
	if !v.CanSet() {
		return newError(nestedPrefix, "value not settable")
	}

	switch v.Kind() {
	case reflect.Interface:
		native, err := c.toNative(nestedPrefix)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(native))
	case reflect.Struct:
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
				err := c.scanNestedWithNestedPrefix(parent, field, nestedName)
				if err != nil {
					return err
				}
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

			val, ok := c.temp[name]

			if !ok || !field.CanSet() {
				continue
			}

			err := setValue(field, string(val))
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		err := c.toMap(v, nestedPrefix)
		if err != nil {
			return err
		}
	default:
		newError(nestedPrefix, "unsupported type: %v+", v.Type())
	}

	return nil
}

// toMap converts the parsed key-value pairs into a map[string]V where V is the map value type.
//
// This method is used when the target type is a map instead of a struct.
// It filters keys by prefix if provided and converts string values to the appropriate type.
//
// For .env files, all keys are stored in UPPER_SNAKE_CASE format. If a prefix is provided,
// only keys that start with that prefix will be included in the resulting map.
//
// The map key in the result includes the prefix, so for prefix "DB_" and key "DB_HOST",
// the map will contain "DB_HOST" as the key (not just "HOST").
//
// Type conversion is handled automatically:
//   - If the target map value type is string, values are used as-is
//   - If the target type is int/float/bool, string values are parsed accordingly
//   - If parsing fails, an error is returned
//
// Parameters:
//   - v: The reflect.Value of the map to populate
//   - prefix: Optional prefix to filter keys (e.g., "DB_" for database configs)
//
// Returns:
//   - error: An error if type conversion or map creation fails
func (c *Codec[T]) toMap(v reflect.Value, prefix string) error {
	if v.Type().Key().Kind() != reflect.String {
		return newError(prefix, "map key must be string, got %s", v.Type().Key())
	}

	mapType := v.Type()
	newMap := reflect.MakeMap(mapType)

	for k, v := range c.temp {
		if prefix != "" {
			if !strings.HasPrefix(k, prefix) {
				continue
			}
		}
		nested := prefix + "_" + k
		if prefix == "" {
			nested = k
		}

		elemValue := reflect.New(mapType.Elem()).Elem()

		err := setValue(elemValue, string(v))
		if err != nil {
			return err
		}

		newMap.SetMapIndex(reflect.ValueOf(nested), elemValue)
	}
	v.Set(newMap)
	return nil
}

// toNative converts the parsed key-value pairs into native Go types for interface{}/any.
//
// This method is used when the target type is interface{} or any. It creates a
// map[string]any where each value is converted to the most appropriate native type:
//   - Boolean strings ("true"/"false") → bool
//   - Integer strings ("123") → int64
//   - Float strings ("3.14") → float64
//   - Everything else → string
//
// Unlike toMap, this method always returns a map[string]any regardless of prefix,
// but it filters keys by prefix if provided.
//
// Parameters:
//   - prefix: Optional prefix to filter keys (empty string means include all)
//
// Returns:
//   - any: A map[string]any containing the filtered and converted values
//   - error: An error if conversion fails
func (c *Codec[T]) toNative(prefix string) (any, error) {
	m := make(map[string]any)
	for k, v := range c.temp {
		if prefix != "" {
			if !strings.HasPrefix(k, prefix) {
				continue
			}
		}
		var converted any
		err := setValue(reflect.ValueOf(&converted).Elem(), string(v))
		if err != nil {
			return nil, newError(prefix, "%v", err)
		}
		m[k] = converted
	}
	return m, nil
}

// setValue sets a struct field value from a string using reflection.
//
// This function handles type conversion from string to the appropriate Go type.
// It supports pointer types by automatically dereferencing them.
//
// Supported types:
//   - string: Direct assignment
//   - int, int64: Parsed as base-10 integer
//   - float64: Parsed as floating-point number
//   - bool: Parsed as boolean (true/false)
//   - any: Parsed any value
//
// Parameters:
//   - field: The reflect.Value of the field to set
//   - val: The string value to convert and assign
//
// return error if type conversion fails.
func setValue(field reflect.Value, val string) error {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	// Basic kinds
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int64:
		i64, err := strconv.ParseInt(val, 0, 64)
		if err != nil {
			return newError("", "convert string to int error: %+v", err)
		}
		field.SetInt(i64)
	case reflect.Float64:
		f64, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return newError("", "convert string to float error: %+v", err)
		}
		field.SetFloat(f64)
	case reflect.Bool:
		bVal, err := strconv.ParseBool(val)
		if err != nil {
			return newError("", "convert string to bool error: %+v", err)
		}
		field.SetBool(bVal)
	case reflect.Interface:
		var converted any
		if b, err := strconv.ParseBool(val); err == nil {
			converted = b
		} else if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			converted = i
		} else if f, err := strconv.ParseFloat(val, 64); err == nil {
			converted = f
		} else {
			converted = val
		}

		field.Set(reflect.ValueOf(converted))
	}
	return nil
}

func setValueAny(field reflect.Value, val any) {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	// Basic kinds
	switch field.Kind() {
	case reflect.String:
		s, ok := val.(string)
		if !ok {
			return
		}
		field.SetString(s)

	case reflect.Int, reflect.Int64:
		i64, ok := val.(int)
		if !ok {
			log.Fatalf("convert string to int error: %+v", ok)
		}
		field.SetInt(int64(i64))
	case reflect.Float64:
		f64, ok := val.(float64)
		if !ok {
			log.Fatalf("convert string to float error: %+v", ok)
		}
		field.SetFloat(f64)
	case reflect.Bool:
		bVal, ok := val.(bool)
		if !ok {
			log.Fatalf("convert string to bool error: %+v", ok)
		}
		field.SetBool(bVal)
	}
}

func newError(path, format string, args ...any) error {
	if path != "" {
		return fmt.Errorf("ast unmarshal error at %s: %w", path, fmt.Errorf(format, args...))
	}
	return fmt.Errorf("ast unmarshal error: "+format, args...)
}
