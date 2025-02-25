# EasyCfg

EasyCfg is a Go tool designed to simplify system configuration management. It automatically converts YAML configuration files into Go structs and uses Viper to read and monitor configuration changes.

## Features

- Automatically converts YAML configuration files to Go structs
- Generates corresponding Go files
- Uses Viper to read YAML configurations
- Supports hot reloading of configurations
- Supports monitoring configuration file changes

## Installation

```bash
go get github.com/chiayu0816/easycfg
```

## Usage

### Generate Configuration Structs

```bash
# Basic usage
go run cmd/easycfg/main.go -yaml path/to/config.yml

# Specify output directory
go run cmd/easycfg/main.go -yaml path/to/config.yml -output myconfig

# Specify package name
go run cmd/easycfg/main.go -yaml path/to/config.yml -package myconfig

# Monitor configuration file changes
go run cmd/easycfg/main.go -yaml path/to/config.yml -watch
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