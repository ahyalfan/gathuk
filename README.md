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

**Requirements:**

- Go 1.21 or higher (for generics support)

**Verify installation:**

```bash
go list -m github.com/ahyalfan/gathuk
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

## Core Concepts

### 1. Generic Type Parameter

Gathuk uses Go generics to provide type-safe configuration loading:

```go
// Concrete struct type (RECOMMENDED)
gt := gathuk.NewGathuk[Config]()

// Generic any type (LIMITED - see warnings)
gt := gathuk.NewGathuk[any]()

// Map type (LIMITED - see warnings)
gt := gathuk.NewGathuk[map[string]any]()
```

**Always prefer concrete struct types for:**

- ✅ Type safety at compile time
- ✅ Proper merging when loading multiple files
- ✅ Better IDE support and autocomplete
- ✅ Self-documenting code

### 2. Configuration Loading Flow

```
Config Files → Tokenize → Parse → Decode → Struct
                                  ↓
                          Environment Variables
                                  ↓
                           Merge & Apply Options
                                  ↓
                            Final Config Struct
```

### 3. Field Name Mapping

Gathuk automatically converts field names to appropriate conventions:

| Go Field Name | .env Format      | JSON Format      |
| ------------- | ---------------- | ---------------- |
| `Port`        | `PORT`           | `port`           |
| `ServerPort`  | `SERVER_PORT`    | `server_port`    |
| `DatabaseURL` | `DATABASE_U_R_L` | `database_u_r_l` |
| `APIKey`      | `A_P_I_KEY`      | `a_p_i_key`      |

**Override with tags:**

```go
type Config struct {
    APIKey string `config:"api_key"`  // → API_KEY or api_key
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

**Features:**

- Simple key-value pairs: `KEY=value`
- Comments start with `#`
- Keys automatically converted to UPPER_SNAKE_CASE
- No quotes needed for string values
- Inline comments supported: `PORT=8080 # server port`

**Example:**

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp

# Feature Flags
DEBUG=true
ENABLE_LOGGING=true
```

**Supported Types:**

- `string`: Direct text
- `int`, `int64`: Integers
- `float64`: Floating-point numbers
- `bool`: `true` or `false`

#### JSON Format

**Features:**

- Full JSON specification compliance
- Nested objects and arrays
- Keys use lower_snake_case by default
- Type-safe parsing
- Pretty-print support for writing

**Example:**

```json
{
  "server": {
    "port": 8080,
    "host": "localhost",
    "timeout": 30
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "credentials": {
      "username": "admin",
      "password": "secret"
    }
  },
  "features": {
    "debug": true,
    "cache_enabled": true
  }
}
```

**Supported Types:**

- All primitive types (string, number, boolean, null)
- Objects (nested structs)
- Arrays (slices)
- Mixed arrays with `[]interface{}`

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

### Merge Behavior

**Files are processed sequentially:**

1. First file loaded → Initial config
2. Second file loaded → Merged with first
3. Third file loaded → Merged with result of 1+2
4. Continue...

**Merge rules:**

- ✅ Non-zero values from later files **override** earlier files
- ❌ Zero values from later files **do NOT override** earlier files
- ✅ New fields from later files **are added**
- ✅ Nested structs **merge recursively**

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

### Zero Value Behavior

**IMPORTANT:** Zero values are NOT merged to prevent accidental clearing:

```go
// base.env
PORT=8080
HOST=localhost
MAX_CONNECTIONS=100

// override.env
PORT=0         # Zero value - IGNORED
HOST=          # Empty string - IGNORED
MAX_CONNECTIONS=50  # Non-zero - USED
```

```go
gt := gathuk.NewGathuk[Config]()
err := gt.LoadConfigFiles("base.env", "override.env")

config := gt.GetConfig()
// Result:
// Port:           8080 (NOT overridden by 0)
// Host:           "localhost" (NOT overridden by "")
// MaxConnections: 50 (overridden by non-zero)
```

**Rationale:** This prevents accidentally clearing important configuration values with empty or zero values in override files.

### Environment-Specific Loading

```go
func LoadConfig() (*Config, error) {
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    files := []string{
        "config/base.env",                     // Always loaded
        fmt.Sprintf("config/%s.env", env),     // Environment-specific
    }

    // Add local overrides if exists
    localFile := "config/local.env"
    if _, err := os.Stat(localFile); err == nil {
        files = append(files, localFile)
    }

    gt := gathuk.NewGathuk[Config]()
    if err := gt.LoadConfigFiles(files...); err != nil {
        return nil, err
    }

    return &config, nil
}
```

**Directory structure:**

```
config/
  ├── base.env          # Common settings
  ├── development.env   # Dev overrides
  ├── staging.env       # Staging overrides
  ├── production.env    # Production overrides
  └── local.env         # Local dev (gitignored)
```

### Priority Examples

#### Example 1: Environment Only

```go
os.Setenv("PORT", "9000")
os.Setenv("HOST", "0.0.0.0")

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true

// No files - only environment
err := gt.LoadConfigFiles()

config := gt.GetConfig()
// Port: 9000, Host: "0.0.0.0"
```

#### Example 2: File + Environment (Env Wins)

```env
# config.env
PORT=8080
HOST=localhost
```

```go
os.Setenv("PORT", "9000") // This will win

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true
err := gt.LoadConfigFiles("config.env")

config := gt.GetConfig()
// Port: 9000 (from env), Host: "localhost" (from file)
```

#### Example 3: File + Environment (File Wins)

```go
os.Setenv("PORT", "9000") // This will be ignored

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true
gt.globalDecodeOpt.PreferFileOverEnv = true
err := gt.LoadConfigFiles("config.env")

config := gt.GetConfig()
// Port: 8080 (from file), Host: "localhost" (from file)
```

#### Example 4: Partial Environment

```env
# config.env
PORT=8080
HOST=localhost
```

```go
os.Setenv("DEBUG", "true")      // Additional env var
os.Setenv("LOG_LEVEL", "info")  // Additional env var

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true
err := gt.LoadConfigFiles("config.env")

config := gt.GetConfig()
// Port: 8080 (file), Host: "localhost" (file)
// Debug: true (env), LogLevel: "info" (env)
```

### Docker/Kubernetes Integration

Gathuk works seamlessly with containerized deployments:

```go
type Config struct {
    Port         int    `config:"port"`
    DatabaseURL  string `config:"database_url"`
    RedisURL     string `config:"redis_url"`
}

func main() {
    gt := gathuk.NewGathuk[Config]()
    gt.globalDecodeOpt.AutomaticEnv = true

    // In Docker/K8s, all config comes from environment
    // Set via docker-compose.yml, Dockerfile ENV, or K8s ConfigMap
    err := gt.LoadConfigFiles()

    config := gt.GetConfig()
    // Ready to use!
}
```

**docker-compose.yml:**

```yaml
services:
  app:
    environment:
      - PORT=8080
      - DATABASE_URL=postgres://localhost:5432/db
      - REDIS_URL=redis://localhost:6379
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

### Validation After Loading

```go
type Config struct {
    Port     int    `config:"port"`
    Host     string `config:"host"`
    LogLevel string `config:"log_level"`
}

func (c *Config) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }

    if c.Host == "" {
        return fmt.Errorf("host is required")
    }

    validLevels := map[string]bool{
        "debug": true, "info": true, "warn": true, "error": true,
    }
    if !validLevels[c.LogLevel] {
        return fmt.Errorf("invalid log level: %s", c.LogLevel)
    }

    return nil
}

func main() {
    gt := gathuk.NewGathuk[Config]()
    err := gt.LoadConfigFiles("config.env")
    if err != nil {
        log.Fatal(err)
    }

    config := gt.GetConfig()
    if err := config.Validate(); err != nil {
        log.Fatal("Config validation failed:", err)
    }
}
```

### Configuration Reloading

```go
type App struct {
    config *Config
    gt     *gathuk.Gathuk[Config]
}

func (app *App) ReloadConfig() error {
    if err := app.gt.LoadConfigFiles("config.env"); err != nil {
        return err
    }

    newConfig := app.gt.GetConfig()

    // Validate before applying
    if err := newConfig.Validate(); err != nil {
        return err
    }

    // Atomic update
    app.config = &newConfig

    log.Println("Configuration reloaded successfully")
    return nil
}

// Reload on signal
func (app *App) WatchConfig() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGHUP)

    for {
        <-sigChan
        if err := app.ReloadConfig(); err != nil {
            log.Printf("Failed to reload config: %v", err)
        }
    }
}
```

## Important Warnings

### ⚠️ Warning 1: Generic Type `any` Behavior

**CRITICAL:** When using `any` or `map[string]any`, multiple file loading does NOT merge:

```go
// ❌ WRONG: Files are NOT merged!
gt := gathuk.NewGathuk[any]()
gt.LoadConfigFiles("base.env", "dev.env")
// Only dev.env values are kept!
// All base.env values are LOST!

// ✅ CORRECT: Use struct type for merging
type Config struct {
    Port int
    Host string
}
gt := gathuk.NewGathuk[Config]()
gt.LoadConfigFiles("base.env", "dev.env")
// Properly merged!
```

**Why?**

- Struct types: Gathuk knows which fields to merge
- `any`/`map`: Gathuk sees generic map, replaces entirely
- Each load creates new map, discarding previous

**Solutions:**

1. **Use concrete struct types** (recommended)
2. **Load files separately** and merge manually
3. **Load single file** at a time

See [Multiple Files & Merging Documentation](docs/multiple-files.md) for details.

### ⚠️ Warning 2: Zero Values

Zero values from later files do NOT override earlier files:

```go
// base.env
PORT=8080

// override.env
PORT=0  # Will NOT override!

gt.LoadConfigFiles("base.env", "override.env")
// Result: Port = 8080 (not 0)
```

**To force zero values:**

- Load only the override file
- Use a non-zero sentinel value
- Manually set after loading

### ⚠️ Warning 3: Field Names with Acronyms

```go
type Config struct {
    APIKey string  // → A_P_I_KEY (not API_KEY)
    HTTPURL string // → H_T_T_P_U_R_L (not HTTP_URL)
}

// Use tags for better names
type Config struct {
    APIKey string `config:"api_key"`  // → API_KEY
    HTTPURL string `config:"http_url"` // → HTTP_URL
}
```

### ⚠️ Warning 4: Concurrent Access

```go
// ❌ NOT safe for concurrent access during load
gt := gathuk.NewGathuk[Config]()

go gt.LoadConfigFiles("config1.env") // Unsafe!
go gt.LoadConfigFiles("config2.env") // Unsafe!

// ✅ Safe: Load once, read concurrently
gt.LoadConfigFiles("config.env")
config := gt.GetConfig()

go func() { use(config) }() // Safe
go func() { use(config) }() // Safe
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

## Complete Examples

### Example 1: Web Server Configuration

```go
type Config struct {
    Server   ServerConfig   `nested:"server"`
    Database DatabaseConfig `nested:"db"`
    Redis    RedisConfig    `nested:"redis"`
    Logging  LogConfig      `nested:"log"`
}

type ServerConfig struct {
    Port         int    `config:"port"`
    Host         string `config:"host"`
    ReadTimeout  int    `config:"read_timeout"`
    WriteTimeout int    `config:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `config:"host"`
    Port     int    `config:"port"`
    User     string `config:"user"`
    Password string `config:"password"`
    Database string `config:"name"`
    MaxConns int    `config:"max_connections"`
}

type RedisConfig struct {
    Host     string `config:"host"`
    Port     int    `config:"port"`
    Password string `config:"password"`
    DB       int    `config:"db"`
}

type LogConfig struct {
    Level  string `config:"level"`
    Format string `config:"format"`
}

func main() {
    // Load configuration
    gt := gathuk.NewGathuk[Config]()
    gt.globalDecodeOpt.AutomaticEnv = true

    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    files := []string{
        "config/base.env",
        fmt.Sprintf("config/%s.env", env),
    }

    if err := gt.LoadConfigFiles(files...); err != nil {
        log.Fatal("Failed to load config:", err)
    }

    config := gt.GetConfig()

    // Validate configuration
    if config.Server.Port < 1 || config.Server.Port > 65535 {
        log.Fatal("Invalid server port")
    }

    // Start server
    addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
    log.Printf("Starting server on %s", addr)

    // Use configuration
    db := connectDatabase(config.Database)
    redis := connectRedis(config.Redis)

    server := &http.Server{
        Addr:         addr,
        ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
    }

    log.Fatal(server.ListenAndServe())
}
```

**`config/base.env`:**

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=app
DB_PASSWORD=secret
DB_NAME=myapp
DB_MAX_CONNECTIONS=25

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

**`config/development.env`:**

```env
# Override for development
SERVER_PORT=3000
DB_MAX_CONNECTIONS=5
LOG_LEVEL=debug
LOG_FORMAT=text
```

### Example 2: Microservice Configuration

```go
type Config struct {
    Service     ServiceConfig     `nested:"service"`
    HTTP        HTTPConfig        `nested:"http"`
    GRPC        GRPCConfig        `nested:"grpc"`
    Observability ObservabilityConfig `nested:"obs"`
    Dependencies  DependenciesConfig  `nested:"deps"`
}

type ServiceConfig struct {
    Name        string `config:"name"`
    Version     string `config:"version"`
    Environment string `config:"environment"`
}

type HTTPConfig struct {
    Enabled bool   `config:"enabled"`
    Port    int    `config:"port"`
    Timeout int    `config:"timeout"`
}

type GRPCConfig struct {
    Enabled bool   `config:"enabled"`
    Port    int    `config:"port"`
    Timeout int    `config:"timeout"`
}

type ObservabilityConfig struct {
    Metrics   MetricsConfig   `nested:"metrics"`
    Tracing   TracingConfig   `nested:"tracing"`
    Logging   LoggingConfig   `nested:"logging"`
}

type MetricsConfig struct {
    Enabled bool   `config:"enabled"`
    Port    int    `config:"port"`
    Path    string `config:"path"`
}

type TracingConfig struct {
    Enabled  bool   `config:"enabled"`
    Endpoint string `config:"endpoint"`
    SampleRate float64 `config:"sample_rate"`
}

type LoggingConfig struct {
    Level  string `config:"level"`
    Format string `config:"format"`
}

type DependenciesConfig struct {
    Database DatabaseConfig `nested:"db"`
    Cache    CacheConfig    `nested:"cache"`
    Queue    QueueConfig    `nested:"queue"`
}

type DatabaseConfig struct {
    Host         string `config:"host"`
    Port         int    `config:"port"`
    User         string `config:"user"`
    Password     string `config:"password"`
    Database     string `config:"name"`
    MaxOpenConns int    `config:"max_open_conns"`
    MaxIdleConns int    `config:"max_idle_conns"`
}

type CacheConfig struct {
    Host     string `config:"host"`
    Port     int    `config:"port"`
    Password string `config:"password"`
    TTL      int    `config:"ttl"`
}

type QueueConfig struct {
    URL          string `config:"url"`
    MaxRetries   int    `config:"max_retries"`
    RetryDelay   int    `config:"retry_delay"`
}

func LoadConfig() (*Config, error) {
    gt := gathuk.NewGathuk[Config]()
    gt.globalDecodeOpt.AutomaticEnv = true

    // Load base + environment-specific config
    env := os.Getenv("SERVICE_ENVIRONMENT")
    if env == "" {
        env = "development"
    }

    files := []string{
        "config/base.env",
        fmt.Sprintf("config/%s.env", env),
    }

    // Add secrets file if exists (for local development)
    if _, err := os.Stat("config/secrets.env"); err == nil {
        files = append(files, "config/secrets.env")
    }

    if err := gt.LoadConfigFiles(files...); err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    config := gt.GetConfig()

    // Validate
    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }

    return &config, nil
}

func validateConfig(cfg *Config) error {
    if cfg.Service.Name == "" {
        return fmt.Errorf("service name is required")
    }

    if !cfg.HTTP.Enabled && !cfg.GRPC.Enabled {
        return fmt.Errorf("at least one protocol (HTTP or GRPC) must be enabled")
    }

    if cfg.HTTP.Enabled && (cfg.HTTP.Port < 1 || cfg.HTTP.Port > 65535) {
        return fmt.Errorf("invalid HTTP port: %d", cfg.HTTP.Port)
    }

    if cfg.GRPC.Enabled && (cfg.GRPC.Port < 1 || cfg.GRPC.Port > 65535) {
        return fmt.Errorf("invalid GRPC port: %d", cfg.GRPC.Port)
    }

    return nil
}
```

### Example 3: CLI Application Configuration

```go
type Config struct {
    App      AppConfig      `nested:"app"`
    API      APIConfig      `nested:"api"`
    Output   OutputConfig   `nested:"output"`
    Advanced AdvancedConfig `nested:"advanced"`
}

type AppConfig struct {
    Name    string `config:"name"`
    Version string `config:"version"`
    Debug   bool   `config:"debug"`
}

type APIConfig struct {
    BaseURL string `config:"base_url"`
    Token   string `config:"token"`
    Timeout int    `config:"timeout"`
}

type OutputConfig struct {
    Format string `config:"format"` // json, yaml, table
    Color  bool   `config:"color"`
    Quiet  bool   `config:"quiet"`
}

type AdvancedConfig struct {
    CacheDir     string `config:"cache_dir"`
    MaxRetries   int    `config:"max_retries"`
    RetryDelay   int    `config:"retry_delay"`
}

func main() {
    // Parse flags
    configFile := flag.String("config", "", "Config file path")
    debug := flag.Bool("debug", false, "Enable debug mode")
    flag.Parse()

    // Load configuration
    gt := gathuk.NewGathuk[Config]()
    gt.globalDecodeOpt.AutomaticEnv = true

    files := []string{}

    // Load from default locations
    homeDir, _ := os.UserHomeDir()
    defaultFiles := []string{
        filepath.Join(homeDir, ".myapp", "config.env"),
        ".myapp.env",
    }

    for _, f := range defaultFiles {
        if _, err := os.Stat(f); err == nil {
            files = append(files, f)
        }
    }

    // Load from specified config file
    if *configFile != "" {
        files = append(files, *configFile)
    }

    if len(files) > 0 {
        if err := gt.LoadConfigFiles(files...); err != nil {
            log.Fatal("Failed to load config:", err)
        }
    }

    config := gt.GetConfig()

    // Override with flags
    if *debug {
        config.App.Debug = true
    }

    // Use configuration
    runCLI(config)
}
```

### Example 4: Testing Configuration

```go
func TestConfigLoading(t *testing.T) {
    tests := []struct {
        name     string
        envFile  string
        envVars  map[string]string
        want     Config
        wantErr  bool
    }{
        {
            name:    "basic config",
            envFile: "testdata/basic.env",
            want: Config{
                Port: 8080,
                Host: "localhost",
            },
            wantErr: false,
        },
        {
            name:    "with environment override",
            envFile: "testdata/basic.env",
            envVars: map[string]string{
                "PORT": "9000",
            },
            want: Config{
                Port: 9000,
                Host: "localhost",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set environment variables
            for k, v := range tt.envVars {
                os.Setenv(k, v)
                defer os.Unsetenv(k)
            }

            // Load config
            gt := gathuk.NewGathuk[Config]()
            gt.globalDecodeOpt.AutomaticEnv = true

            err := gt.LoadConfigFiles(tt.envFile)
            if (err != nil) != tt.wantErr {
                t.Errorf("LoadConfigFiles() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            got := gt.GetConfig()
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetConfig() = %v, want %v", got, tt.want)
            }
        })
    }
}
func TestConfigMerging(t *testing.T) {
    // Create temporary files
    baseFile := createTempFile(t, `
PORT=8080
HOST=localhost
DEBUG=false
`)
    defer os.Remove(baseFile)

    devFile := createTempFile(t, `
DEBUG=true
LOG_LEVEL=debug
`)
    defer os.Remove(devFile)

    // Load and merge
    gt := gathuk.NewGathuk[Config]()
    err := gt.LoadConfigFiles(baseFile, devFile)
    if err != nil {
        t.Fatal(err)
    }

    config := gt.GetConfig()

    // Verify merged results
    if config.Port != 8080 {
        t.Errorf("Port = %d, want 8080", config.Port)
    }
    if config.Host != "localhost" {
        t.Errorf("Host = %s, want localhost", config.Host)
    }
    if config.Debug != true {
        t.Errorf("Debug = %v, want true", config.Debug)
    }
    if config.LogLevel != "debug" {
        t.Errorf("LogLevel = %s, want debug", config.LogLevel)
    }
}

func createTempFile(t *testing.T, content string) string {
    t.Helper()
    file, err := os.CreateTemp("", "config-*.env")
    if err != nil {
        t.Fatal(err)
    }
    if _, err := file.WriteString(content); err != nil {
        t.Fatal(err)
    }
    file.Close()
    return file.Name()
}

```

### Example 5: Dynamic Configuration Switching

```go
type ConfigManager struct {
    configs map[string]*Config
    active  string
    mu      sync.RWMutex
}

func NewConfigManager() *ConfigManager {
    return &ConfigManager{
        configs: make(map[string]*Config),
        active:  "default",
    }
}

func (cm *ConfigManager) LoadProfile(name, file string) error {
    gt := gathuk.NewGathuk[Config]()
    gt.globalDecodeOpt.AutomaticEnv = true

    if err := gt.LoadConfigFiles(file); err != nil {
        return fmt.Errorf("failed to load profile %s: %w", name, err)
    }

    config := gt.GetConfig()

    cm.mu.Lock()
    cm.configs[name] = &config
    cm.mu.Unlock()

    return nil
}

func (cm *ConfigManager) SwitchProfile(name string) error {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    if _, exists := cm.configs[name]; !exists {
        return fmt.Errorf("profile %s not found", name)
    }

    cm.active = name
    log.Printf("Switched to profile: %s", name)
    return nil
}

func (cm *ConfigManager) GetConfig() Config {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    return *cm.configs[cm.active]
}

func main() {
    manager := NewConfigManager()

    // Load multiple profiles
    profiles := map[string]string{
        "development": "config/dev.env",
        "staging":     "config/staging.env",
        "production":  "config/prod.env",
    }

    for name, file := range profiles {
        if err := manager.LoadProfile(name, file); err != nil {
            log.Printf("Warning: %v", err)
        }
    }

    // Use active profile
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    if err := manager.SwitchProfile(env); err != nil {
        log.Fatal(err)
    }

    config := manager.GetConfig()
    log.Printf("Running with config: %+v", config)
}
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

**What this means:**

- **Simple Load**: ~10 microseconds per operation
- **Nested Struct**: ~11 microseconds per operation
- **Multiple Files**: ~20 microseconds per operation (loading 2 files)

**Memory efficiency:**

- Simple config: ~3.2 KB per load
- Nested struct: ~3.4 KB per load
- Multiple files: ~6.1 KB per load

### Performance Tips

1. **Reuse Gathuk instances when possible**

```go
// ✅ Good: Reuse instance
gt := gathuk.NewGathuk[Config]()
for _, file := range files {
    gt.LoadConfigFiles(file)
}

// ❌ Avoid: Creating new instance each time
for _, file := range files {
    gt := gathuk.NewGathuk[Config]() // Unnecessary allocation
    gt.LoadConfigFiles(file)
}
```

2. **Load once, use many times**

```go
// ✅ Good: Load once at startup
var globalConfig Config

func init() {
    gt := gathuk.NewGathuk[Config]()
    gt.LoadConfigFiles("config.env")
    globalConfig = gt.GetConfig()
}

func handler1() {
    // Use globalConfig
}

func handler2() {
    // Use globalConfig
}
```

3. **Use concrete struct types**

```go
// ✅ Good: Concrete type (faster)
gt := gathuk.NewGathuk[Config]()

// ❌ Slower: Generic any type
gt := gathuk.NewGathuk[any]()
```

## Best Practices

### 1. Always Use Struct Types for Multiple Files

```go
// ✅ DO: Use concrete struct for merging
type Config struct {
    Port int
    Host string
}
gt := gathuk.NewGathuk[Config]()
gt.LoadConfigFiles("base.env", "dev.env")

// ❌ DON'T: Use any with multiple files
gt := gathuk.NewGathuk[any]()
gt.LoadConfigFiles("base.env", "dev.env") // Only last file kept!
```

### 2. Organize Config Files by Environment

```
config/
  ├── base.env          # Common settings (all environments)
  ├── development.env   # Dev-specific overrides
  ├── staging.env       # Staging-specific overrides
  ├── production.env    # Production-specific overrides
  ├── secrets.env       # Secrets (gitignored, optional)
  └── local.env         # Local dev overrides (gitignored)
```

### 3. Validate Configuration After Loading

```go
type Config struct {
    Port     int
    Host     string
    LogLevel string
}

func (c *Config) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }

    if c.Host == "" {
        return fmt.Errorf("host is required")
    }

    validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !validLevels[c.LogLevel] {
        return fmt.Errorf("invalid log level: %s", c.LogLevel)
    }

    return nil
}

func main() {
    gt := gathuk.NewGathuk[Config]()
    gt.LoadConfigFiles("config.env")

    config := gt.GetConfig()
    if err := config.Validate(); err != nil {
        log.Fatal("Configuration error:", err)
    }
}
```

### 4. Use Struct Tags Consistently

```go
// ✅ Good: Explicit and consistent
type Config struct {
    ServerPort int    `config:"server_port"`
    DBHost     string `config:"db_host"`
    APIKey     string `config:"api_key"`
}

// ❌ Avoid: Mixing conventions
type Config struct {
    ServerPort int              // Auto-mapped
    DBHost     string `config:"database_host"` // Custom
    APIKey     string            // Auto-mapped (becomes A_P_I_KEY!)
}
```

### 5. Document Your Configuration

```go
type Config struct {
    // Server listening port (default: 8080, range: 1-65535)
    Port int `config:"port"`

    // Server bind address (default: "localhost")
    // Use "0.0.0.0" to listen on all interfaces
    Host string `config:"host"`

    // Maximum number of database connections (default: 100)
    MaxConnections int `config:"max_connections"`

    // Enable debug logging (default: false)
    // WARNING: Debug logging may expose sensitive information
    Debug bool `config:"debug"`
}
```

### 6. Use Environment Variables for Secrets

```go
// ❌ Don't: Store secrets in config files
// config.env (committed to git)
DATABASE_PASSWORD=mysecret123

// ✅ Do: Use environment variables for secrets
// config.env (committed to git)
DATABASE_HOST=localhost
DATABASE_PORT=5432

// Set secrets via environment
export DATABASE_PASSWORD=mysecret123

// Or use separate secrets file (gitignored)
// config/secrets.env (gitignored)
DATABASE_PASSWORD=mysecret123
```

### 7. Provide Sensible Defaults

```go
type Config struct {
    Port         int    `config:"port"`
    Host         string `config:"host"`
    ReadTimeout  int    `config:"read_timeout"`
    WriteTimeout int    `config:"write_timeout"`
}

func LoadConfigWithDefaults() (*Config, error) {
    // Set defaults
    config := Config{
        Port:         8080,
        Host:         "localhost",
        ReadTimeout:  30,
        WriteTimeout: 30,
    }

    // Override with file values
    gt := gathuk.NewGathuk[Config]()
    gt.globalDecodeOpt.AutomaticEnv = true

    // LoadConfigFiles will only override non-zero values
    if err := gt.LoadConfigFiles("config.env"); err != nil {
        // If config file doesn't exist, use defaults
        if !os.IsNotExist(err) {
            return nil, err
        }
    } else {
        config = gt.GetConfig()
    }

    return &config, nil
}
```

### 8. Test Configuration Loading

```go
func TestConfigLoading(t *testing.T) {
    // Create test config file
    content := `
PORT=8080
HOST=localhost
DEBUG=true
`
    tmpfile, err := os.CreateTemp("", "config-*.env")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpfile.Name())

    if _, err := tmpfile.WriteString(content); err != nil {
        t.Fatal(err)
    }
    tmpfile.Close()

    // Load config
    gt := gathuk.NewGathuk[Config]()
    if err := gt.LoadConfigFiles(tmpfile.Name()); err != nil {
        t.Fatal(err)
    }

    config := gt.GetConfig()

    // Assert values
    if config.Port != 8080 {
        t.Errorf("Port = %d, want 8080", config.Port)
    }
    if config.Host != "localhost" {
        t.Errorf("Host = %s, want localhost", config.Host)
    }
    if config.Debug != true {
        t.Errorf("Debug = %v, want true", config.Debug)
    }
}
```

### 9. Handle Missing Files Gracefully

```go
func LoadConfig(files ...string) (*Config, error) {
    gt := gathuk.NewGathuk[Config]()
    gt.globalDecodeOpt.AutomaticEnv = true

    // Filter existing files
    existingFiles := []string{}
    for _, file := range files {
        if _, err := os.Stat(file); err == nil {
            existingFiles = append(existingFiles, file)
        } else {
            log.Printf("Config file not found (skipping): %s", file)
        }
    }

    if len(existingFiles) == 0 {
        return nil, fmt.Errorf("no config files found")
    }

    if err := gt.LoadConfigFiles(existingFiles...); err != nil {
        return nil, err
    }

    config := gt.GetConfig()
    return &config, nil
}
```

### 10. Use Feature Flags Pattern

```go
type Config struct {
    Features FeatureFlags `nested:"feature"`
}

type FeatureFlags struct {
    EnableNewUI      bool `config:"new_ui"`
    EnableBetaAPI    bool `config:"beta_api"`
    EnableCaching    bool `config:"caching"`
}

// Load base config with all features disabled
// Then override based on environment

// base.env
FEATURE_NEW_UI=false
FEATURE_BETA_API=false
FEATURE_CACHING=true

// production.env
FEATURE_CACHING=true

// development.env
FEATURE_NEW_UI=true
FEATURE_BETA_API=true
FEATURE_CACHING=false
```

## Migration Guide

### From Viper

**Before (Viper):**

```go
import "github.com/spf13/viper"

viper.SetConfigName("config")
viper.SetConfigType("json")
viper.AddConfigPath(".")
viper.AutomaticEnv()

if err := viper.ReadInConfig(); err != nil {
    log.Fatal(err)
}

port := viper.GetInt("server.port")
host := viper.GetString("server.host")
```

**After (Gathuk):**

```go
import "github.com/ahyalfan/gathuk"

type Config struct {
    Server struct {
        Port int    `config:"port"`
        Host string `config:"host"`
    } `config:"server"`
}

gt := gathuk.NewGathuk[Config]()
gt.globalDecodeOpt.AutomaticEnv = true

if err := gt.LoadConfigFiles("config.json"); err != nil {
    log.Fatal(err)
}

config := gt.GetConfig()
port := config.Server.Port
host := config.Server.Host
```

**Benefits:**

- ✅ Type-safe access (no Get\* methods)
- ✅ Compile-time checking
- ✅ Better IDE support
- ✅ No string keys to remember

### From godotenv

**Before (godotenv):**

```go
import "github.com/joho/godotenv"

if err := godotenv.Load(); err != nil {
    log.Fatal(err)
}

port, _ := strconv.Atoi(os.Getenv("PORT"))
host := os.Getenv("HOST")
debug := os.Getenv("DEBUG") == "true"
```

**After (Gathuk):**

```go
import "github.com/ahyalfan/gathuk"

type Config struct {
    Port  int
    Host  string
    Debug bool
}

gt := gathuk.NewGathuk[Config]()
if err := gt.LoadConfigFiles(".env"); err != nil {
    log.Fatal(err)
}

config := gt.GetConfig()
port := config.Port    // Automatically converted to int
host := config.Host
debug := config.Debug  // Automatically converted to bool
```

**Benefits:**

- ✅ Automatic type conversion
- ✅ No manual parsing
- ✅ Type-safe struct
- ✅ Less boilerplate

### From encoding/json

**Before (encoding/json):**

```go
import "encoding/json"

file, err := os.Open("config.json")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

var config Config
decoder := json.NewDecoder(file)
if err := decoder.Decode(&config); err != nil {
    log.Fatal(err)
}
```

**After (Gathuk):**

```go
import "github.com/ahyalfan/gathuk"

gt := gathuk.NewGathuk[Config]()
if err := gt.LoadConfigFiles("config.json"); err != nil {
    log.Fatal(err)
}

config := gt.GetConfig()
```

**Benefits:**

- ✅ Environment variable support
- ✅ Multiple file merging
- ✅ Format-agnostic (same code for .env, JSON, etc.)
- ✅ Less boilerplate

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
