package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Client wraps the AWS DynamoDB client
type Client struct {
	dynamoDB *dynamodb.Client
	config   aws.Config
}

// NewClient creates a new DynamoDB client wrapper
func NewClient(cfg aws.Config) (*Client, error) {
	client := dynamodb.NewFromConfig(cfg)

	return &Client{
		dynamoDB: client,
		config:   cfg,
	}, nil
}

// GetDynamoDB returns the underlying DynamoDB client
func (c *Client) GetDynamoDB() *dynamodb.Client {
	return c.dynamoDB
}

// GetConfig returns the AWS configuration
func (c *Client) GetConfig() aws.Config {
	return c.config
}
