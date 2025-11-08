// Package gathuk
package gathuk

import (
	"errors"
	"strings"
	"sync"

	"github.com/ahyalfan/gathuk/internal/encoding/dotenv"
	"github.com/ahyalfan/gathuk/option"
)

// DefaultCodecRegistry is a thread-safe registry that manages codecs for
// different configuration file formats. It provides a default implementation
// of the CodecRegistry interface.
//
// The registry includes built-in support for .env files and allows registration
// of additional codecs for other formats like JSON, YAML, TOML, etc.
//
// Type parameter T represents the configuration struct type that codecs will
// encode to and decode from.
//
// Example:
//
//	registry := gathuk.NewDefaultCodecRegister[Config]()
//	registry.RegisterCodec("json", &JSONCodec[Config]{})
//	registry.RegisterCodec("yaml", &YAMLCodec[Config]{})
type DefaultCodecRegistry[T any] struct {
	codecs map[string]option.Codec[T] // map of format names to their codec implementations
	mu     sync.Mutex                 // mutex to protect concurrent access to codecs map
}

// NewDefaultCodecRegister creates and initializes a new DefaultCodecRegistry.
//
// The returned registry is ready to use and includes built-in support for
// .env file format. Additional codecs can be registered using RegisterCodec.
//
// Type parameter T should be a struct type representing your application's configuration.
//
// Returns a pointer to DefaultCodecRegistry with an initialized codec map.
//
// Example:
//
//	type Config struct {
//	    Port int
//	    Host string
//	}
//
//	registry := gathuk.NewDefaultCodecRegister[Config]()
func NewDefaultCodecRegister[T any]() *DefaultCodecRegistry[T] {
	dcr := new(DefaultCodecRegistry[T])
	dcr.codecs = make(map[string]option.Codec[T])
	return dcr
}

// RegisterCodec registers a codec for a specific file format.
//
// Format names are case-insensitive and will be normalized to lowercase.
// If a codec already exists for the given format, it will be replaced.
//
// This method is thread-safe and can be called concurrently.
//
// Parameters:
//   - format: The format name (e.g., "json", "yaml", "toml"). Case-insensitive
//   - codec: The codec implementation for this format
//
// Panics if the codec map is nil (should never happen if created with NewDefaultCodecRegister).
//
// Example:
//
//	registry := gathuk.NewDefaultCodecRegister[Config]()
//
//	// Register JSON codec
//	registry.RegisterCodec("json", &JSONCodec[Config]{})
//
//	// Register YAML codec
//	registry.RegisterCodec("yaml", &YAMLCodec[Config]{})
//
//	// Register TOML codec (case-insensitive)
//	registry.RegisterCodec("TOML", &TOMLCodec[Config]{})
func (dcr *DefaultCodecRegistry[T]) RegisterCodec(format string, codec option.Codec[T]) {
	dcr.mu.Lock()
	defer dcr.mu.Unlock()

	if dcr.codecs == nil {
		panic("dcr.codecs is nil: codec map must be initialized before use")
	}

	format = strings.ToLower(format)

	dcr.codecs[format] = codec
}

// Encoder returns an encoder for the specified format.
//
// Format names are case-insensitive. If no encoder is registered for the format,
// this method returns an error.
//
// This method is thread-safe and can be called concurrently.
//
// Parameters:
//   - format: The format name (e.g., "json", "yaml", "env")
//
// Returns:
//   - option.Encoder[T]: The encoder implementation for the format
//   - error: An error if no encoder is found for the format
//
// Example:
//
//	encoder, err := registry.Encoder("json")
//	if err != nil {
//	    log.Fatal("JSON encoder not found:", err)
//	}
//
//	data, err := encoder.Encode(config)
func (dcr *DefaultCodecRegistry[T]) Encoder(format string) (option.Encoder[T], error) {
	if v, ok := dcr.codec(format); ok {
		return v, nil
	}
	return nil, errors.New("encoder not found for this format")
}

// Decoder returns a decoder for the specified format.
//
// Format names are case-insensitive. If no decoder is registered for the format,
// this method returns an error.
//
// This method is thread-safe and can be called concurrently.
//
// Parameters:
//   - format: The format name (e.g., "json", "yaml", "env")
//
// Returns:
//   - option.Decoder[T]: The decoder implementation for the format
//   - error: An error if no decoder is found for the format
//
// Example:
//
//	decoder, err := registry.Decoder("env")
//	if err != nil {
//	    log.Fatal("ENV decoder not found:", err)
//	}
//
//	config, err := decoder.Decode(fileData)
func (dcr *DefaultCodecRegistry[T]) Decoder(format string) (option.Decoder[T], error) {
	if v, ok := dcr.codec(format); ok {
		return v, nil
	}
	return nil, errors.New("decoder not found for this format")
}

// codec is an internal method that retrieves a codec for the specified format.
//
// This method first checks the registered codecs map. If no codec is found,
// it checks for built-in codecs (currently only "env" format is built-in).
//
// Format names are case-insensitive.
//
// Parameters:
//   - format: The format name to look up
//
// Returns:
//   - option.Codec[T]: The codec implementation if found
//   - bool: true if a codec was found, false otherwise
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
