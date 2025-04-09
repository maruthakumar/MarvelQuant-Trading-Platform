#!/bin/bash
set -e

# Change to the backend directory
cd "$(dirname "$0")/.."

# Run all tests with verbose output
go test -v ./tests/...

# Run tests with coverage
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html

echo "Tests completed successfully!"
echo "Coverage report generated at: coverage.html"
