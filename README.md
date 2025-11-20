# Gathuk

![Go Version](https://img.shields.io/badge/Go-1.24.3-blue.svg)
![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ahyalfan/gathuk)](https://goreportcard.com/report/github.com/ahyalfan/gathuk)
[![GoDoc](https://godoc.org/github.com/ahyalfan/gathuk?status.svg)](https://godoc.org/github.com/ahyalfan/gathuk)

**Gathuk** is a type-safe, flexible configuration management library for Go that converts configuration files into strongly-typed structs. It supports multiple file formats (currently `.env` `.json`), nested structures, and automatic environment variable binding.

## Features

- 🎯 **Type-Safe**: Uses Go generics for compile-time type safety
- 📁 **Multiple File Formats**: Support for `.env`, `.json` (YAML, TOML coming soon)
- 🔄 **Multiple File Loading**: Load and merge configurations from multiple files
- 🔄 **Automatic Environment Variables**: Automatically bind OS environment variables to struct fields
- 🏗️ **Nested Structures**: Full support for nested struct configurations with custom prefixes
- 🔧 **Flexible Options**: Configure priority between file configs and environment variables
- 💾 **Write Support**: Export configurations back to files
- 🎨 **Custom Codecs**: Extensible codec system for adding new file formats
- 🚀 **Zero Dependencies**: Minimal external dependencies
- ⚡ **High Performance**: Optimized for speed with efficient parsing

## Installation

```bash
go get github.com/ahyalfan/gathuk
```

## Quick Start

### Basic Usage with .env

```go
package main

import (
    "fmt"
    "log"
    "github.com/ahyalfan/gathuk"
)

type Config struct {
    Port int
    Host string
}

func main() {
    // Create a new Gathuk instance
    gt := gathuk.NewGathuk[Config]()

    // Load configuration from file
    if err := gt.LoadConfigFiles(".env"); err != nil {
        log.Fatal(err)
    }

    // Get the parsed configuration
    config := gt.GetConfig()
    fmt.Printf("Server: %s:%d\n", config.Host, config.Port)
}
```

**`.env` file:**

```env
PORT=8080
HOST=localhost
```

### Basic Usage with JSON

```go
type Config struct {
    Port int    `config:"port"`
    Host string `config:"host"`
}

func main() {
    gt := gathuk.NewGathuk[Config]()

    if err := gt.LoadConfigFiles("config.json"); err != nil {
        log.Fatal(err)
    }

    config := gt.GetConfig()
    fmt.Printf("Server: %s:%d\n", config.Host, config.Port)
}
```

**`config.json` file:**

```json
{
  "port": 8080,
  "host": "localhost"
}
```

## Supported Formats

| Format                | Extension       | Status         | Tag Convention   |
| --------------------- | --------------- | -------------- | ---------------- |
| Environment Variables | `.env`          | ✅ Stable      | UPPER_SNAKE_CASE |
| JSON                  | `.json`         | ✅ Stable      | lower_snake_case |
| YAML                  | `.yaml`, `.yml` | 🚧 Coming Soon | lower_snake_case |
| TOML                  | `.toml`         | 🚧 Coming Soon | lower_snake_case |

### Format-Specific Behavior

#### .env Format

- Keys are automatically converted to `UPPER_SNAKE_CASE`
- Supports comments with `#`
- Example: `DATABASE_URL=postgres://localhost`

#### JSON Format

- Keys use `lower_snake_case` by default
- Supports nested objects and arrays
- Full JSON specification compliance
- Example: `{"database_url": "postgres://localhost"}`

## Struct Tags

Gathuk supports two main struct tags for customization:

### `config` Tag

Maps struct fields to specific configuration keys:

```go
type Config struct {
    // .env: SERVER_PORT | JSON: server_port
    Port int `config:"server_port"`

    // .env: API_KEY | JSON: api_key
    APIKey string `config:"api_key"`
}
```

### `nested` Tag

Defines prefix for nested structures:

```go
type Config struct {
    // All Database fields will have DB_ prefix in .env
    // In JSON: nested under "db" object
    Database Database `config:"db"`
    // or
    // Database Database `config:"db"`
}

type Database struct {
    Host string  // .env: DB_HOST | JSON: db.host
    Port int     // .env: DB_PORT | JSON: db.port
}
```

**Example `.env`:**

```env
DB_HOST=localhost
DB_PORT=5432
```

**Example JSON:**

```json
{
  "db": {
    "host": "localhost",
    "port": 5432
  }
}
```

### Ignoring Fields

Use `-` to exclude fields from configuration:

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

### Priority Examples

**Scenario 1: File Only (Default)**

```go
gt := gathuk.NewGathuk[Config]()
// Only reads from config.env
err := gt.LoadConfigFiles("config.env")
```

**Scenario 2: Environment Override**

```go
gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true
// Environment variables override file values
err := gt.LoadConfigFiles("config.env")
```

**Scenario 3: File Override**

```go
gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true
gt.globalDecodeOpt.PreferFileOverEnv = true
// File values override environment variables
err := gt.LoadConfigFiles("config.env")
```

## Multiple File Loading

Load and merge configurations from multiple files:

```go
gt := gathuk.NewGathuk[Config]()

// Method 1: Load multiple files at once
err := gt.LoadConfigFiles("base.env", "dev.env", "local.env")

// Method 2: Set base files, then load additional files
gt.SetConfigFiles("base.env", "defaults.env")
err := gt.LoadConfigFiles("override.env")

// Method 3: Mix different formats
err := gt.LoadConfigFiles("base.json", "override.env")
```

**Merge Behavior**:

- Later files override values from earlier files
- Zero values are not merged (existing non-zero values are preserved)
- Nested structs are merged recursively

### Example: Multi-Environment Setup

```go
// Load base config + environment-specific config
env := os.Getenv("APP_ENV") // "development", "staging", "production"
if env == "" {
    env = "development"
}

gt := gathuk.NewGathuk[Config]()
gt.SetConfigFiles("config/base.json")
err := gt.LoadConfigFiles(fmt.Sprintf("config/%s.json", env))
```

## Writing Configuration

Export your configuration to files:

```go
config := Config{
    Port: 8080,
    Host: "localhost",
}

// Write to .env file
err := gt.WriteConfigFile("output.env", 0644, config)

// Write to JSON file
err := gt.WriteConfigFile("output.json", 0644, config)

// Write to io.Writer
var buf bytes.Buffer
err := gt.WriteConfig(&buf, "json", config)
fmt.Println(buf.String())
```

## Advanced Usage

### Custom Codec Registry

Create and register custom codecs for different file formats:

```go
// Create custom codec
type JSONCodec[T any] struct {
    option.DefaultCodec[T]
}

func (c *JSONCodec[T]) Decode(buf []byte,val  *T) error {
    err := json.Unmarshal(buf, val)
    return  err
}

func (c *JSONCodec[T]) Encode(val T) ([]byte, error) {
    return json.Marshal(val)
}

// Register codec
func main() {
    registry := gathuk.NewDefaultCodecRegister[Config]()
    registry.RegisterCodec("json", &JSONCodec[Config]{})

    gt := gathuk.NewGathuk[Config]()
    gt.SetCustomCodecRegistry(registry)

    err := gt.LoadConfigFiles("config.json")
}
```

### Loading from io.Reader

```go
// From file
file, err := os.Open("config.env")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

gt := gathuk.NewGathuk[Config]()
err = gt.LoadConfig(file, "env")

// From HTTP response
resp, err := http.Get("https://api.example.com/config")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

err = gt.LoadConfig(resp.Body, "json")

// From string
configStr := `{"port": 8080, "host": "localhost"}`
reader := strings.NewReader(configStr)
err = gt.LoadConfig(reader, "json")
```

### Format-Specific Options

```go
gt := gathuk.NewGathuk[Config]()

// Set decode options for specific format
envOpt := &option.DecodeOption{
    AutomaticEnv: true,
    PreferFileOverEnv: true,
}
gt.SetDecodeOption("env", envOpt)

// JSON doesn't need env options
jsonOpt := &option.DecodeOption{
    AutomaticEnv: false,
}
gt.SetDecodeOption("json", jsonOpt)
```

## Field Name Mapping

### Environment Variables (.env)

By default, converts PascalCase to UPPER_SNAKE_CASE:

| Struct Field     | .env Variable     |
| ---------------- | ----------------- |
| `SimpleValue`    | `SIMPLE_VALUE`    |
| `DatabaseURL`    | `DATABASE_U_R_L`  |
| `APIKey`         | `A_P_I_KEY`       |
| `MaxConnections` | `MAX_CONNECTIONS` |

Use `config` tag for custom mapping:

```go
type Config struct {
    APIKey string `config:"api_key"`  // Maps to API_KEY
}
```

### JSON Format

By default, converts PascalCase to lower_snake_case:

| Struct Field     | JSON Key          |
| ---------------- | ----------------- |
| `SimpleValue`    | `simple_value`    |
| `DatabaseURL`    | `database_u_r_l`  |
| `APIKey`         | `a_p_i_key`       |
| `MaxConnections` | `max_connections` |

Use `config` tag for custom mapping:

```go
type Config struct {
    APIKey string `config:"api_key"`  // Maps to "api_key" in JSON
}
```

## Examples

### Example 1: Simple Configuration

```go
// .env file
// PORT=8080
// HOST=localhost

type Config struct {
    Port int
    Host string
}

gt := gathuk.NewGathuk[Config]()
err := gt.LoadConfigFiles(".env")
// Result: {Port: 8080, Host: "localhost"}
```

### Example 2: With Environment Variable Override

```go
// config.json
// {
//   "server": {
//     "port": 8080,
//     "host": "localhost"
//   },
//   "database": {
//     "host": "db.example.com",
//     "port": 5432
//   }
// }

type Server struct {
    Port int    `config:"port"`
    Host string `config:"host"`
}

type Database struct {
    Host string `config:"host"`
    Port int    `config:"port"`
}

type Config struct {
    Server   Server   `config:"server"`
    Database Database `config:"database"`
}
```

### Example 3: Environment Variable Override

```go
// config.env
// USER=file_user
// PORT=8080

type Config struct {
    User   string
    Port   int
    Editor string
}

// Set environment variables
os.Setenv("USER", "env_user")
os.Setenv("EDITOR", "nvim")

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true

err := gt.LoadConfigFiles("config.env")
// Result: {User: "env_user", Port: 8080, Editor: "nvim"}
// USER from env overrides file, EDITOR only in env, PORT from file

// With PreferFileOverEnv
gt.globalDecodeOpt.PreferFileOverEnv = true
err = gt.LoadConfigFiles("config.env")
// Result: {User: "file_user", Port: 8080, Editor: "nvim"}
// USER from file overrides env, EDITOR still from env
```

### Example 4: Multi-Format Configuration

```go
type Config struct {
    Server   ServerConfig   `config:"server"`
    Database DatabaseConfig `config:"database"`
}

gt := gathuk.NewGathuk[Config]()

// Load base config from JSON, override with .env
err := gt.LoadConfigFiles("config.json", "override.env")
```

### Example 5: Dynamic Configuration

```go
// Load based on environment
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}

configFiles := []string{
    "config/base.json",
    fmt.Sprintf("config/%s.json", env),
}

// Add local override if exists
if _, err := os.Stat("config/local.json"); err == nil {
    configFiles = append(configFiles, "config/local.json")
}

gt := gathuk.NewGathuk[Config]()
err := gt.LoadConfigFiles(configFiles...)
```

## Performance

Gathuk is optimized for performance with efficient parsing and minimal allocations.

### Benchmark Results

```
goos: linux
goarch: amd64
pkg: github.com/ahyalfan/gathuk
cpu: AMD Ryzen 5 6600H with Radeon Graphics

BenchmarkGathuk/Simple_Load-12                  116638    10046 ns/op    3240 B/op    50 allocs/op
BenchmarkGathuk/Nested_Struct-12                113784    10685 ns/op    3432 B/op    62 allocs/op
BenchmarkGathuk/Multiple_Files-12                57288    20465 ns/op    6152 B/op   102 allocs/op
```

### Run Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkGathuk/Simple -benchmem

# With CPU profiling
go test -bench=. -benchmem -cpuprofile=cpu.prof

# With memory profiling
go test -bench=. -benchmem -memprofile=mem.prof
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

#### `SetDecodeOption(format string, opt *option.DecodeOption)`

Sets decode options for a specific format.

#### `SetEncodeOption(format string, opt *option.EncodeOption)`

Sets encode options for a specific format.

### For complete API documentation, see [GoDoc](https://godoc.org/github.com/ahyalfan/gathuk)

## Best Practices

### 1. Use Environment-Specific Configs

```go
// config/
//   base.json       - Common settings
//   development.json - Dev overrides
//   staging.json    - Staging overrides
//   production.json - Production overrides
//   local.json      - Local developer overrides (gitignored)
```

### 2. Validate Configuration

```go
config := gt.GetConfig()

// Validate required fields
if config.Database.Host == "" {
    log.Fatal("Database host is required")
}

// Validate ranges
if config.Server.Port < 1 || config.Server.Port > 65535 {
    log.Fatal("Invalid port number")
}
```

### 3. Use Struct Tags Consistently

```go
// Good: Explicit mapping
type Config struct {
    ServerPort int `config:"server_port"`
    DBHost     string `config:"db_host"`
}

// Avoid: Mixing conventions
type Config struct {
    ServerPort int              // Auto-mapped
    DBHost     string `config:"database_host"` // Custom
}
```

### 4. Document Your Configuration

```go
type Config struct {
    // Server port (default: 8080)
    ServerPort int `config:"server_port"`

    // Maximum number of database connections (default: 100)
    MaxConnections int `config:"max_connections"`
}
```

## Migration Guide

### From Viper

```go
// Before (Viper)
viper.SetConfigName("config")
viper.SetConfigType("json")
viper.AddConfigPath(".")
err := viper.ReadInConfig()
port := viper.GetInt("server.port")

// After (Gathuk)
type Config struct {
    Server struct {
        Port int `config:"port"`
    } `config:"server"`
}

gt := gathuk.NewGathuk[Config]()
err := gt.LoadConfigFiles("config.json")
port := gt.GetConfig().Server.Port
```

### From godotenv

```go
// Before (godotenv)
err := godotenv.Load()
port := os.Getenv("PORT")

// After (Gathuk)
type Config struct {
    Port int
}

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true
err := gt.LoadConfigFiles(".env")
port := gt.GetConfig().Port
```

## FAQ

**Q: Can I use multiple formats simultaneously?**  
A: Yes! You can load different formats in sequence: `gt.LoadConfigFiles("base.json", "override.env")`

**Q: How do I handle missing configuration files?**  
A: Check for `os.IsNotExist(err)` and provide defaults or use fallback files.

**Q: Can I reload configuration at runtime?**  
A: Yes, call `LoadConfigFiles()` again. Values will be merged with existing configuration.

**Q: Does it support configuration validation?**  
A: Validate after loading using your own validation logic or libraries like `go-playground/validator`.

**Q: How do I set default values?**  
A: Initialize your struct with defaults before loading: `config := Config{Port: 8080}`

**Q: Can I use with Docker/Kubernetes?**  
A: Yes! Use `AutomaticEnv` to read from environment variables set by orchestration tools.

**Q: Is it thread-safe?**  
A: Reading config after loading is thread-safe. Loading config should be done during initialization.

## Roadmap

- [x] .env format support
- [x] JSON format support
- [x] Environment variable binding
- [x] Multiple file merging
- [x] Nested structure support
- [x] Write support
- [ ] YAML format support
- [ ] TOML format support
- [ ] Configuration validation
- [ ] Hot reload support
- [ ] Configuration encryption
- [ ] Remote config sources (etcd, consul)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/ahyalfan/gathuk.git
cd gathuk

# Run tests
go test ./...

# Run benchmarks
go test -bench=. -benchmem

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Contribution Guidelines

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Write tests for your changes
4. Ensure all tests pass (`go test ./...`)
5. Run `go fmt ./...` to format code
6. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
7. Push to the branch (`git push origin feature/AmazingFeature`)
8. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Author

- [@ahyalfan](https://github.com/ahyalfan)

## Acknowledgments

- Inspired by [Viper](https://github.com/spf13/viper) and [godotenv](https://github.com/joho/godotenv)
- Thanks to all contributors

## Support

If you find this project helpful, please give it a ⭐️!

For issues and questions, please use the [GitHub issue tracker](https://github.com/ahyalfan/gathuk/issues).
