// Package json
package json

import (
	"github.com/ahyalfan/gathuk/option"
)

type Codec[T any] struct {
	option.DefaultCodec[T]

	do *option.DecodeOption
	eo *option.EncodeOption
}

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

func (c *Codec[T]) Encode(val T) ([]byte, error) {
	astN, err := c.StructToAST(&val)
	if err != nil {
		return nil, err
	}
	r, err := c.serialize(astN)

	return r, err
}

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
