package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	appConfig "checkout-go/internal/config"
)

// NewConfig creates a new AWS configuration using the provided app configuration
func NewConfig(cfg *appConfig.Config) (aws.Config, error) {
	var configOptions []func(*config.LoadOptions) error

	// Set the region
	configOptions = append(configOptions, config.WithRegion(cfg.AWSRegion))

	// Use explicit credentials if provided for DynamoDB
	if cfg.HasDynamoDBCredentials() {
		creds := credentials.NewStaticCredentialsProvider(
			cfg.AWSDynamoDBAccessKeyID,
			cfg.AWSDynamoDBSecretAccessKey,
			"", // token is optional for static credentials
		)
		configOptions = append(configOptions, config.WithCredentialsProvider(creds))
	}

	// Load AWS configuration
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return aws.Config{}, err
	}

	return awsConfig, nil
}

// GetTableName returns the DynamoDB table name with environment prefix
// This function is kept for backward compatibility but now uses the config
func GetTableName(cfg *appConfig.Config, baseName string) string {
	return cfg.GetTableName(baseName)
}

// IsProduction checks if the current environment is production
// This function is kept for backward compatibility but now uses the config
func IsProduction(cfg *appConfig.Config) bool {
	return cfg.IsProduction()
}
