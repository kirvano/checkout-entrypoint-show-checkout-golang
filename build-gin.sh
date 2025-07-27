#!/bin/bash

# Build script for Gin-based checkout server

echo "Building Gin-based checkout server..."

# Clean previous builds
rm -f checkout-gin
rm -f bin/checkout-gin*

# Build for current platform
echo "Building for current platform..."
go build -o checkout-gin ./cmd/gin/

# Build for Linux (useful for Docker/deployment)
echo "Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -o bin/checkout-gin-linux ./cmd/gin/

# Build for Windows
echo "Building for Windows amd64..."
GOOS=windows GOARCH=amd64 go build -o bin/checkout-gin.exe ./cmd/gin/

# Build for macOS
echo "Building for macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -o bin/checkout-gin-darwin ./cmd/gin/

echo "Build completed!"
echo "Files generated:"
echo "  - checkout-gin (current platform)"
echo "  - bin/checkout-gin-linux (Linux)"
echo "  - bin/checkout-gin.exe (Windows)"
echo "  - bin/checkout-gin-darwin (macOS)"
echo ""
echo "To run the server:"
echo "  ./checkout-gin"
echo ""
echo "Or with custom port:"
echo "  PORT=3000 ./checkout-gin" 