# Makefile for Flush Wars Backend

# Define variables
GO = go
BIN_NAME = flush-wars-backend
LINTER = revive

# Go Build target: build the Go binary
build:
	$(GO) build -o $(BIN_NAME)

# Test target: run Go tests
test:
	$(GO) test -v ./...

# Lint target: run the linter
lint:
	$(LINTER) ./...

# Clean target: remove the compiled binary
clean:
	rm -f $(BIN_NAME)

# Run target: build and run the Go application
run: build
	./$(BIN_NAME)

# Install dependencies (optional, depending on how you manage dependencies)
deps:
	$(GO) get -v ./...

# Format the code using gofmt
fmt:
	$(GO) fmt ./...

# Go mod tidy target: tidy up go.mod and go.sum
mod-tidy:
	$(GO) mod tidy

# All target: build, lint, and test
all: deps build fmt mod-tidy lint test 
