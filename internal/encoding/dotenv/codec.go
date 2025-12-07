// Package dotenv provides encoding and decoding functionality for .env file format.
//
// This package implements a codec that can parse .env configuration files and
// convert them to Go structs, as well as encode Go structs back to .env format.
//
// .env file format:
//   - Each line contains a key-value pair: KEY=value
//   - Lines starting with # are comments and are ignored
//   - Empty lines are ignored
//   - Keys are case-sensitive
//   - Values are read as-is until end of line or comment
//
// Example .env file:
//
//	# Database configuration
//	DB_HOST=localhost
//	DB_PORT=5432
//	DB_USER=admin
//
//	# Server configuration
//	SERVER_PORT=8080
//	DEBUG=true
//
// Supported Go types:
//   - string: Plain text values
//   - int, int64: Integer values
//   - float64: Floating-point values
//   - bool: Boolean values (true/false)
//
// Nested structs are supported using the `nested` tag to define prefixes:
//
//	type Config struct {
//	    Database struct {
//	        Host string  // Maps to DB_HOST
//	        Port int     // Maps to DB_PORT
//	    } 'config':"db"`
//	}
package dotenv

import (
	"bytes"
	"os"
	"reflect"
	"strconv"
	"strings"

	utility "github.com/ahyalfan/gathuk/internal/utils"
	"github.com/ahyalfan/gathuk/option"
	"github.com/ahyalfan/gathuk/shared"
)

// Codec implements the option.Codec interface for .env file format.
//
// It provides both encoding (struct to .env) and decoding (.env to struct)
// functionality with support for nested structures, custom field mapping,
// and environment variable integration.
//
// Type parameter T represents the configuration struct type.
//
// Example usage:
//
//	codec := &Codec[Config]{}
//	codec.ApplyDecodeOption(&option.DecodeOption{
//	    AutomaticEnv: true,
//	})
//
//	config, err := codec.Decode(fileData)
//	if err != nil {
//	    log.Fatal(err)
//	}
type Codec[T any] struct {
	option.DefaultCodec[T]

	// do contains decode options that control how values are read
	do *option.DecodeOption
	// eo contains encode options that control how values are written
	eo *option.EncodeOption

	// temp is a temporary map used during encoding/decoding to store
	// key-value pairs before converting them to/from the struct
	temp map[string][]byte
}

// ApplyEncodeOption sets the encode options for this codec.
//
// These options control how the codec behaves when encoding structs to .env format,
// such as whether to include environment variables in the output.
//
// Parameters:
//   - eo: Pointer to EncodeOption containing the options to apply
//
// Example:
//
//	codec.ApplyEncodeOption(&option.EncodeOption{
//	    AutomaticEnv: true,
//	})
func (c *Codec[T]) ApplyEncodeOption(eo *option.EncodeOption) {
	c.eo = eo
}

// CheckEncodeOption checks whether encode options have been applied to this codec.
//
// Returns:
//   - bool: true if encode options have been set, false otherwise
func (c *Codec[T]) CheckEncodeOption() bool {
	return c.eo != nil
}

// Encode converts a configuration struct to .env file format.
//
// The encoding process:
//  1. Flattens nested structures using prefixes defined by `nested` tags
//  2. Converts field names to UPPER_SNAKE_CASE
//  3. Applies custom field names from `config` tags
//  4. Formats each key-value pair as KEY=value
//
// Parameters:
//   - val: The configuration struct to encode
//
// Returns:
//   - []byte: The encoded .env file content
//   - error: An error if encoding fails
//
// Example:
//
//	type Config struct {
//	    Port int
//	    Host string
//	}
//
//	codec := &Codec[Config]{}
//	data, err := codec.Encode(Config{Port: 8080, Host: "localhost"})
//	// data contains:
//	// PORT=8080
//	// HOST=localhost
func (c *Codec[T]) Encode(val T) ([]byte, error) {
	if c.temp == nil {
		c.temp = make(map[string][]byte)
	}

	c.flattenWithNestedPrefix(val)
	// var build strings.Builder
	// for k, v := range c.temp {
	// 	build.WriteString(k)
	// 	build.WriteRune('=')
	// 	build.Write(v)
	// 	build.WriteRune('\n')
	// }
	//
	// return []byte(build.String()), nil

	var build []byte
	for k, v := range c.temp {
		build = append(build, []byte(k)...)
		build = append(build, '=')
		build = append(build, v...)
		build = append(build, '\n')
	}
	return build, nil
}

// ApplyDecodeOption sets the decode options for this codec.
//
// These options control how the codec behaves when decoding .env files to structs,
// such as whether to read from environment variables and how to handle conflicts
// between file and environment values.
//
// Parameters:
//   - do: Pointer to DecodeOption containing the options to apply
//
// Example:
//
//	codec.ApplyDecodeOption(&option.DecodeOption{
//	    AutomaticEnv:      true,
//	    PreferFileOverEnv: true,
//	})
func (c *Codec[T]) ApplyDecodeOption(do *option.DecodeOption) {
	c.do = do
}

// CheckDecodeOption checks whether decode options have been applied to this codec.
//
// Returns:
//   - bool: true if decode options have been set, false otherwise
func (c *Codec[T]) CheckDecodeOption() bool {
	return c.do != nil
}

// Decode parses .env file content and populates a configuration struct.
//
// The decoding process:
//  1. Parses each line of the .env file
//  2. Extracts key-value pairs (KEY=value format)
//  3. Ignores comments (lines starting with #) and empty lines
//  4. Maps keys to struct fields using field names or `config` tags
//  5. Handles nested structures using `nested` tag prefixes
//  6. Optionally reads from environment variables based on DecodeOption
//  7. Converts string values to appropriate Go types
//
// Supported line formats:
//   - KEY=value          # Standard format
//   - KEY=value # comment # With inline comment
//   - # comment          # Comment line (ignored)
//   - # Empty line (ignored)
//
// Parameters:
//   - buf: Byte slice containing .env file content
//
// Returns:
//   - T: The populated configuration struct
//   - error: An error if decoding fails
//
// Example:
//
//	codec := &Codec[Config]{}
//	codec.ApplyDecodeOption(&option.DecodeOption{
//	    AutomaticEnv: true,
//	})
//
//	data := []byte("PORT=8080\nHOST=localhost")
//	config, err := codec.Decode(data)
func (c *Codec[T]) Decode(buf []byte, val *T) error {
	if c.temp == nil {
		c.temp = make(map[string][]byte)
	}

	lines := bytes.SplitSeq(buf, []byte{'\n'})

	for line := range lines {

		line = bytes.TrimSpace(line)
		escape := bytes.IndexByte(line, '#')
		if escape != -1 {
			line = line[:escape]
		}

		bs := bytes.Split(line, []byte(" "))

		if len(bs) < 1 {
			continue
		}
		bs = bytes.Split(bs[0], []byte("="))

		if len(bs) < 2 {
			continue
		}

		c.temp[string(bs[0])] = bs[1]

		if c.do.PersistToOSEnv {
			err := os.Setenv(string(bs[0]), string(bs[1]))
			if err != nil {
				return nil
			}
		}
	}

	if c.do.AutomaticEnv {
		if c.do.PreferFileOverEnv {
			for _, e := range os.Environ() {
				pair := strings.SplitN(e, "=", 2)
				if _, ok := c.temp[pair[0]]; !ok {
					c.temp[pair[0]] = []byte(pair[1])
				}
			}
		} else {
			for _, e := range os.Environ() {
				pair := strings.SplitN(e, "=", 2)
				c.temp[pair[0]] = []byte(pair[1])
			}
		}
	}

	err := c.scanWithNestedPrefix(val)

	return err
}

// flattenWithNestedPrefix initiates the flattening process for encoding.
//
// This method prepares a struct for encoding by flattening nested structures
// and converting field names to configuration keys.
//
// Parameters:
//   - v: The configuration struct to flatten
//
// Returns:
//   - error: An error if flattening fails
func (c *Codec[T]) flattenWithNestedPrefix(v T) error {
	vt := reflect.ValueOf(v)
	if vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	parent := reflect.TypeOf(v)
	c.flattenNestedWithNestedPrefix(parent, vt, "")

	return nil
}

// flattenNestedWithNestedPrefix recursively flattens a struct into key-value pairs
// for encoding to .env format.
//
// This method processes each field of the struct:
//   - For nested structs: Recursively processes with the appropriate prefix
//   - For basic types: Converts to string and stores in temp map
//   - Respects `config` and `nested` struct tags
//
// Parameters:
//   - parent: The parent type (used to prevent infinite recursion)
//   - v: The reflect.Value of the struct to flatten
//   - nestedPrefix: The prefix to prepend to field names
func (c *Codec[T]) flattenNestedWithNestedPrefix(
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

		c.temp[name] = parseToBytes(field)
	}
}

// parseToBytes converts a struct field value to its byte representation.
//
// This function is used during encoding to convert Go values to strings
// that can be written to .env files.
//
// Supported types:
//   - string: Direct conversion to bytes
//   - int, int8, int16, int32, int64: Formatted as base-10 integer
//   - uint, uint8, uint16, uint32, uint64: Formatted as base-10 unsigned integer
//   - float32, float64: Formatted as floating-point number
//   - bool: Formatted as "true" or "false"
//
// Parameters:
//   - field: The reflect.Value of the field to convert
//
// Returns:
//   - []byte: The byte representation of the field value, or nil for unsupported types
func parseToBytes(field reflect.Value) []byte {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}

	// Basic kinds
	switch field.Kind() {
	case reflect.String:
		return []byte(field.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(field.Int(), 10))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(field.Uint(), 10))

	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(field.Float(), 'f', -1, 64))

	case reflect.Bool:
		return []byte(strconv.FormatBool(field.Bool()))
	}
	return nil
}
