package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"checkout-go/internal/infrastructure/di"
	"checkout-go/internal/usecases/showcheckout"
)

var validate = validator.New()

// CheckoutHandlers contains the HTTP handlers for checkout endpoints
type CheckoutHandlers struct {
	container *di.Container
}

// NewCheckoutHandlers creates a new CheckoutHandlers instance
func NewCheckoutHandlers(container *di.Container) *CheckoutHandlers {
	return &CheckoutHandlers{
		container: container,
	}
}

// HealthCheck handles GET /health
func (h *CheckoutHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "checkout-api",
	})
}

// ShowCheckout handles GET /checkout/:uuid
func (h *CheckoutHandlers) ShowCheckout(c *gin.Context) {
	offerUUID := c.Param("uuid")
	if offerUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Missing offer UUID",
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Build request from HTTP context
	req, err := h.buildRequestFromGin(offerUUID, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Validation failed: " + err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Get use case from container
	useCase := h.container.GetShowCheckoutUseCase()

	// Execute use case
	result, err := useCase.Execute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	// Send successful response
	c.JSON(http.StatusOK, result)
}

// buildRequestFromGin constructs the ShowCheckoutRequest from Gin context
func (h *CheckoutHandlers) buildRequestFromGin(offerUUID string, c *gin.Context) (*showcheckout.ShowCheckoutRequest, error) {
	// Extract client info from query parameters and headers
	clientInfo := showcheckout.ClientInfo{
		IsMobile: c.Query("isMobile") == "true",
	}

	// Map optional string parameters
	if ip := c.Query("ip"); ip != "" {
		clientInfo.IP = &ip
	}
	if userAgent := c.GetHeader("User-Agent"); userAgent != "" {
		clientInfo.UserAgent = &userAgent
	}
	if browser := c.Query("browser"); browser != "" {
		clientInfo.Browser = &browser
	}
	if browserVersion := c.Query("browserVersion"); browserVersion != "" {
		clientInfo.BrowserVersion = &browserVersion
	}
	if os := c.Query("os"); os != "" {
		clientInfo.OS = &os
	}
	if osVersion := c.Query("osVersion"); osVersion != "" {
		clientInfo.OSVersion = &osVersion
	}
	if country := c.Query("country"); country != "" {
		clientInfo.Country = &country
	}
	if state := c.Query("state"); state != "" {
		clientInfo.State = &state
	}
	if city := c.Query("city"); city != "" {
		clientInfo.City = &city
	}
	if lat := c.Query("lat"); lat != "" {
		clientInfo.Lat = &lat
	}
	if lon := c.Query("lon"); lon != "" {
		clientInfo.Lon = &lon
	}

	// Extract UTM info
	utmInfo := showcheckout.UTMInfo{}
	if src := c.Query("src"); src != "" {
		utmInfo.Src = &src
	}
	if utmSource := c.Query("utm_source"); utmSource != "" {
		utmInfo.UTMSource = &utmSource
	}
	if utmMedium := c.Query("utm_medium"); utmMedium != "" {
		utmInfo.UTMMedium = &utmMedium
	}
	if utmCampaign := c.Query("utm_campaign"); utmCampaign != "" {
		utmInfo.UTMCampaign = &utmCampaign
	}
	if utmTerm := c.Query("utm_term"); utmTerm != "" {
		utmInfo.UTMTerm = &utmTerm
	}
	if utmContent := c.Query("utm_content"); utmContent != "" {
		utmInfo.UTMContent = &utmContent
	}

	// Build the request
	req := &showcheckout.ShowCheckoutRequest{
		OfferUUID:  offerUUID,
		ClientInfo: clientInfo,
		UTMInfo:    utmInfo,
	}

	// Extract optional tracking parameters
	if aff := c.Query("aff"); aff != "" {
		req.Aff = &aff
	}
	if fbclid := c.Query("fbclid"); fbclid != "" {
		req.Fbclid = &fbclid
	}
	if gclid := c.Query("gclid"); gclid != "" {
		req.Gclid = &gclid
	}
	if ttclid := c.Query("ttclid"); ttclid != "" {
		req.Ttclid = &ttclid
	}
	if clickId := c.Query("clickId"); clickId != "" {
		req.ClickID = &clickId
	}
	if originalUrl := c.Query("originalUrl"); originalUrl != "" {
		// Validate URL
		if _, err := url.Parse(originalUrl); err == nil {
			req.OriginalURL = &originalUrl
		}
	}

	// Extract cookie from headers
	if cookie := c.GetHeader("Cookie"); cookie != "" {
		req.Cookie = &cookie
	}

	return req, nil
} 