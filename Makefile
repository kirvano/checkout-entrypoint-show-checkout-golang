.PHONY: build clean deploy test fmt lint deps local-dev

# Variables
BINARY_NAME=bootstrap
LAMBDA_ZIP=lambda-function.zip
GO_VERSION=1.21

# Build the Lambda function for deployment
build:
	@echo "Building Lambda function..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) cmd/lambda/main.go
	@echo "Creating deployment package..."
	zip $(LAMBDA_ZIP) $(BINARY_NAME)
	@echo "Build complete: $(LAMBDA_ZIP)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME) $(LAMBDA_ZIP)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod verify

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Local development build (for testing)
local-dev:
	@echo "Building for local development..."
	go build -o checkout-dev cmd/lambda/main.go
	@echo "Local development binary: checkout-dev"

# Deploy to AWS Lambda (requires AWS CLI)
deploy: build
	@echo "Deploying to AWS Lambda..."
	@if [ -z "$(FUNCTION_NAME)" ]; then \
		echo "Error: FUNCTION_NAME environment variable is required"; \
		echo "Usage: make deploy FUNCTION_NAME=your-function-name"; \
		exit 1; \
	fi
	aws lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://$(LAMBDA_ZIP)
	@echo "Deployment complete"

# Create a new Lambda function (requires AWS CLI)
create-function:
	@echo "Creating new Lambda function..."
	@if [ -z "$(FUNCTION_NAME)" ] || [ -z "$(ROLE_ARN)" ]; then \
		echo "Error: FUNCTION_NAME and ROLE_ARN environment variables are required"; \
		echo "Usage: make create-function FUNCTION_NAME=your-function-name ROLE_ARN=your-role-arn"; \
		exit 1; \
	fi
	aws lambda create-function \
		--function-name $(FUNCTION_NAME) \
		--runtime provided.al2 \
		--role $(ROLE_ARN) \
		--handler bootstrap \
		--zip-file fileb://$(LAMBDA_ZIP) \
		--timeout 30 \
		--memory-size 512
	@echo "Function created successfully"

# Update function configuration
update-config:
	@echo "Updating function configuration..."
	@if [ -z "$(FUNCTION_NAME)" ]; then \
		echo "Error: FUNCTION_NAME environment variable is required"; \
		exit 1; \
	fi
	aws lambda update-function-configuration \
		--function-name $(FUNCTION_NAME) \
		--timeout 30 \
		--memory-size 512 \
		--environment Variables='{ENVIRONMENT=production,AWS_REGION=us-east-1}'
	@echo "Configuration updated"

# Run security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Generate project documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation server running at http://localhost:6060"

# Benchmark tests
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development tools installed"

# Docker build (optional)
docker-build:
	@echo "Building Docker image..."
	docker build -t checkout-go:latest .

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build Lambda function for deployment"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  deps           - Install dependencies"
	@echo "  tidy           - Tidy dependencies"
	@echo "  local-dev      - Build for local development"
	@echo "  deploy         - Deploy to AWS Lambda (requires FUNCTION_NAME)"
	@echo "  create-function- Create new Lambda function (requires FUNCTION_NAME and ROLE_ARN)"
	@echo "  update-config  - Update function configuration"
	@echo "  security       - Run security scan"
	@echo "  docs           - Generate documentation"
	@echo "  benchmark      - Run benchmark tests"
	@echo "  install-tools  - Install development tools"
	@echo "  docker-build   - Build Docker image"
	@echo "  help           - Show this help message"

# Default target
all: clean deps fmt lint test build