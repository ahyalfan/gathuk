// Package json provides encoding and decoding functionality for JSON format.
package json

import (
	"github.com/ahyalfan/gathuk/option"
)

// Codec implements the option.Codec interface for JSON format.
//
// It provides both encoding (struct to JSON) and decoding (JSON to struct)
// functionality with support for nested structures, arrays, maps, and
// custom field mapping via struct tags.
//
// Type parameter T represents the configuration struct type.
//
// The codec uses a multi-stage pipeline:
//   - Decode: JSON bytes → Tokenize → Parse → AST → Map → Struct
//   - Encode: Struct → Map → AST → Serialize → JSON bytes
//
// Example usage:
//
//	type Config struct {
//	    Port int    `config:"port"`
//	    Host string `config:"host"`
//	}
//
//	codec := &Codec[Config]{}
//
//	// Decode
//	var config Config
//	err := codec.Decode(jsonData, &config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Encode
//	data, err := codec.Encode(config)
type Codec[T any] struct {
	option.DefaultCodec[T]

	do *option.DecodeOption
	eo *option.EncodeOption
}

// ApplyEncodeOption sets the encode options for this codec.
//
// These options control how the codec behaves when encoding structs to JSON,
// such as formatting preferences and field selection.
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

// Encode converts a configuration struct to JSON bytes.
//
// The encoding process follows these steps:
//  1. Convert struct to AST using mapper
//  2. Serialize AST to JSON bytes
//
// Parameters:
//   - val: The configuration struct to encode
//
// Returns:
//   - []byte: The encoded JSON data
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
//	// data contains: {"port": 8080, "host": "localhost"}
func (c *Codec[T]) Encode(val T) ([]byte, error) {
	astN, err := c.StructToAST(&val)
	if err != nil {
		return nil, err
	}
	r, err := c.serialize(astN)

	return r, err
}

// ApplyDecodeOption sets the decode options for this codec.
//
// These options control how the codec behaves when decoding JSON to structs,
// such as handling of unknown fields and type conversion.
//
// Parameters:
//   - do: Pointer to DecodeOption containing the options to apply
//
// Example:
//
//	codec.ApplyDecodeOption(&option.DecodeOption{
//	    AutomaticEnv: true,
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

// Decode parses JSON bytes and populates a configuration struct.
//
// The decoding process follows these steps:
//  1. Tokenize JSON bytes into tokens
//  2. Parse tokens into AST
//  3. Map AST to struct using reflection
//
// Parameters:
//   - val: JSON bytes to decode
//   - dst: Pointer to the destination struct
//
// Returns:
//   - error: An error if decoding fails at any stage
//
// Example:
//
//	type Config struct {
//	    Port int    `config:"port"`
//	    Host string `config:"host"`
//	}
//
//	codec := &Codec[Config]{}
//	var config Config
//
//	jsonData := []byte(`{"port": 8080, "host": "localhost"}`)
//	err := codec.Decode(jsonData, &config)
//	if err != nil {
//	    log.Fatal(err)
//	}
func (c *Codec[T]) Decode(val []byte, dst *T) error {
	if val == nil {
		return nil
	}
	tokens, err := Tokenize(val)
	if err != nil {
		return err
	}
	ast, err := Parser(tokens)
	if err != nil {
		return err
	}
	err = c.ASTToStruct(ast, dst)
	if err != nil {
		return err
	}
	return nil
}
