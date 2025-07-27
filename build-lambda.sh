#!/bin/bash
set -e

echo "Building AWS Lambda deployment package..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Build for Lambda (Linux)
echo "Building Go binary for Linux..."
GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bin/bootstrap ./cmd/lambda

# Create deployment package
echo "Creating deployment package..."
cd bin
zip -r checkout-lambda.zip bootstrap
cd ..

echo "✅ Build successful!"
echo ""
echo "Deployment package created: bin/checkout-lambda.zip"
echo ""
echo "To deploy to AWS Lambda:"
echo "1. Upload bin/checkout-lambda.zip to your Lambda function"
echo "2. Set handler to: bootstrap"
echo "3. Set runtime to: provided.al2 or provided.al2023"
echo "4. Configure environment variables as needed"