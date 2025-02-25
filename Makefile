.PHONY: build run-example generate-config clean run-complete-example test test-verbose benchmark

# Build tool
build:
	go build -o bin/easycfg cmd/easycfg/main.go

# Run example
run-example:
	go run cmd/example/main.go

# Run complete example
run-complete-example:
	go run examples/complete/main.go examples/complete/config.go

# Generate configuration struct
generate-config:
	go run cmd/easycfg/main.go -yaml test_config.yml -output generated1

# Generate configuration struct and watch for changes
watch-config:
	go run cmd/easycfg/main.go -yaml test_config.yml -output generated -watch

# Run tests
test:
	go test -v ./...

# Run verbose tests
test-verbose:
	go test -v -cover ./...

# Run benchmark tests
benchmark:
	go test -bench=. -benchmem ./...

# Clean generated files
clean:
	rm -rf bin
	rm -rf generated

# Initialize directory structure
init:
	mkdir -p bin generated

# Install dependencies
deps:
	go mod tidy

# Default target
all: init deps build generate-config 