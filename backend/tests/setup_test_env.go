package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Define the project root directory
	projectRoot := "/home/ubuntu/trading_platform_implementation"
	
	// Create test directories if they don't exist
	testDirs := []string{
		filepath.Join(projectRoot, "backend", "tests"),
		filepath.Join(projectRoot, "backend", "tests", "mocks"),
		filepath.Join(projectRoot, "backend", "tests", "fixtures"),
	}
	
	for _, dir := range testDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
	
	// Install test dependencies
	cmd := exec.Command("go", "get", "-u", "github.com/stretchr/testify/assert", "github.com/stretchr/testify/require", "github.com/stretchr/testify/mock")
	cmd.Dir = filepath.Join(projectRoot, "backend")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	fmt.Println("Installing test dependencies...")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error installing test dependencies: %v\n", err)
		os.Exit(1)
	}
	
	// Create a test configuration file
	testConfig := `{
		"database": {
			"host": "localhost",
			"port": 5432,
			"user": "test_user",
			"password": "test_password",
			"dbname": "test_db"
		},
		"redis": {
			"host": "localhost",
			"port": 6379
		},
		"rabbitmq": {
			"host": "localhost",
			"port": 5672,
			"user": "guest",
			"password": "guest"
		},
		"jwt": {
			"secret": "test-secret-key",
			"expiration": 86400
		},
		"broker": {
			"xts": {
				"endpoint": "https://xts-api.example.com",
				"api_key": "test-api-key",
				"api_secret": "test-api-secret"
			}
		}
	}`
	
	configPath := filepath.Join(projectRoot, "backend", "tests", "config.json")
	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		fmt.Printf("Error writing test configuration: %v\n", err)
		os.Exit(1)
	}
	
	// Create a test runner script
	testRunner := `#!/bin/bash
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
`
	
	runnerPath := filepath.Join(projectRoot, "backend", "tests", "run_tests.sh")
	if err := os.WriteFile(runnerPath, []byte(testRunner), 0755); err != nil {
		fmt.Printf("Error writing test runner script: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Test environment setup completed successfully!")
	fmt.Println("To run tests, execute: ./backend/tests/run_tests.sh")
}
