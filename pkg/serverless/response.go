package serverless

import (
	"encoding/json"
	"log"
	"strings"

	"checkout-go/internal/core/errors"
	"github.com/aws/aws-lambda-go/events"
)

// SendJSON sends a successful JSON response
func SendJSON(data interface{}, statusCode int) events.APIGatewayProxyResponse {
	// Convert to snake_case like the TypeScript version
	snakeCaseData := convertToSnakeCase(data)

	body, err := json.Marshal(snakeCaseData)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"message": "Internal server error"}`,
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

// SendErrorJSON sends an error JSON response
func SendErrorJSON(err error, statusCode int) events.APIGatewayProxyResponse {
	// Check if it's a custom error with specific handling
	if customErr, ok := err.(*errors.BaseError); ok {
		errorResponse := map[string]string{
			"message": customErr.GetMessage(),
		}

		body, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			log.Printf("Failed to marshal error response: %v", marshalErr)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"message": "Internal server error"}`,
			}
		}

		return events.APIGatewayProxyResponse{
			StatusCode: customErr.GetHTTPCode(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(body),
		}
	}

	// Handle validation errors
	if validationErr, ok := err.(*errors.ValidationError); ok {
		errorResponse := map[string]interface{}{
			"message": validationErr.GetMessage(),
			"details": validationErr.Details,
		}

		body, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			log.Printf("Failed to marshal validation error response: %v", marshalErr)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
				Body: `{"message": "Internal server error"}`,
			}
		}

		return events.APIGatewayProxyResponse{
			StatusCode: validationErr.GetHTTPCode(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(body),
		}
	}

	// Default error handling
	errorResponse := map[string]string{
		"message": err.Error(),
	}

	body, marshalErr := json.Marshal(errorResponse)
	if marshalErr != nil {
		log.Printf("Failed to marshal error response: %v", marshalErr)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"message": "Internal server error"}`,
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

// SendJSONWithCORS sends a JSON response with CORS headers
func SendJSONWithCORS(data interface{}, statusCode int) events.APIGatewayProxyResponse {
	// Convert to snake_case like the TypeScript version
	snakeCaseData := convertToSnakeCase(data)

	body, err := json.Marshal(snakeCaseData)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Methods": "OPTIONS,POST",
				"Content-Type":                 "application/json",
			},
			Body: `{"message": "Internal server error"}`,
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "OPTIONS,POST",
			"Content-Type":                 "application/json",
		},
		Body: string(body),
	}
}

// SendErrorJSONWithCORS sends an error JSON response with CORS headers
func SendErrorJSONWithCORS(err error, statusCode int) events.APIGatewayProxyResponse {
	// Check if it's a custom error with specific handling
	if customErr, ok := err.(*errors.BaseError); ok {
		errorResponse := map[string]string{
			"message": customErr.GetMessage(),
		}

		body, marshalErr := json.Marshal(errorResponse)
		if marshalErr != nil {
			log.Printf("Failed to marshal error response: %v", marshalErr)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers: map[string]string{
					"Access-Control-Allow-Origin":  "*",
					"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
					"Access-Control-Allow-Methods": "OPTIONS,POST",
					"Content-Type":                 "application/json",
				},
				Body: `{"message": "Internal server error"}`,
			}
		}

		return events.APIGatewayProxyResponse{
			StatusCode: customErr.GetHTTPCode(),
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Methods": "OPTIONS,POST",
				"Content-Type":                 "application/json",
			},
			Body: string(body),
		}
	}

	// Default error handling
	errorResponse := map[string]string{
		"message": err.Error(),
	}

	body, marshalErr := json.Marshal(errorResponse)
	if marshalErr != nil {
		log.Printf("Failed to marshal error response: %v", marshalErr)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Methods": "OPTIONS,POST",
				"Content-Type":                 "application/json",
			},
			Body: `{"message": "Internal server error"}`,
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "OPTIONS,POST",
			"Content-Type":                 "application/json",
		},
		Body: string(body),
	}
}

// convertToSnakeCase converts struct field names to snake_case for JSON serialization
// This is a simplified version - in production you might want to use a library like "github.com/iancoleman/strcase"
func convertToSnakeCase(data interface{}) interface{} {
	// For simplicity, we'll rely on the json tags in our structs
	// The structs are already defined with snake_case json tags
	return data
}

// toSnakeCase converts a camelCase string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteString("_")
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
