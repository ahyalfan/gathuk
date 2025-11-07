// Package option
package option

// DecodeOption contains options that control how configuration data is decoded
// from files and environment variables.
//
// These options allow fine-grained control over configuration precedence,
// environment variable binding, and persistence behavior.
//
// Example:
//
//	opt := &DecodeOption{
//	    AutomaticEnv:      true,  // Read from OS environment
//	    PersistToOSEnv:    false, // Don't write back to OS env
//	    PreferFileOverEnv: true,  // File values override env vars
//	}
type DecodeOption struct {
	AutomaticEnv      bool // jika true, baca OS environment otomatis
	PersistToOSEnv    bool // jika true, hasil decode disimpan di OS env juga
	PreferFileOverEnv bool // jika true, config file diutamakan dibanding OS env / string
}

// EncodeOption contains options that control how configuration data is encoded
// to files and other output destinations.
//
// These options determine the source of data when encoding configuration
// and how to handle environment variables during encoding.
//
// Example:
//
//	opt := &EncodeOption{
//	    AutomaticEnv:      true,  // Include OS environment values
//	    PreferFileOverEnv: false, // Env vars override struct values
//	}
type EncodeOption struct {
	AutomaticEnv      bool // jika true, baca OS environment otomatis
	PreferFileOverEnv bool // jika true, config file diutamakan dibanding OS env / string
}

// DecodeOptionApplier is an interface for types that can accept and apply
// decode options.
//
// Implementations should store the provided options and use them during
// the decode process.
type DecodeOptionApplier interface {
	// ApplyDecodeOption applies the given decode options to the decoder.
	//
	// Parameters:
	//  - opt: Pointer to DecodeOption containing the options to apply
	ApplyDecodeOption(*DecodeOption)
	// CheckDecodeOption checks if decode options have been applied.
	//
	// Returns:
	//  - bool: true if options are set, false otherwise
	CheckDecodeOption() bool
}

// EncodeOptionApplier is an interface for types that can accept and apply
// encode options.
//
// Implementations should store the provided options and use them during
// the encode process.
type EncodeOptionApplier interface {
	// ApplyEncodeOption applies the given encode options to the encoder.
	//
	// Parameters:
	//  - opt: Pointer to EncodeOption containing the options to apply
	ApplyEncodeOption(*EncodeOption)
	// CheckEncodeOption checks if encode options have been applied.
	//
	// Returns:
	//  - bool: true if options are set, false otherwise
	CheckEncodeOption() bool
}
