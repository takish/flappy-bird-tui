.PHONY: build run clean test install

# Build the game
build:
	go build -o flappy-bird .

# Run the game directly
run:
	go run .

# Clean build artifacts
clean:
	rm -f flappy-bird
	go clean

# Run tests
test:
	go test -v ./...

# Install to GOPATH/bin
install:
	go install

# Build for multiple platforms
release:
	GOOS=darwin GOARCH=arm64 go build -o dist/flappy-bird-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o dist/flappy-bird-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o dist/flappy-bird-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -o dist/flappy-bird-windows-amd64.exe .

# Default target
all: build
