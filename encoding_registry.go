// Package gathuk
package gathuk

import (
	"errors"
	"strings"
	"sync"

	"github.com/ahyalfan/gathuk/internal/encoding/dotenv"
)

type Encoder[T any] interface {
	Encode(T) ([]byte, error)
}

type Decoder[T any] interface {
	Decode([]byte) (T, error)
}

type Codec[T any] interface {
	Encoder[T]
	Decoder[T]
}

type EncoderRegistry[T any] interface {
	Encoder(format string) (Encoder[T], error)
}

type DecoderRegistry[T any] interface {
	Decoder(format string) (Decoder[T], error)
}

type CodecRegistry[T any] interface {
	DecoderRegistry[T]
	EncoderRegistry[T]
}

type DefaultCodecRegistry[T any] struct {
	codecs map[string]Codec[T]

	sync.Mutex
}

func NewDefaultCodecRegister[T any]() *DefaultCodecRegistry[T] {
	dcr := new(DefaultCodecRegistry[T])
	dcr.codecs = make(map[string]Codec[T])
	return dcr
}

func (dcr *DefaultCodecRegistry[T]) RegisterCodec(format string, codec Codec[T]) {
	dcr.Lock()
	defer dcr.Unlock()

	if dcr.codecs == nil {
		panic("dcr.codecs is nil: codec map must be initialized before use")
	}

	format = strings.ToLower(format)

	dcr.codecs[format] = codec
}

func (dcr *DefaultCodecRegistry[T]) Encoder(format string) (Encoder[T], error) {
	if v, ok := dcr.codec(format); ok {
		return v, nil
	}
	return nil, errors.New("encoder not found for this format")
}

func (dcr *DefaultCodecRegistry[T]) Decoder(format string) (Decoder[T], error) {
	if v, ok := dcr.codec(format); ok {
		return v, nil
	}
	return nil, errors.New("decoder not found for this format")
}

func (dcr *DefaultCodecRegistry[T]) codec(format string) (Codec[T], bool) {
	if v, ok := dcr.codecs[format]; ok {
		return v, true
	}
	format = strings.ToLower(format)
	switch format {
	case "env":
		return &dotenv.Codec[T]{}, true
	default:
		return nil, false
	}
}
