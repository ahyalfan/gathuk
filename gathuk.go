// Package gathuk provides a type-safe, flexible configuration management library
// for Go applications. It supports loading configurations from multiple file formats
// (currently .env), automatic environment variable binding, nested structures, and
// type-safe configuration merging.
//
// Basic usage:
//
//	type Config struct {
//	    Port int
//	    Host string
//	}
//
//	gt := gathuk.NewGathuk[Config]()
//	err := gt.LoadConfigFiles("config.env")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	config := gt.GetConfig()
//
// The library uses Go generics to provide compile-time type safety and supports
// custom struct tags for field mapping and nested structure prefixes.
package gathuk

import (
	"bytes"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ahyalfan/gathuk/option"
)

// Gathuk is the main configuration manager that handles loading, parsing,
// and merging configuration from multiple sources.
//
// Type parameter T represents the configuration struct type that will be
// populated with values from configuration files and environment variables.
//
// Example:
//
//	type AppConfig struct {
//	    Database struct {
//	        Host string
//	        Port int
//	    } `nested:"db"`
//	}
//
//	gt := gathuk.NewGathuk[AppConfig]()
type Gathuk[T any] struct {
	// globalDecodeOpt contains decode options applied to all decoders
	// unless overridden by format-specific options
	globalDecodeOpt option.DecodeOption
	// globalEncodeOpt contains encode options applied to all encoders
	// unless overridden by format-specific options
	globalEncodeOpt option.EncodeOption

	Mode string // dev, staging, production. mungkin set modenya di taruh di flag pas jalanin binary
	// mode file example dev.env,stag.env,dev.json

	// ConfigFiles is a list of base configuration file paths that will be
	// loaded when LoadConfigFiles is called without arguments
	ConfigFiles []string

	// value stores the parsed and merged configuration struct
	value T

	// map value -> if map feature ready, like convert to map or write use map

	// CodecRegistry manages encoders and decoders for different file formats.
	// By default, it includes support for .env files
	CodecRegistry option.CodecRegistry[T]

	// logger is used for internal logging. Defaults to text handler writing to stdout
	logger *slog.Logger
}

// Option is an interface for applying configuration options to Gathuk instance.
// This interface follows the functional options pattern for extensible configuration.
type Option[T any] interface {
	apply(g *Gathuk[T])
}

// optionFunc is a function type that implements the Option interface,
// allowing functions to be used as options.
type optionFunc[T any] func(g *Gathuk[T])

// apply implements the Option interface for optionFunc.
func (fn optionFunc[T]) apply(g *Gathuk[T]) {
	fn(g)
}

// NewGathuk creates and initializes a new Gathuk instance with default settings.
//
// The returned instance includes:
//   - A default codec registry with support for .env files
//   - A default logger writing to stdout
//   - Empty configuration ready to be populated
//
// Type parameter T should be a struct type representing your application's configuration.
//
// Example:
//
//	type Config struct {
//	    Port int
//	    Host string
//	}
//
//	gt := gathuk.NewGathuk[Config]()
func NewGathuk[T any]() *Gathuk[T] {
	g := &Gathuk[T]{}
	g.CodecRegistry = NewDefaultCodecRegister[T]()
	g.logger = slog.New(slog.NewTextHandler(os.Stdout, nil)) // default slog
	return g
}

// SetCustomCodecRegistry replaces the default codec registry with a custom one.
// This allows you to add support for additional file formats beyond .env.
//
// The provided codec registry must not be nil, otherwise this method will panic.
//
// Returns the Gathuk instance for method chaining.
//
// Example:
//
//	registry := gathuk.NewDefaultCodecRegister[Config]()
//	registry.RegisterCodec("json", &JSONCodec[Config]{})
//
//	gt := gathuk.NewGathuk[Config]()
//	gt.SetCustomCodecRegistry(registry)
func (g *Gathuk[T]) SetCustomCodecRegistry(c option.CodecRegistry[T]) *Gathuk[T] {
	if c == nil {
		panic("codec registry not nil")
	}
	g.CodecRegistry = c
	return g
}

// SetDecodeOption sets the decode options for a specific file format.
// These options control how configuration values are read and merged.
//
// Parameters:
//   - format: The file format (e.g., "env", "json", "yaml")
//   - decodeOption: Pointer to DecodeOption containing the configuration
//
// Panics if no decoder exists for the specified format.
//
// Example:
//
//	opt := &option.DecodeOption{
//	    AutomaticEnv: true,
//	    PreferFileOverEnv: true,
//	}
//	gt.SetDecodeOption("env", opt)
func (g *Gathuk[T]) SetDecodeOption(format string, decodeOption *option.DecodeOption) {
	c, err := g.CodecRegistry.Decoder(format)
	if err != nil {
		g.logger.Error(err.Error())
		panic("set decode option failed")
	}
	c.ApplyDecodeOption(decodeOption)
}

// SetEncodeOption sets the encode options for a specific file format.
// These options control how configuration values are written to files.
//
// Parameters:
//   - format: The file format (e.g., "env", "json", "yaml")
//   - encodeOption: Pointer to EncodeOption containing the configuration
//
// Panics if no encoder exists for the specified format.
//
// Example:
//
//	opt := &option.EncodeOption{
//	    AutomaticEnv: true,
//	}
//	gt.SetEncodeOption("env", opt)
func (g *Gathuk[T]) SetEncodeOption(format string, encodeOption *option.EncodeOption) {
	c, err := g.CodecRegistry.Encoder(format)
	if err != nil {
		g.logger.Error(err.Error())
		panic("set decode option failed")
	}
	c.ApplyEncodeOption(encodeOption)
}

// SetConfigFiles sets the base configuration files that will be loaded
// when LoadConfigFiles is called without arguments.
//
// This method does not load the files immediately; it only stores the file paths.
// Call LoadConfigFiles to actually load and parse the configuration.
//
// Parameters:
//   - srcFiles: Variable number of file paths to set as base configuration files
//
// Example:
//
//	gt.SetConfigFiles("base.env", "default.env")
//	// Later, load these plus additional files
//	err := gt.LoadConfigFiles("override.env")
func (g *Gathuk[T]) SetConfigFiles(srcFiles ...string) {
	g.ConfigFiles = srcFiles
}

// LoadConfigFiles loads configuration from one or more files and merges them
// into the configuration struct. Files are processed in order, with later files
// overriding values from earlier ones.
//
// If no files are specified and no base files are set via SetConfigFiles,
// this method will attempt to load from ".env" by default.
//
// The merge behavior:
//   - Later files override earlier files
//   - Zero values are not merged (existing non-zero values are preserved)
//   - Nested structs are merged recursively
//
// Parameters:
//   - srcFiles: Variable number of configuration file paths to load
//
// Returns an error if any file cannot be read or parsed.
//
// Example:
//
//	// Load single file
//	err := gt.LoadConfigFiles("config.env")
//
//	// Load and merge multiple files
//	err := gt.LoadConfigFiles("base.env", "dev.env", "local.env")
//
//	// Load base files plus additional file
//	gt.SetConfigFiles("base.env")
//	err := gt.LoadConfigFiles("override.env")
func (g *Gathuk[T]) LoadConfigFiles(srcFiles ...string) error {
	srcFiles = resolveFilenames(append(g.ConfigFiles, srcFiles...)...)
	for _, filename := range srcFiles {
		err := g.loadFile(filename, &g.value)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadConfig loads configuration from an io.Reader with the specified format
// and merges it into the configuration struct.
//
// This method is useful when you want to load configuration from sources
// other than files, such as HTTP responses, embedded resources, or in-memory buffers.
//
// Parameters:
//   - src: io.Reader containing the configuration data
//   - format: The format of the configuration (e.g., "env", "json", "yaml")
//
// Returns an error if reading or parsing fails.
//
// Example:
//
//	file, err := os.Open("config.env")
//	if err != nil {
//	    return err
//	}
//	defer file.Close()
//
//	err = gt.LoadConfig(file, "env")
//
//	// Or from a string
//	config := strings.NewReader("PORT=8080\nHOST=localhost")
//	err = gt.LoadConfig(config, "env")
func (g *Gathuk[T]) LoadConfig(src io.Reader, format string) error {
	err := g.load(src, format, &g.value)
	if err != nil {
		return err
	}
	return nil
}

// loadFile is an internal method that opens and loads a single configuration file.
// It automatically determines the file format from the file extension.
//
// Parameters:
//   - filename: Path to the configuration file
//
// Returns the parsed configuration struct and any error encountered.
func (g *Gathuk[T]) loadFile(filename string, val *T) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	ext := strings.Trim(filepath.Ext(filename), ".")

	return g.load(f, ext, val)
}

// load is an internal method that reads and parses configuration data from an io.Reader.
//
// Parameters:
//   - src: io.Reader containing the configuration data
//   - format: The format of the configuration data
//
// Returns the parsed configuration struct and any error encountered.
func (g *Gathuk[T]) load(src io.Reader, format string, val *T) error {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, src)
	if err != nil {
		return err
	}

	by := buf.Bytes()

	dc, err := g.CodecRegistry.Decoder(format)
	if err != nil {
		return err
	}

	if ok := dc.CheckDecodeOption(); !ok {
		dc.ApplyDecodeOption(&g.globalDecodeOpt)
	}

	err = dc.Decode(by, val)
	if err != nil {
		return err
	}

	return nil
}

// WriteConfigFile writes the configuration struct to a file with the specified permissions.
//
// The file format is automatically determined from the file extension.
// If the file already exists, it will be truncated.
//
// Parameters:
//   - dst: Destination file path
//   - mode: File permissions (e.g., 0644). Use 0 for default permissions
//   - config: Configuration struct to write
//
// Returns an error if the file cannot be created or written.
//
// Example:
//
//	config := Config{
//	    Port: 8080,
//	    Host: "localhost",
//	}
//
//	err := gt.WriteConfigFile("output.env", 0644, config)
func (g *Gathuk[T]) WriteConfigFile(dst string, mode fs.FileMode, config T) error {
	err := g.writeFile(dst, mode, config)
	if err != nil {
		return err
	}
	return nil
}

// WriteConfig writes the configuration struct to an io.Writer in the specified format.
//
// This method is useful when you want to write configuration to destinations
// other than files, such as HTTP responses, network connections, or in-memory buffers.
//
// Parameters:
//   - out: io.Writer to write the configuration to
//   - format: The output format (e.g., "env", "json", "yaml")
//   - config: Configuration struct to write
//
// Returns an error if encoding or writing fails.
//
// Example:
//
//	var buf bytes.Buffer
//	err := gt.WriteConfig(&buf, "env", config)
//	fmt.Println(buf.String())
//
//	// Or write to stdout
//	err := gt.WriteConfig(os.Stdout, "env", config)
func (g *Gathuk[T]) WriteConfig(out io.Writer, format string, config T) error {
	err := g.write(out, format, config)
	if err != nil {
		return err
	}
	return nil
}

// writeFile is an internal method that creates a file and writes configuration to it.
//
// Parameters:
//   - dst: Destination file path
//   - mode: File permissions (use 0 for default)
//   - config: Configuration struct to write
//
// Returns an error if file creation or writing fails.
func (g *Gathuk[T]) writeFile(dst string, mode fs.FileMode, config T) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer f.Close()

	if mode != 0 {
		err = f.Chmod(mode)
		if err != nil {
			return err
		}
	}

	ext := strings.Trim(filepath.Ext(dst), ".")

	err = g.write(f, ext, config)
	if err != nil {
		return err
	}

	return nil
}

// write is an internal method that encodes and writes configuration to an io.Writer.
//
// Parameters:
//   - out: io.Writer to write to
//   - format: Output format
//   - config: Configuration struct to write
//
// Returns an error if encoding or writing fails.
func (g *Gathuk[T]) write(out io.Writer, format string, config T) error {
	enc, err := g.CodecRegistry.Encoder(format)
	if err != nil {
		return err
	}

	bys, err := enc.Encode(config)
	if err != nil {
		return err
	}
	_, err = out.Write(bys)
	if err != nil {
		return err
	}

	return nil
}

// GetConfig returns the current configuration struct.
//
// This method returns a copy of the configuration, so modifications to the
// returned struct will not affect the internal configuration state.
//
// Returns the configuration struct of type T.
//
// Example:
//
//	config := gt.GetConfig()
//	fmt.Printf("Port: %d, Host: %s\n", config.Port, config.Host)
func (g *Gathuk[T]) GetConfig() T {
	return g.value
}

// mergeStruct recursively merges configuration from src into dst.
//
// The merge behavior:
//   - Only non-zero values from src are copied to dst
//   - Nested structs are merged recursively
//   - Unexported fields are skipped
//
// Parameters:
//   - dst: Destination struct pointer to merge into
//   - src: Source struct pointer to merge from
//
// Returns an error if merging fails.
func (g *Gathuk[T]) mergeStruct(dst, src any) error {
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src).Elem()

	for i := 0; i < dv.NumField(); i++ {
		df := dv.Field(i)
		sf := sv.Field(i)

		if !df.CanSet() {
			continue
		}

		switch df.Kind() {
		case reflect.Struct:
			if err := g.mergeStruct(df.Addr().Interface(), sf.Addr().Interface()); err != nil {
				return err
			}
		default:
			if !isZeroValue(sf) {
				df.Set(sf)
			}
		}
	}

	return nil
}
