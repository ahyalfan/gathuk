// Package option
package option

type Encoder[T any] interface {
	Encode(T) ([]byte, error)
	EncodeOptionApplier

	marker()
}

type Decoder[T any] interface {
	Decode([]byte) (T, error)

	DecodeOptionApplier
	marker()
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
