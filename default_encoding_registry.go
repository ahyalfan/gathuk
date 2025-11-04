// Package gathuk
package gathuk

import (
	"errors"
	"strings"
	"sync"

	"github.com/ahyalfan/gathuk/internal/encoding/dotenv"
	"github.com/ahyalfan/gathuk/option"
)

type DefaultCodecRegistry[T any] struct {
	codecs map[string]option.Codec[T]

	mu sync.Mutex
}

func NewDefaultCodecRegister[T any]() *DefaultCodecRegistry[T] {
	dcr := new(DefaultCodecRegistry[T])
	dcr.codecs = make(map[string]option.Codec[T])
	return dcr
}

func (dcr *DefaultCodecRegistry[T]) RegisterCodec(format string, codec option.Codec[T]) {
	dcr.mu.Lock()
	defer dcr.mu.Unlock()

	if dcr.codecs == nil {
		panic("dcr.codecs is nil: codec map must be initialized before use")
	}

	format = strings.ToLower(format)

	dcr.codecs[format] = codec
}

func (dcr *DefaultCodecRegistry[T]) Encoder(format string) (option.Encoder[T], error) {
	if v, ok := dcr.codec(format); ok {
		return v, nil
	}
	return nil, errors.New("encoder not found for this format")
}

func (dcr *DefaultCodecRegistry[T]) Decoder(format string) (option.Decoder[T], error) {
	if v, ok := dcr.codec(format); ok {
		return v, nil
	}
	return nil, errors.New("decoder not found for this format")
}

func (dcr *DefaultCodecRegistry[T]) codec(format string) (option.Codec[T], bool) {
	format = strings.ToLower(format)
	if v, ok := dcr.codecs[format]; ok {
		return v, true
	}

	switch format {
	case "env":
		return &dotenv.Codec[T]{}, true
	default:
		return nil, false
	}
}
