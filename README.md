# EasyCfg

[![GitHub release](https://img.shields.io/github/v/release/chiayu0816/easycfg)](https://github.com/chiayu0816/easycfg/releases/latest)

EasyCfg is a Go tool designed to simplify system configuration management. It automatically converts YAML configuration files into Go structs and uses Viper to read and monitor configuration changes.

## Features

- Automatically converts YAML configuration files to Go structs
- Generates corresponding Go files
- Uses Viper to read YAML configurations
- Supports hot reloading of configurations
- Supports monitoring configuration file changes

## Installation

```bash
# Latest version
go get github.com/chiayu0816/easycfg

# Specific version
go get github.com/chiayu0816/easycfg@v1.0.0

# Install CLI tool
go install github.com/chiayu0816/easycfg/cmd/easycfgcli@latest
```

## Usage

### Generate Configuration Structs

```bash
# Using go run (if you've installed the package with go get)
go run github.com/chiayu0816/easycfg/cmd/easycfg -yaml path/to/config.yml

# Using the installed CLI tool
easycfgcli -yaml path/to/config.yml

# Specify output directory
easycfgcli -yaml path/to/config.yml -output myconfig

# Specify package name
easycfgcli -yaml path/to/config.yml -package myconfig

# Monitor configuration file changes
easycfgcli -yaml path/to/config.yml -watch
```

### Using Generated Configurations in Your Program

```go
package main

import (
    "fmt"
    "log"

    "github.com/chiayu0816/easycfg"
)

func main() {
    // Create a configuration struct instance
    cfg := &MyConfig{}

    // Load configuration
    if err := easycfg.LoadConfig("config.yml", cfg); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Use configuration
    fmt.Printf("Configuration value: %s\n", cfg.SomeField)

    // Monitor configuration changes
    easycfg.WatchConfig("config.yml", cfg, func() {
        fmt.Println("Configuration has been updated")
    })
}
```

## Examples

Check the `examples/complete` directory for a complete example.

Run examples:

```bash
# Run basic example
make run-example

# Run complete example
make run-complete-example
```

## Development

If you're working on the EasyCfg codebase, you can use the Makefile to simplify common tasks:

```bash
# Build the main tool
make build

# Build the CLI tool
make build-cli

# Install the CLI tool locally
make install

# Generate configuration from test_config.yml
make generate-config

# Generate configuration using the CLI tool
make generate-config-cli

# Watch for changes in the configuration file
make watch-config

# Clean generated files
make clean
```

## Testing

EasyCfg includes unit tests and benchmarks to ensure code correctness and performance.

Run tests:

```bash
# Run all tests
make test

# Run verbose tests (with coverage)
make test-verbose

# Run benchmarks
make benchmark
```

## Dependencies

- [github.com/go-yaml/yaml](https://github.com/go-yaml/yaml)
- [github.com/spf13/viper](https://github.com/spf13/viper)

## License

MIT 