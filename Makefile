.PHONY: build run clean test install lint help

BINARY_NAME=flappy-bird-tui
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE?=$(shell date -u +%Y-%m-%d)

# Build the game
build:
	go build -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(BINARY_NAME) .

# Run the game directly
run:
	go run .

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
	go clean

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run ./...

# Install to GOPATH/bin
install:
	go install

# Build for multiple platforms
release:
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o dist/$(BINARY_NAME)-windows-amd64.exe .

# Show help
help:
	@echo "Available targets:"
	@echo "  make        - Build the binary (default)"
	@echo "  make build  - Build the binary"
	@echo "  make run    - Run the game directly"
	@echo "  make test   - Run tests"
	@echo "  make lint   - Run linter (golangci-lint)"
	@echo "  make clean  - Remove build artifacts"
	@echo "  make install - Install to GOPATH/bin"
	@echo "  make release - Build for multiple platforms"
	@echo "  make help   - Show this help message"

# Default target
all: build
