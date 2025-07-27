# Environment Variables

This document lists all supported environment variables for the checkout-go application.

## Application Configuration

### `APP_ENV`
- **Description**: The environment the application is running in
- **Values**: `development`, `production`, `test`
- **Default**: `development`
- **Example**: `APP_ENV=production`

## Server Configuration

### `PORT`
- **Description**: Port for the HTTP server
- **Default**: `8080`
- **Example**: `PORT=8080`

### `GIN_MODE`
- **Description**: Gin web framework mode
- **Values**: `debug`, `release`, `test`
- **Default**: Auto-detected from APP_ENV
- **Example**: `GIN_MODE=release`

### `LOG_LEVEL`
- **Description**: Log level for the application
- **Values**: `debug`, `info`, `warn`, `error`
- **Default**: `info`
- **Example**: `LOG_LEVEL=info`

## AWS Configuration

### `AWS_REGION`
- **Description**: AWS region for all AWS services
- **Default**: `us-east-1`
- **Example**: `AWS_REGION=us-east-1`

### `AWS_DYNAMODB_ACCESS_KEY_ID`
- **Description**: AWS Access Key ID specifically for DynamoDB operations
- **Required**: Optional (uses AWS default credential chain if not provided)
- **Example**: `AWS_DYNAMODB_ACCESS_KEY_ID=AKIAY3`

### `AWS_DYNAMODB_SECRET_ACCESS_KEY`
- **Description**: AWS Secret Access Key specifically for DynamoDB operations
- **Required**: Required if `AWS_DYNAMODB_ACCESS_KEY_ID` is provided
- **Example**: `AWS_DYNAMODB_SECRET_ACCESS_KEY=Th6swIw`

### `AWS_S3_BUCKET`
- **Description**: S3 bucket name for file storage
- **Required**: Yes
- **Example**: `AWS_S3_BUCKET=production.kirvano.com`

### `S3_BASE_PATH`
- **Description**: Custom S3 base path for file URLs
- **Default**: Auto-generated from bucket if not set
- **Example**: `S3_BASE_PATH=https://s3.amazonaws.com/production.kirvano.com/`

## Google Pay Configuration

### `GOOGLE_PAY_MERCHANT_ID_D15`
- **Description**: Google Pay Merchant ID for D15 transactions
- **Example**: `GOOGLE_PAY_MERCHANT_ID_D15=d3d2a`

### `GOOGLE_PAY_MERCHANT_ID_D2`
- **Description**: Google Pay Merchant ID for D2 transactions
- **Example**: `GOOGLE_PAY_MERCHANT_ID_D2=f20cd38`

## Legacy Environment Variables

These variables are maintained for backward compatibility and are automatically mapped to their new equivalents:

### `ENVIRONMENT`
- **Maps to**: `APP_ENV`
- **Description**: Legacy environment variable (use `APP_ENV` instead)

### `S3_BUCKET`
- **Maps to**: `AWS_S3_BUCKET`
- **Description**: Legacy S3 bucket variable (use `AWS_S3_BUCKET` instead)

## Example Configuration

```bash
# Production environment configuration
APP_ENV=production
PORT=8080
AWS_REGION=us-east-1
AWS_DYNAMODB_ACCESS_KEY_ID=AKIAY3
AWS_DYNAMODB_SECRET_ACCESS_KEY=Th6swIw
AWS_S3_BUCKET=production.kirvano.com
GOOGLE_PAY_MERCHANT_ID_D15=d3d2a
GOOGLE_PAY_MERCHANT_ID_D2=f20cd38
```

## Configuration Validation

The application performs validation on startup and will fail to start if:

1. `AWS_DYNAMODB_ACCESS_KEY_ID` is provided without `AWS_DYNAMODB_SECRET_ACCESS_KEY`
2. `AWS_DYNAMODB_SECRET_ACCESS_KEY` is provided without `AWS_DYNAMODB_ACCESS_KEY_ID`
3. `AWS_S3_BUCKET` is not provided

## Accessing Configuration in Code

Configuration is centrally managed through the `internal/config` package and is available throughout the application via the dependency injection container:

```go
// Get configuration from DI container
config := container.GetConfig()

// Use configuration values
bucket := config.AWSS3Bucket
isProduction := config.IsProduction()
tableName := config.GetTableName("users")
``` 