#!/bin/bash

# Setup script for Workforce AI Transition Simulator

echo "Setting up Workforce AI Transition Simulator..."

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "Error: Go is not installed. Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Found Go version: $GO_VERSION"

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Run tests to verify setup
echo "Running tests..."
go test ./...

if [ $? -eq 0 ]; then
    echo "Setup complete! All tests passed."
    echo ""
    echo "To build the simulator, run: go build -o simulator ./cmd/simulator"
    echo "To run tests, run: go test ./..."
else
    echo "Setup complete, but some tests failed. Please review the output above."
fi
