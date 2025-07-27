package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"

	"checkout-go/internal/infrastructure/di"
	"checkout-go/internal/usecases/showcheckout"
	"checkout-go/pkg/serverless"
)

var validate *validator.Validate
var container *di.Container

func init() {
	validate = validator.New()
	log.Println("Initializing Lambda function...")
	
	// Initialize dependency injection container (which loads configuration)
	var err error
	container, err = di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}
	
	// Get configuration for logging
	config := container.GetConfig()
	log.Printf("Lambda initialized successfully - Environment: %s", config.AppEnv)
}

// handleCheckoutEntrypoint is the main Lambda handler function
func handleCheckoutEntrypoint(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request for path: %s", event.Path)

	// Initialize DI container
	// container, err := di.NewContainer() // This line is removed as per the new_code
	// if err != nil {
	// 	log.Printf("Failed to initialize DI container: %v", err)
	// 	return serverless.SendErrorJSON(fmt.Errorf("internal server error"), 500), nil
	// }

	// Get use case from container
	useCase := container.GetShowCheckoutUseCase()

	// Extract offer UUID from path parameters or proxy path
	var offerUUID string
	var exists bool

	// First try to get from direct path parameter
	offerUUID, exists = event.PathParameters["offerUuid"]
	
	// If not found, try to extract from proxy path or main path
	if !exists || offerUUID == "" {
		// Try proxy path parameter
		if proxyPath, proxyExists := event.PathParameters["proxy"]; proxyExists {
			offerUUID = extractUUIDFromPath(proxyPath)
		}
		
		// If still not found, try the main path
		if offerUUID == "" {
			offerUUID = extractUUIDFromPath(event.Path)
		}
	}

	if offerUUID == "" {
		log.Printf("Missing offerUuid in path: %s, pathParams: %v", event.Path, event.PathParameters)
		return serverless.SendErrorJSON(fmt.Errorf("missing offerUuid parameter"), 400), nil
	}

	// Parse query parameters and build request
	req, err := buildShowCheckoutRequest(offerUUID, event.QueryStringParameters, event.Headers)
	if err != nil {
		log.Printf("Failed to build request: %v", err)
		return serverless.SendErrorJSON(err, 400), nil
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		log.Printf("Request validation failed: %v", err)
		return serverless.SendErrorJSON(fmt.Errorf("validation failed: %v", err), 400), nil
	}

	// Execute use case
	result, err := useCase.Execute(ctx, req)
	if err != nil {
		log.Printf("Use case execution failed: %v", err)
		return serverless.SendErrorJSON(err, 400), nil
	}

	log.Printf("Successfully processed checkout request for offer: %s", offerUUID)
	return serverless.SendJSON(result, 200), nil
}

// buildShowCheckoutRequest constructs the request from Lambda event data
func buildShowCheckoutRequest(offerUUID string, queryParams map[string]string, headers map[string]string) (*showcheckout.ShowCheckoutRequest, error) {
	if queryParams == nil {
		queryParams = make(map[string]string)
	}

	// Extract client info from query parameters
	clientInfo := showcheckout.ClientInfo{
		IsMobile: queryParams["isMobile"] == "true",
	}

	// Map optional string parameters
	if ip := queryParams["ip"]; ip != "" {
		clientInfo.IP = &ip
	}
	if userAgent := queryParams["userAgent"]; userAgent != "" {
		clientInfo.UserAgent = &userAgent
	}
	if browser := queryParams["browser"]; browser != "" {
		clientInfo.Browser = &browser
	}
	if browserVersion := queryParams["browserVersion"]; browserVersion != "" {
		clientInfo.BrowserVersion = &browserVersion
	}
	if os := queryParams["os"]; os != "" {
		clientInfo.OS = &os
	}
	if osVersion := queryParams["osVersion"]; osVersion != "" {
		clientInfo.OSVersion = &osVersion
	}
	if country := queryParams["country"]; country != "" {
		clientInfo.Country = &country
	}
	if state := queryParams["state"]; state != "" {
		clientInfo.State = &state
	}
	if city := queryParams["city"]; city != "" {
		clientInfo.City = &city
	}
	if lat := queryParams["lat"]; lat != "" {
		clientInfo.Lat = &lat
	}
	if lon := queryParams["lon"]; lon != "" {
		clientInfo.Lon = &lon
	}

	// Extract UTM info
	utmInfo := showcheckout.UTMInfo{}
	if src := queryParams["src"]; src != "" {
		utmInfo.Src = &src
	}
	if utmSource := queryParams["utm_source"]; utmSource != "" {
		utmInfo.UTMSource = &utmSource
	}
	if utmMedium := queryParams["utm_medium"]; utmMedium != "" {
		utmInfo.UTMMedium = &utmMedium
	}
	if utmCampaign := queryParams["utm_campaign"]; utmCampaign != "" {
		utmInfo.UTMCampaign = &utmCampaign
	}
	if utmTerm := queryParams["utm_term"]; utmTerm != "" {
		utmInfo.UTMTerm = &utmTerm
	}
	if utmContent := queryParams["utm_content"]; utmContent != "" {
		utmInfo.UTMContent = &utmContent
	}

	// Build the request
	req := &showcheckout.ShowCheckoutRequest{
		OfferUUID:  offerUUID,
		ClientInfo: clientInfo,
		UTMInfo:    utmInfo,
	}

	// Extract optional tracking parameters
	if aff := queryParams["aff"]; aff != "" {
		req.Aff = &aff
	}
	if fbclid := queryParams["fbclid"]; fbclid != "" {
		req.Fbclid = &fbclid
	}
	if gclid := queryParams["gclid"]; gclid != "" {
		req.Gclid = &gclid
	}
	if ttclid := queryParams["ttclid"]; ttclid != "" {
		req.Ttclid = &ttclid
	}
	if clickId := queryParams["clickId"]; clickId != "" {
		req.ClickID = &clickId
	}
	if originalUrl := queryParams["originalUrl"]; originalUrl != "" {
		// Validate URL
		if _, err := url.Parse(originalUrl); err == nil {
			req.OriginalURL = &originalUrl
		}
	}

	// Extract cookie from headers
	if cookie := headers["cookie"]; cookie != "" {
		req.Cookie = &cookie
	}
	// Also check lowercase (some proxies might lowercase headers)
	if cookie := headers["Cookie"]; cookie != "" {
		req.Cookie = &cookie
	}

	return req, nil
}

// parseIntParam safely parses an integer parameter
func parseIntParam(value string) (*int, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// parseFloatParam safely parses a float parameter
func parseFloatParam(value string) (*float64, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// extractUUIDFromPath extracts UUID from a path string
func extractUUIDFromPath(path string) string {
	// Split path by / and look for UUID pattern
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	for _, part := range parts {
		// UUID pattern: 8-4-4-4-12 characters (with hyphens) or 32 characters (without hyphens)
		if len(part) == 36 && strings.Count(part, "-") == 4 {
			// Standard UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
			if isValidUUIDFormat(part) {
				return part
			}
		} else if len(part) == 32 {
			// UUID without hyphens: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
			if isValidHexString(part) {
				// Convert to standard UUID format
				return part[:8] + "-" + part[8:12] + "-" + part[12:16] + "-" + part[16:20] + "-" + part[20:]
			}
		}
	}
	
	return ""
}

// isValidUUIDFormat checks if string matches UUID format
func isValidUUIDFormat(s string) bool {
	if len(s) != 36 {
		return false
	}
	
	for i, r := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if r != '-' {
				return false
			}
		} else {
			if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
				return false
			}
		}
	}
	
	return true
}

// isValidHexString checks if string contains only hex characters
func isValidHexString(s string) bool {
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

func main() {
	log.Println("Starting checkout Lambda function")
	lambda.Start(handleCheckoutEntrypoint)
}
