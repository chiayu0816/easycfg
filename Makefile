.PHONY: build build-cli run-example generate-config clean run-complete-example test test-verbose benchmark install

# Build tool
build:
	go build -o bin/easycfg cmd/easycfg/main.go

# Build CLI tool
build-cli:
	go build -o bin/easycfgcli cmd/easycfgcli/main.go

# Install CLI tool
install:
	go install ./cmd/easycfgcli

# Run example
run-example:
	go run cmd/example/main.go

# Run complete example
run-complete-example:
	go run examples/complete/main.go examples/complete/config.go

# Generate configuration struct
generate-config:
	go run cmd/easycfg/main.go -yaml test_config.yml -output generated

# Generate configuration struct using CLI tool
generate-config-cli:
	./bin/easycfgcli -yaml test_config.yml -output generated

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
	rm -rf test_output

# Initialize directory structure
init:
	mkdir -p bin generated

# Install dependencies
deps:
	go mod tidy

# Default target
all: init deps build build-cli generate-config 

