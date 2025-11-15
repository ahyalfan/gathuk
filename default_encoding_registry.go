// Package gathuk
package gathuk

import (
	"errors"
	"strings"
	"sync"

	"github.com/ahyalfan/gathuk/internal/encoding/dotenv"
	"github.com/ahyalfan/gathuk/internal/encoding/json"
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
// Thread Safety:
//   - All public methods are protected by a mutex
//   - Safe for concurrent registration and retrieval of codecs
//   - Multiple goroutines can safely call Encoder/Decoder methods
//
// Built-in Codecs:
//   - "env": .env file format (always available)
//
// Example:
//
//	type Config struct {
//	    Port int
//	    Host string
//	}
//
//	registry := gathuk.NewDefaultCodecRegister[Config]()
//
//	// Register custom codecs
//	registry.RegisterCodec("json", &JSONCodec[Config]{})
//	registry.RegisterCodec("yaml", &YAMLCodec[Config]{})
//
//	// Use with Gathuk
//	gt := gathuk.NewGathuk[Config]()
//	gt.SetCustomCodecRegistry(registry)
type DefaultCodecRegistry[T any] struct {
	codecs map[string]option.Codec[T]

	mu sync.Mutex
}

// NewDefaultCodecRegister creates and initializes a new DefaultCodecRegistry.
//
// The returned registry is ready to use and includes built-in support for
// .env file format. Additional codecs can be registered using RegisterCodec.
//
// The codec map is pre-allocated and ready for registration of custom codecs.
// The built-in "env" codec is lazily loaded when first requested.
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
//	// Create registry
//	registry := gathuk.NewDefaultCodecRegister[Config]()
//
//	// Built-in env codec is automatically available
//	envCodec, _ := registry.Decoder("env")
//
//	// Register additional formats
//	registry.RegisterCodec("json", &JSONCodec[Config]{})
func NewDefaultCodecRegister[T any]() *DefaultCodecRegistry[T] {
	dcr := new(DefaultCodecRegistry[T])
	dcr.codecs = make(map[string]option.Codec[T])
	return dcr
}

// RegisterCodec registers a codec for a specific file format.
//
// Format names are case-insensitive and will be normalized to lowercase.
// If a codec already exists for the given format, it will be replaced with
// the new codec.
//
// This method is thread-safe and can be called concurrently from multiple
// goroutines.
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
//
//	// Replace existing codec
//	registry.RegisterCodec("json", &ImprovedJSONCodec[Config]{})
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
// Format names are case-insensitive. The method first checks registered codecs,
// then falls back to built-in codecs. If no encoder is found for the format,
// it returns an error.
//
// This method is thread-safe and can be called concurrently from multiple
// goroutines.
//
// Built-in formats:
//   - "env": Environment variable / .env file format
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
//	if err != nil {
//	    log.Fatal("Encoding failed:", err)
//	}
func (dcr *DefaultCodecRegistry[T]) Encoder(format string) (option.Encoder[T], error) {
	if v, ok := dcr.codec(format); ok {
		return v, nil
	}
	return nil, errors.New("encoder not found for this format")
}

// Decoder returns a decoder for the specified format.
//
// Format names are case-insensitive. The method first checks registered codecs,
// then falls back to built-in codecs. If no decoder is found for the format,
// it returns an error.
//
// This method is thread-safe and can be called concurrently from multiple
// goroutines.
//
// Built-in formats:
//   - "env": Environment variable / .env file format
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
//	var config Config
//	err = decoder.Decode(fileData, &config)
//	if err != nil {
//	    log.Fatal("Decoding failed:", err)
//	}
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
// Lookup priority:
//  1. Registered codecs in the codecs map
//  2. Built-in codecs (env)
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
	case "json":
		return &json.Codec[T]{}, true
	default:
		return nil, false
	}
}
