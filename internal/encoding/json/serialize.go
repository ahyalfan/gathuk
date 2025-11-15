// Package json
package json

import (
	"bytes"
	"fmt"
	"strconv"
)

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

func (c *Codec[T]) serialize(node ASTNode) ([]byte, error) {
	var buf bytes.Buffer
	err := c.serializeNode(&buf, node, 0)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

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

func (c *Codec[T]) serializeString(buf *bytes.Buffer, str StringNode) error {
	escape := escepeStringByte([]byte(str.Value))
	buf.WriteByte('"')
	buf.Write(escape)
	buf.WriteByte('"')
	return nil
}

func (c *Codec[T]) serializeNumber(buf *bytes.Buffer, num NumberNode) error {
	buf.WriteString(strconv.FormatFloat(num.Value, 'g', -1, 64))
	return nil
}

func (c *Codec[T]) serializeBoolean(buf *bytes.Buffer, bool BooleanNode) error {
	if bool.Value {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	return nil
}

func (c *Codec[T]) serializeNull(buf *bytes.Buffer) error {
	buf.WriteString("null")
	return nil
}
