package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration
type Config struct {
	// Application Environment
	AppEnv string

	// Server Configuration
	Port     string
	GinMode  string
	LogLevel string

	// AWS Configuration
	AWSRegion                    string
	AWSDynamoDBAccessKeyID       string
	AWSDynamoDBSecretAccessKey   string
	AWSS3Bucket                  string
	AWSS3BasePath                string

	// Google Pay Configuration
	GooglePayMerchantIDD15 string
	GooglePayMerchantIDD2  string

	// Legacy Environment Variables (for backward compatibility)
	Environment string // maps to AppEnv
	S3Bucket    string // maps to AWSS3Bucket
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		// Application defaults
		AppEnv:   getEnvWithDefault("APP_ENV", getEnvWithDefault("ENVIRONMENT", "development")),
		Port:     getEnvWithDefault("PORT", "8080"),
		GinMode:  getEnvWithDefault("GIN_MODE", ""),
		LogLevel: getEnvWithDefault("LOG_LEVEL", "info"),

		// AWS defaults
		AWSRegion:                    getEnvWithDefault("AWS_REGION", "us-east-1"),
		AWSDynamoDBAccessKeyID:       os.Getenv("AWS_DYNAMODB_ACCESS_KEY_ID"),
		AWSDynamoDBSecretAccessKey:   os.Getenv("AWS_DYNAMODB_SECRET_ACCESS_KEY"),
		AWSS3Bucket:                  getEnvWithDefault("AWS_S3_BUCKET", getEnvWithDefault("S3_BUCKET", "")),
		AWSS3BasePath:                os.Getenv("S3_BASE_PATH"),

		// Google Pay defaults
		GooglePayMerchantIDD15: os.Getenv("GOOGLE_PAY_MERCHANT_ID_D15"),
		GooglePayMerchantIDD2:  os.Getenv("GOOGLE_PAY_MERCHANT_ID_D2"),

		// Legacy compatibility
		Environment: getEnvWithDefault("ENVIRONMENT", getEnvWithDefault("APP_ENV", "development")),
		S3Bucket:    getEnvWithDefault("S3_BUCKET", getEnvWithDefault("AWS_S3_BUCKET", "")),
	}

	// Validate required configuration
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// validate checks that required configuration values are present
func (c *Config) validate() error {
	var errors []string

	// Validate required AWS configuration if DynamoDB credentials are provided
	if c.AWSDynamoDBAccessKeyID != "" && c.AWSDynamoDBSecretAccessKey == "" {
		errors = append(errors, "AWS_DYNAMODB_SECRET_ACCESS_KEY is required when AWS_DYNAMODB_ACCESS_KEY_ID is provided")
	}
	if c.AWSDynamoDBSecretAccessKey != "" && c.AWSDynamoDBAccessKeyID == "" {
		errors = append(errors, "AWS_DYNAMODB_ACCESS_KEY_ID is required when AWS_DYNAMODB_SECRET_ACCESS_KEY is provided")
	}

	// Validate S3 bucket configuration
	if c.AWSS3Bucket == "" {
		errors = append(errors, "AWS_S3_BUCKET is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errors, ", "))
	}

	return nil
}

// IsProduction returns true if the application is running in production
func (c *Config) IsProduction() bool {
	env := strings.ToLower(c.AppEnv)
	return env == "production" || env == "prod"
}

// IsDevelopment returns true if the application is running in development
func (c *Config) IsDevelopment() bool {
	env := strings.ToLower(c.AppEnv)
	return env == "development" || env == "dev"
}

// IsTest returns true if the application is running in test mode
func (c *Config) IsTest() bool {
	env := strings.ToLower(c.AppEnv)
	return env == "test" || env == "testing"
}

// GetGinMode returns the appropriate Gin mode based on environment
func (c *Config) GetGinMode() string {
	if c.GinMode != "" {
		return c.GinMode
	}

	if c.IsProduction() {
		return "release"
	}
	if c.IsTest() {
		return "test"
	}
	return "debug"
}

// GetTableName returns the DynamoDB table name with environment prefix
func (c *Config) GetTableName(baseName string) string {
	return baseName
}

// HasDynamoDBCredentials returns true if explicit DynamoDB credentials are configured
func (c *Config) HasDynamoDBCredentials() bool {
	return c.AWSDynamoDBAccessKeyID != "" && c.AWSDynamoDBSecretAccessKey != ""
}

// GetS3BasePath returns the S3 base path, generating it if not explicitly set
func (c *Config) GetS3BasePath() string {
	if c.AWSS3BasePath != "" {
		return c.AWSS3BasePath
	}
	
	if c.AWSS3Bucket != "" {
		path := "https://s3.amazonaws.com/" + c.AWSS3Bucket + "/"
		return path
	}
	
	return ""
}

// getEnvWithDefault returns the environment variable value or the default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool returns the environment variable as a boolean or the default if not set
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvInt returns the environment variable as an integer or the default if not set
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
} 