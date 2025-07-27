package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"

	"checkout-go/internal/infrastructure/di"
	"checkout-go/internal/usecases/showcheckout"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func main() {
	log.Println("Starting local checkout server...")
	
	// Initialize dependency injection container (which loads configuration)
	log.Println("Initializing dependency injection container...")
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}
	log.Println("DI container initialized successfully")
	
	// Get configuration from container
	config := container.GetConfig()
	
	// Test that we can get the use case
	useCase := container.GetShowCheckoutUseCase()
	if useCase == nil {
		log.Fatalf("Failed to get ShowCheckoutUseCase from container")
	}
	log.Println("ShowCheckoutUseCase retrieved successfully")
	
	http.HandleFunc("/checkout/", handleCheckout)
	http.HandleFunc("/health", handleHealth)
	
	log.Println("Server endpoints registered:")
	log.Println("  - GET /health")
	log.Println("  - GET /checkout/{uuid}")
	log.Println("")
	log.Printf("Environment: %s", config.AppEnv)
	log.Printf("Server running at http://localhost:%s", config.Port)
	log.Printf("Example: http://localhost:%s/checkout/123e4567-e89b-12d3-a456-426614174000", config.Port)
	log.Println("")
	log.Println("Press Ctrl+C to stop the server")
	
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func handleCheckout(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract offer UUID from path
	path := strings.TrimPrefix(r.URL.Path, "/checkout/")
	if path == "" {
		sendError(w, "Missing offer UUID", http.StatusBadRequest)
		return
	}

	log.Printf("Processing request for offer: %s", path)

	// Initialize DI container
	container, err := di.NewContainer()
	if err != nil {
		log.Printf("Failed to initialize DI container: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get use case from container
	useCase := container.GetShowCheckoutUseCase()

	// Build request from HTTP request
	req, err := buildRequestFromHTTP(path, r)
	if err != nil {
		log.Printf("Failed to build request: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		log.Printf("Request validation failed: %v", err)
		sendError(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}

	// Execute use case
	result, err := useCase.Execute(context.Background(), req)
	if err != nil {
		log.Printf("Use case execution failed: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}

	log.Printf("Successfully processed checkout request for offer: %s", path)
}

func buildRequestFromHTTP(offerUUID string, r *http.Request) (*showcheckout.ShowCheckoutRequest, error) {
	queryParams := r.URL.Query()

	// Extract client info from query parameters
	clientInfo := showcheckout.ClientInfo{
		IsMobile: queryParams.Get("isMobile") == "true",
	}

	// Map optional string parameters
	if ip := queryParams.Get("ip"); ip != "" {
		clientInfo.IP = &ip
	}
	if userAgent := r.Header.Get("User-Agent"); userAgent != "" {
		clientInfo.UserAgent = &userAgent
	}
	if browser := queryParams.Get("browser"); browser != "" {
		clientInfo.Browser = &browser
	}
	if browserVersion := queryParams.Get("browserVersion"); browserVersion != "" {
		clientInfo.BrowserVersion = &browserVersion
	}
	if os := queryParams.Get("os"); os != "" {
		clientInfo.OS = &os
	}
	if osVersion := queryParams.Get("osVersion"); osVersion != "" {
		clientInfo.OSVersion = &osVersion
	}
	if country := queryParams.Get("country"); country != "" {
		clientInfo.Country = &country
	}
	if state := queryParams.Get("state"); state != "" {
		clientInfo.State = &state
	}
	if city := queryParams.Get("city"); city != "" {
		clientInfo.City = &city
	}
	if lat := queryParams.Get("lat"); lat != "" {
		clientInfo.Lat = &lat
	}
	if lon := queryParams.Get("lon"); lon != "" {
		clientInfo.Lon = &lon
	}

	// Extract UTM info
	utmInfo := showcheckout.UTMInfo{}
	if src := queryParams.Get("src"); src != "" {
		utmInfo.Src = &src
	}
	if utmSource := queryParams.Get("utm_source"); utmSource != "" {
		utmInfo.UTMSource = &utmSource
	}
	if utmMedium := queryParams.Get("utm_medium"); utmMedium != "" {
		utmInfo.UTMMedium = &utmMedium
	}
	if utmCampaign := queryParams.Get("utm_campaign"); utmCampaign != "" {
		utmInfo.UTMCampaign = &utmCampaign
	}
	if utmTerm := queryParams.Get("utm_term"); utmTerm != "" {
		utmInfo.UTMTerm = &utmTerm
	}
	if utmContent := queryParams.Get("utm_content"); utmContent != "" {
		utmInfo.UTMContent = &utmContent
	}

	// Build the request
	req := &showcheckout.ShowCheckoutRequest{
		OfferUUID:  offerUUID,
		ClientInfo: clientInfo,
		UTMInfo:    utmInfo,
	}

	// Extract optional tracking parameters
	if aff := queryParams.Get("aff"); aff != "" {
		req.Aff = &aff
	}
	if fbclid := queryParams.Get("fbclid"); fbclid != "" {
		req.Fbclid = &fbclid
	}
	if gclid := queryParams.Get("gclid"); gclid != "" {
		req.Gclid = &gclid
	}
	if ttclid := queryParams.Get("ttclid"); ttclid != "" {
		req.Ttclid = &ttclid
	}
	if clickId := queryParams.Get("clickId"); clickId != "" {
		req.ClickID = &clickId
	}
	if originalUrl := queryParams.Get("originalUrl"); originalUrl != "" {
		// Validate URL
		if _, err := url.Parse(originalUrl); err == nil {
			req.OriginalURL = &originalUrl
		}
	}

	// Extract cookie from headers
	if cookie := r.Header.Get("Cookie"); cookie != "" {
		req.Cookie = &cookie
	}

	return req, nil
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   true,
		"message": message,
		"status":  status,
	})
}