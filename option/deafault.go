// Package option
package option

type DefaultCodec[T any] struct{}

func (DefaultCodec[T]) marker() {}

func (dc *DefaultCodec[T]) ApplyDecodeOption(*DecodeOption) {}
func (dc *DefaultCodec[T]) CheckDecodeOption() bool         { return false }

func (dc *DefaultCodec[T]) Decode([]byte) (T, error) {
	var zeroValue T
	return zeroValue, nil
}

func (dc *DefaultCodec[T]) ApplyEncodeOption(*EncodeOption) {}
func (dc *DefaultCodec[T]) CheckEncodeOption() bool         { return false }

func (dc *DefaultCodec[T]) Encode([]byte) (T, error) {
	var zeroValue T
	return zeroValue, nil
}
