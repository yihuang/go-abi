.PHONY: build test generate example clean

# Build the project
build:
	go build ./...

# Run tests
test:
	go test ./...

# Generate code from ABI
generate:
	@echo "Usage: go run cmd/generator/main.go -abi <abi_file> -out <output_file> -pkg <package_name>"

# Run the example
example:
	go run examples/simple/main.go

# Clean build artifacts
clean:
	go clean