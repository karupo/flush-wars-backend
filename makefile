# Makefile

# Define the binary name
BIN_NAME=flush-wars-backend

# Go build target
build:
	go build -o $(BIN_NAME)

# Run revive linter
lint:
	revive -config revive.toml ./...

# Clean up generated files
clean:
	rm -f $(BIN_NAME)

# Default target to build and lint at the same time
all: build lint
