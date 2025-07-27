#!/bin/bash
set -e

echo "Building local development server..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build the local server
go build -o bin/checkout-local.exe ./cmd/local

echo "✅ Build successful!"
echo ""
echo "To start the server, run:"
echo "   ./bin/checkout-local.exe"
echo ""
echo "Server will be available at: http://localhost:8080"
echo "Test endpoint: http://localhost:8080/checkout/123e4567-e89b-12d3-a456-426614174000"