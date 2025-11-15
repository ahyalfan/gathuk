// Package json provides encoding and decoding functionality for JSON format.
package json

import (
	"bytes"
	"fmt"
	"strconv"
)

// escapeStringByte escapes special characters in a byte slice for JSON.
//
// Characters that need escaping:
//   - " (quote) → \"
//   - \ (backslash) → \\
//   - \n (newline) → \\n
//   - \t (tab) → \\t
//   - \r (carriage return) → \\r
//
// Parameters:
//   - s: The byte slice to escape
//
// Returns:
//   - []byte: The escaped byte slice
func escepeStringByte(s []byte) []byte {
	result := []byte{}
	for _, ch := range s {
		switch ch {
		case '"':
			result = append(result, []byte("\\\"")...)
		case '\\':
			result = append(result, []byte("\\\\")...)
		case '\n':
			result = append(result, []byte("\\n")...)
		case '\t':
			result = append(result, []byte("\\t")...)
		case '\r':
			result = append(result, []byte("\\r")...)
		default:
			result = append(result, ch)
		}
	}
	return result
}

// serialize converts an AST node to JSON bytes.
//
// This is the final step in the encoding pipeline, converting the
// Abstract Syntax Tree representation into a valid JSON byte sequence.
//
// The serializer handles:
//   - Proper JSON formatting (with optional pretty-printing)
//   - String escaping (quotes, backslashes, newlines, etc.)
//   - Number formatting
//   - Boolean and null literals
//
// Parameters:
//   - node: The root AST node to serialize
//
// Returns:
//   - []byte: The JSON representation as bytes
//   - error: An error if serialization fails
//
// Example:
//
//	ast := ObjectNode{
//	    Value: map[string]ASTNode{
//	        "name": StringNode{"John"},
//	        "age": NumberNode{30},
//	    },
//	}
//	data, err := codec.serialize(ast)
//	// data: []byte(`{"name": "John", "age": 30}`)
func (c *Codec[T]) serialize(node ASTNode) ([]byte, error) {
	var buf bytes.Buffer
	err := c.serializeNode(&buf, node, 0)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// serializeNode recursively serializes an AST node to a buffer.
//
// Parameters:
//   - buf: The buffer to write JSON bytes to
//   - node: The AST node to serialize
//   - depth: Current nesting depth (for pretty printing)
//
// Returns:
//   - error: An error if serialization fails
func (c *Codec[T]) serializeNode(buf *bytes.Buffer, node ASTNode, depth int) error {
	switch n := node.(type) {
	case ObjectNode:
		c.serializeObject(buf, n, depth)
	case ArrayNode:
		c.serializeArray(buf, n, depth)
	case StringNode:
		c.serializeString(buf, n)
	case NumberNode:
		c.serializeNumber(buf, n)
	case BooleanNode:
		c.serializeBoolean(buf, n)
	case NullNode:
		c.serializeNull(buf)
	default:
		return fmt.Errorf("unknown node type: %T", node)
	}

	return nil
}

// serializeObject serializes an ObjectNode to JSON format.
//
// Output format: {"key": value, "key": value}
//
// Parameters:
//   - buf: The buffer to write to
//   - obj: The ObjectNode to serialize
//   - depth: Current nesting depth
//
// Returns:
//   - error: An error if serialization fails
func (c *Codec[T]) serializeObject(buf *bytes.Buffer, obj ObjectNode, depth int) error {
	buf.WriteByte('{')

	first := true
	for key, value := range obj.Value {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		escape := escepeStringByte([]byte(key))

		buf.WriteByte('"')
		buf.Write(escape)
		buf.WriteByte('"')
		buf.WriteByte(':')
		buf.WriteByte(' ')

		if err := c.serializeNode(buf, value, depth+1); err != nil {
			return err
		}
	}

	buf.WriteByte('}')
	return nil
}

// serializeArray serializes an ArrayNode to JSON format.
//
// Output format: [value, value, value]
//
// Parameters:
//   - buf: The buffer to write to
//   - arr: The ArrayNode to serialize
//   - depth: Current nesting depth
//
// Returns:
//   - error: An error if serialization fails
func (c *Codec[T]) serializeArray(buf *bytes.Buffer, arr ArrayNode, depth int) error {
	buf.WriteByte('[')

	for i, elem := range arr.Value {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := c.serializeNode(buf, elem, depth+1); err != nil {
			return err
		}
	}

	buf.WriteByte(']')
	return nil
}

// serializeString serializes a StringNode to JSON format.
//
// Handles proper escaping of special characters:
//   - " → \"
//   - \ → \\
//   - \n → \\n
//   - \t → \\t
//   - \r → \\r
//
// Parameters:
//   - buf: The buffer to write to
//   - str: The StringNode to serialize
//
// Returns:
//   - error: An error if serialization fails
func (c *Codec[T]) serializeString(buf *bytes.Buffer, str StringNode) error {
	escape := escepeStringByte([]byte(str.Value))
	buf.WriteByte('"')
	buf.Write(escape)
	buf.WriteByte('"')
	return nil
}

// serializeNumber serializes a NumberNode to JSON format.
//
// Parameters:
//   - buf: The buffer to write to
//   - num: The NumberNode to serialize
//
// Returns:
//   - error: An error if serialization fails
func (c *Codec[T]) serializeNumber(buf *bytes.Buffer, num NumberNode) error {
	buf.WriteString(strconv.FormatFloat(num.Value, 'g', -1, 64))
	return nil
}

// serializeBoolean serializes a BooleanNode to JSON format.
//
// Output: "true" or "false"
//
// Parameters:
//   - buf: The buffer to write to
//   - bool: The BooleanNode to serialize
//
// Returns:
//   - error: An error if serialization fails
func (c *Codec[T]) serializeBoolean(buf *bytes.Buffer, bool BooleanNode) error {
	if bool.Value {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	return nil
}

// serializeNull serializes a NullNode to JSON format.
//
// Output: "null"
//
// Parameters:
//   - buf: The buffer to write to
//
// Returns:
//   - error: An error if serialization fails
func (c *Codec[T]) serializeNull(buf *bytes.Buffer) error {
	buf.WriteString("null")
	return nil
}
