// Package option provides configuration options and interfaces for encoding
// and decoding configuration data in various formats.
package option

// DefaultCodec is a base implementation that provides default (no-op) implementations
// of the Codec interface methods.
//
// This type is intended to be embedded in custom codec implementations to satisfy
// the Codec interface. By embedding DefaultCodec, custom implementations only need
// to override the methods they actually use, while inheriting safe defaults for
// the rest.
//
// Type parameter T represents the configuration struct type.
//
// Example usage:
//
//	type JSONCodec[T any] struct {
//	    option.DefaultCodec[T]
//	    decodeOpt *option.DecodeOption
//	}
//
//	func (c *JSONCodec[T]) Decode(buf []byte) (T, error) {
//	    var val T
//	    err := json.Unmarshal(buf, &val)
//	    return val, err
//	}
//
//	func (c *JSONCodec[T]) Encode(val T) ([]byte, error) {
//	    return json.Marshal(val)
//	}
//
// By embedding DefaultCodec, JSONCodec automatically gets:
//   - marker() implementation (required by interface)
//   - Default option applier methods
//   - Default option checker methods
type DefaultCodec[T any] struct{}

// marker is an internal method that prevents external implementations of the
// Codec interface.
//
// This method exists to ensure that all codec implementations embed DefaultCodec,
// which provides a stable foundation for future interface changes without breaking
// external implementations.
//
// This method should never be called directly.
func (DefaultCodec[T]) marker() {}

// ApplyDecodeOption is a no-op implementation of the DecodeOptionApplier interface.
//
// This default implementation does nothing. Custom codec implementations that
// support decode options should override this method to store the options for
// later use during decoding.
//
// Parameters:
//   - opt: Pointer to DecodeOption (ignored in default implementation)
//
// Example override:
//
//	func (c *CustomCodec[T]) ApplyDecodeOption(opt *DecodeOption) {
//	    c.decodeOpt = opt
//	}
func (dc *DefaultCodec[T]) ApplyDecodeOption(*DecodeOption) {}

// CheckDecodeOption is a no-op implementation that always returns false.
//
// This default implementation indicates that no decode options have been applied.
// Custom codec implementations should override this method to return true when
// options have been set.
//
// Returns:
//   - bool: Always returns false in default implementation
//
// Example override:
//
//	func (c *CustomCodec[T]) CheckDecodeOption() bool {
//	    return c.decodeOpt != nil
//	}
func (dc *DefaultCodec[T]) CheckDecodeOption() bool { return false }

// Decode is a no-op implementation that returns a zero value.
//
// This default implementation does nothing and returns the zero value for type T.
// Custom codec implementations MUST override this method to provide actual
// decoding functionality.
//
// Parameters:
//   - buf: Byte slice to decode (ignored in default implementation)
//
// Returns:
//   - T: Zero value of type T
//   - error: Always returns nil
//
// Example override:
//
//	func (c *JSONCodec[T]) Decode(buf []byte) (T, error) {
//	    var val T
//	    err := json.Unmarshal(buf, &val)
//	    return val, err
//	}
func (dc *DefaultCodec[T]) Decode([]byte) (T, error) {
	var zeroValue T
	return zeroValue, nil
}

func (dc *DefaultCodec[T]) DecodePointer([]byte, *T) error {
	return nil
}

// ApplyEncodeOption is a no-op implementation of the EncodeOptionApplier interface.
//
// This default implementation does nothing. Custom codec implementations that
// support encode options should override this method to store the options for
// later use during encoding.
//
// Parameters:
//   - opt: Pointer to EncodeOption (ignored in default implementation)
//
// Example override:
//
//	func (c *CustomCodec[T]) ApplyEncodeOption(opt *EncodeOption) {
//	    c.encodeOpt = opt
//	}
func (dc *DefaultCodec[T]) ApplyEncodeOption(*EncodeOption) {}

// CheckEncodeOption is a no-op implementation that always returns false.
//
// This default implementation indicates that no encode options have been applied.
// Custom codec implementations should override this method to return true when
// options have been set.
//
// Returns:
//   - bool: Always returns false in default implementation
//
// Example override:
//
//	func (c *CustomCodec[T]) CheckEncodeOption() bool {
//	    return c.encodeOpt != nil
//	}
func (dc *DefaultCodec[T]) CheckEncodeOption() bool { return false }

// Encode is a no-op implementation that returns a zero value.
//
// This default implementation does nothing and returns the zero value for type T.
// Custom codec implementations MUST override this method to provide actual
// encoding functionality.
//
// Note: The signature appears to be incorrect (should accept T and return []byte),
// but this matches the current implementation. This is likely a bug that should
// be fixed in a future version.
//
// Parameters:
//   - buf: Byte slice (parameter appears incorrect, should be T)
//
// Returns:
//   - T: Zero value of type T (return type appears incorrect, should be []byte)
//   - error: Always returns nil
//
// Example correct override:
//
//	func (c *JSONCodec[T]) Encode(val T) ([]byte, error) {
//	    return json.Marshal(val)
//	}
func (dc *DefaultCodec[T]) Encode([]byte) (T, error) {
	var zeroValue T
	return zeroValue, nil
}
