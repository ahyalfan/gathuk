# Gathuk

![Go Version](https://img.shields.io/badge/Go-1.24.3-blue.svg)
![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)

**Gathuk** is a type-safe, flexible configuration management library for Go that converts configuration files into strongly-typed structs. It supports multiple file formats (currently `.env`), nested structures, and automatic environment variable binding.

## Features

- 🎯 **Type-Safe**: Uses Go generics for compile-time type safety
- 📁 **Multiple File Support**: Load and merge configurations from multiple files
- 🔄 **Automatic Environment Variables**: Automatically bind OS environment variables to struct fields
- 🏗️ **Nested Structures**: Full support for nested struct configurations with custom prefixes
- 🔧 **Flexible Options**: Configure priority between file configs and environment variables
- 💾 **Write Support**: Export configurations back to files
- 🎨 **Custom Codecs**: Extensible codec system for adding new file formats
- 🚀 **Zero Dependencies**: Minimal external dependencies

## Installation

```bash
go get github.com/ahyalfan/gathuk
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/ahyalfan/gathuk"
)

type Config struct {
    SimpleC string
    SimpleE int
}

func main() {
    // Create a new Gathuk instance
    gt := gathuk.NewGathuk[Config]()

    // Load configuration from file
    err := gt.LoadConfigFiles(".env")
    if err != nil {
        panic(err)
    }

    // Get the parsed configuration
    config := gt.GetConfig()
    fmt.Printf("Config: %+v\n", config)
}
```

### Nested Structures

```go
type Database struct {
    User       string
    Server     string `config:"server_port"`
    PoolingMax int    `config:"poling_max_pool"`
}

type AppConfig struct {
    Debug    bool     `config:"debug_c"`
    Database Database `nested:"db"`
}

func main() {
    gt := gathuk.NewGathuk[AppConfig]()
    err := gt.LoadConfigFiles("config.env")
    if err != nil {
        panic(err)
    }

    config := gt.GetConfig()
    fmt.Printf("Database User: %s\n", config.Database.User)
}
```

**config.env:**

```env
DEBUG_C=true
DB_USER=dbtest
DB_SERVER_PORT=5432
DB_POLING_MAX_POOL=100
```

## Struct Tags

Gathuk supports two main struct tags:

### `config` Tag

Maps struct fields to specific configuration keys:

```go
type Config struct {
    Port int `config:"server_port"`  // Maps to SERVER_PORT in .env
}
```

### `nested` Tag

Defines prefix for nested structures:

```go
type Config struct {
    Database Database `nested:"db"`  // All Database fields will have DB_ prefix
}

type Database struct {
    Host string  // Maps to DB_HOST
    Port int     // Maps to DB_PORT
}
```

Use `-` to ignore fields:

```go
type Config struct {
    Internal string `config:"-"`  // Will be ignored
}
```

## Configuration Options

### Decode Options

```go
gt := gathuk.NewGathuk[Config]()

// Enable automatic environment variable binding
gt.globalDecodeOpt.AutomaticEnv = true

// Prefer file values over environment variables
gt.globalDecodeOpt.PreferFileOverEnv = true

// Persist decoded values to OS environment
gt.globalDecodeOpt.PersistToOSEnv = true

err := gt.LoadConfigFiles("config.env")
```

### Option Behavior

| Option              | Description                                                                               |
| ------------------- | ----------------------------------------------------------------------------------------- |
| `AutomaticEnv`      | When `true`, automatically reads from OS environment variables                            |
| `PreferFileOverEnv` | When `true`, prioritizes file config over environment variables (requires `AutomaticEnv`) |
| `PersistToOSEnv`    | When `true`, saves decoded values to OS environment variables                             |

## Multiple File Loading

Load and merge configurations from multiple files:

```go
gt := gathuk.NewGathuk[Config]()

// Method 1: Load multiple files at once
err := gt.LoadConfigFiles("base.env", "override.env")

// Method 2: Set base files, then load additional files
gt.SetConfigFiles("base.env")
err := gt.LoadConfigFiles("override.env")
```

**Merge Behavior**: Later files override values from earlier files. Zero values are not merged.

## Writing Configuration

Export your configuration to files:

```go
config := Config{
    SimpleC: "value",
    SimpleE: 100,
}

// Write to file
err := gt.WriteConfigFile("output.env", 0644, config)

// Or write to io.Writer
var buf bytes.Buffer
err := gt.WriteConfig(&buf, "env", config)
```

## Advanced Usage

### Custom Codec Registry

Create and register custom codecs for different file formats:

```go
// Create custom codec
type JSONCodec[T any] struct {
    option.DefaultCodec[T]
}

func (c *JSONCodec[T]) Decode(buf []byte) (T, error) {
    var value T
    err := json.Unmarshal(buf, &value)
    return value, err
}

func (c *JSONCodec[T]) Encode(val T) ([]byte, error) {
    return json.Marshal(val)
}

// Register codec
registry := gathuk.NewDefaultCodecRegister[Config]()
registry.RegisterCodec("json", &JSONCodec[Config]{})

gt := gathuk.NewGathuk[Config]()
gt.SetCustomCodecRegistry(registry)
```

### Loading from io.Reader

```go
file, err := os.Open("config.env")
if err != nil {
    panic(err)
}
defer file.Close()

err = gt.LoadConfig(file, "env")
```

### Format-Specific Options

```go
gt := gathuk.NewGathuk[Config]()

// Set decode options for specific format
decodeOpt := &option.DecodeOption{
    AutomaticEnv: true,
}
gt.SetDecodeOption("env", decodeOpt)

// Set encode options for specific format
encodeOpt := &option.EncodeOption{
    AutomaticEnv: true,
}
gt.SetEncodeOption("env", encodeOpt)
```

## Field Name Mapping

By default, Gathuk converts PascalCase field names to UPPER_SNAKE_CASE:

| Struct Field  | Env Variable     |
| ------------- | ---------------- |
| `SimpleValue` | `SIMPLE_VALUE`   |
| `DatabaseURL` | `DATABASE_U_R_L` |
| `APIKey`      | `A_P_I_KEY`      |

Use the `config` tag for custom mapping:

```go
type Config struct {
    APIKey string `config:"api_key"`  // Maps to API_KEY
}
```

## Examples

### Example 1: Simple Configuration

```go
// config.env
// SIMPLE_C=hore
// SIMPLE_E=2

type Simple struct {
    SimpleC string
    SimpleE int
}

gt := gathuk.NewGathuk[Simple]()
err := gt.LoadConfigFiles("config.env")
// Result: {SimpleC: "hore", SimpleE: 2}
```

### Example 2: With Environment Variable Override

```go
// config.env
// USER=file_user

type Config struct {
    User   string
    Editor string
}

// Set environment variable
os.Setenv("USER", "env_user")
os.Setenv("EDITOR", "nvim")

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true
gt.LoadConfigFiles("config.env")
// Result: {User: "env_user", Editor: "nvim"}

// With PreferFileOverEnv
gt.globalDecodeOpt.PreferFileOverEnv = true
gt.LoadConfigFiles("config.env")
// Result: {User: "file_user", Editor: "nvim"}
```

### Example 3: Nested Configuration

```go
// app.env
// DEBUG_C=true
// DB_USER=admin
// DB_SERVER_PORT=5432
// DB_POLING_MAX_POOL=100

type Database struct {
    User       string
    Server     string `config:"server_port"`
    PoolingMax int    `config:"poling_max_pool"`
}

type AppConfig struct {
    Debug    bool     `config:"debug_c"`
    Database Database `nested:"db"`
}

gt := gathuk.NewGathuk[AppConfig]()
err := gt.LoadConfigFiles("app.env")
// Result: {Debug: true, Database: {User: "admin", Server: "5432", PoolingMax: 100}}
```

## API Reference

### Core Functions

#### `NewGathuk[T any]() *Gathuk[T]`

Creates a new Gathuk instance with default configuration.

#### `LoadConfigFiles(srcFiles ...string) error`

Loads and merges configurations from one or more files.

#### `LoadConfig(src io.Reader, format string) error`

Loads configuration from an io.Reader with specified format.

#### `GetConfig() T`

Returns the parsed configuration struct.

#### `WriteConfigFile(dst string, mode fs.FileMode, config T) error`

Writes configuration to a file.

#### `WriteConfig(out io.Writer, format string, config T) error`

Writes configuration to an io.Writer.

#### `SetConfigFiles(srcFiles ...string)`

Sets base configuration files without loading them.

#### `SetCustomCodecRegistry(c option.CodecRegistry[T]) *Gathuk[T]`

Sets a custom codec registry for handling different file formats.

## Performance

Gathuk is designed for performance. Here are benchmark results:

```
goos: linux
goarch: amd64
pkg: github.com/ahyalfan/gathuk
cpu: AMD Ryzen 5 6600H with Radeon Graphics
BenchmarkGathuk/Benchmark_1_:_Simple_Load_Gathuk_config-12                 82419             12633 ns/op            3352 B/op         55 allocs/op
BenchmarkGathuk/Benchmark_2_:_Nested_Struct_Load_Gathuk_config-12         103670             11256 ns/op            3760 B/op         71 allocs/op
BenchmarkGathuk/Benchmark_3_:_Nested_Struct_Load_Gathuk_config_with_multiple_files-12              56586             21522 ns/op            6808 B/op         120 allocs/op
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Author

- [@ahyalfan](https://github.com/ahyalfan)

## Support

If you find this project helpful, please give it a ⭐️!

For issues and questions, please use the [GitHub issue tracker](https://github.com/ahyalfan/gathuk/issues).
