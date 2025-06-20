#!/bin/bash

# Install swag if not already installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate swagger documentation
echo "Generating Swagger documentation..."
swag init -g main.go -o docs

echo "Swagger documentation generated successfully!"
echo "You can now access the Swagger UI at: http://localhost:8080/swagger/index.html" 