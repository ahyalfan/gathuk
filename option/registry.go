// Package option provides configuration options and interfaces for encoding
// and decoding configuration data in various formats.
package option

// Encoder is an interface for encoding configuration structs to byte slices.
//
// Implementations should convert a configuration struct of type T into
// a byte representation in a specific format (e.g., JSON, YAML, ENV).
//
// Type parameter T represents the configuration struct type.
//
// Example implementation:
//
//	type JSONEncoder[T any] struct {
//	    option.DefaultCodec[T]
//	}
//
//	func (e *JSONEncoder[T]) Encode(val T) ([]byte, error) {
//	    return json.Marshal(val)
//	}
type Encoder[T any] interface {
	// Encode converts a configuration struct to bytes in the encoder's format.
	//
	// Parameters:
	//  - val: The configuration struct to encode
	//
	// Returns:
	//  - []byte: The encoded configuration data
	//  - error: An error if encoding fails
	Encode(T) ([]byte, error)

	// EncodeOptionApplier provides methods for applying encode options
	EncodeOptionApplier

	// marker is an internal method used to prevent external implementations
	// of this interface. Only types that embed DefaultCodec can satisfy this.
	marker()
}

// Decoder is an interface for decoding byte slices to configuration structs.
//
// Implementations should parse a byte slice in a specific format (e.g., JSON,
// YAML, ENV) and populate a configuration struct of type T.
//
// Type parameter T represents the configuration struct type.
//
// Example implementation:
//
//	type JSONDecoder[T any] struct {
//	    option.DefaultCodec[T]
//	}
//
//	func (d *JSONDecoder[T]) Decode(buf []byte,val *T)  error {
//	    err := json.Unmarshal(buf, val)
//	    return  err
//	}
type Decoder[T any] interface {
	// Decode parses byte data and returns a populated configuration struct.
	//
	// Parameters:
	//  - buf: The byte slice containing configuration data to decode
	//
	// Returns:
	//  - T: The decoded configuration struct
	//  - error: An error if decoding fails
	Decode([]byte, *T) error

	// DecodeOptionApplier provides methods for applying decode options
	DecodeOptionApplier

	// marker is an internal method used to prevent external implementations
	// of this interface. Only types that embed DefaultCodec can satisfy this.
	marker()
}

// Codec is an interface that combines both Encoder and Decoder functionality.
//
// A Codec can both encode configuration structs to bytes and decode bytes
// back to configuration structs in a specific format.
//
// Type parameter T represents the configuration struct type.
//
// Example implementation:
//
//	type JSONCodec[T any] struct {
//	    option.DefaultCodec[T]
//	}
//
//	func (c *JSONCodec[T]) Encode(val T) ([]byte, error) {
//	    return json.Marshal(val)
//	}
//
//	func (c *JSONCodec[T]) Decode(buf []byte,val *T) error {
//	    err := json.Unmarshal(buf, val)
//	    return  err
//	}
type Codec[T any] interface {
	Encoder[T]
	Decoder[T]
}

// EncoderRegistry is an interface for registries that provide encoders for
// different file formats.
//
// Implementations maintain a mapping from format names (e.g., "json", "yaml")
// to their corresponding encoder implementations.
//
// Type parameter T represents the configuration struct type.
//
// Example:
//
//	encoder, err := registry.Encoder("json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	data, err := encoder.Encode(config)
type EncoderRegistry[T any] interface {
	// Encoder returns an encoder for the specified format.
	//
	// Format names are typically case-insensitive (e.g., "json", "JSON", "Json"
	// should all work).
	//
	// Parameters:
	//  - format: The name of the format (e.g., "json", "yaml", "env")
	//
	// Returns:
	//  - Encoder[T]: The encoder implementation for the format
	//  - error: An error if no encoder is found for the format
	Encoder(format string) (Encoder[T], error)
}

// DecoderRegistry is an interface for registries that provide decoders for
// different file formats.
//
// Implementations maintain a mapping from format names (e.g., "json", "yaml")
// to their corresponding decoder implementations.
//
// Type parameter T represents the configuration struct type.
//
// Example:
//
//	decoder, err := registry.Decoder("env")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	config, err := decoder.Decode(fileData)
type DecoderRegistry[T any] interface {
	// Decoder returns a decoder for the specified format.
	//
	// Format names are typically case-insensitive (e.g., "env", "ENV", "Env"
	// should all work).
	//
	// Parameters:
	//  - format: The name of the format (e.g., "json", "yaml", "env")
	//
	// Returns:
	//  - Decoder[T]: The decoder implementation for the format
	//  - error: An error if no decoder is found for the format
	Decoder(format string) (Decoder[T], error)
}

// CodecRegistry is an interface that combines both EncoderRegistry and
// DecoderRegistry functionality.
//
// A CodecRegistry can provide both encoders and decoders for various
// configuration file formats.
//
// Type parameter T represents the configuration struct type.
//
// Example:
//
//	type Config struct {
//	    Port int
//	    Host string
//	}
//
//	registry := gathuk.NewDefaultCodecRegister[Config]()
//	registry.RegisterCodec("json", &JSONCodec[Config]{})
//
//	// Can get both encoder and decoder
//	encoder, _ := registry.Encoder("json")
//	decoder, _ := registry.Decoder("json")
type CodecRegistry[T any] interface {
	DecoderRegistry[T]
	EncoderRegistry[T]
}
