package entities

import (
	"time"

	"checkout-go/internal/core/valueobjects"
)

type CheckoutStatus string

const (
	CheckoutStatusAccessed      CheckoutStatus = "ACCESSED"
	CheckoutStatusAbandonedCart CheckoutStatus = "ABANDONED_CART"
	CheckoutStatusRecovered     CheckoutStatus = "RECOVERED"
	CheckoutStatusSaleFinalized CheckoutStatus = "SALE_FINALIZED"
)

type Checkout struct {
	ID                         *int                   `json:"id,omitempty" dynamodb:"id,omitempty"`
	UUID                       string                 `json:"uuid" dynamodb:"uuid"`
	Code                       *string                `json:"code,omitempty" dynamodb:"code,omitempty"`
	OfferID                    *int                   `json:"offer_id,omitempty" dynamodb:"offer_id,omitempty"`
	ProductID                  int                    `json:"product_id" dynamodb:"product_id"`
	AffiliateID                *int                   `json:"affiliate_id,omitempty" dynamodb:"affiliate_id,omitempty"`
	Status                     CheckoutStatus         `json:"status" dynamodb:"status"`
	UserAgent                  *string                `json:"user_agent,omitempty" dynamodb:"user_agent,omitempty"`
	OS                         *string                `json:"os,omitempty" dynamodb:"os,omitempty"`
	Browser                    *string                `json:"browser,omitempty" dynamodb:"browser,omitempty"`
	BrowserVersion             *string                `json:"browser_version,omitempty" dynamodb:"browser_version,omitempty"`
	IsMobile                   bool                   `json:"is_mobile" dynamodb:"is_mobile"`
	IP                         *string                `json:"ip,omitempty" dynamodb:"ip,omitempty"`
	City                       *string                `json:"city,omitempty" dynamodb:"city,omitempty"`
	State                      *string                `json:"state,omitempty" dynamodb:"state,omitempty"`
	Lat                        *string                `json:"lat,omitempty" dynamodb:"lat,omitempty"`
	Lon                        *string                `json:"lon,omitempty" dynamodb:"lon,omitempty"`
	Country                    *string                `json:"country,omitempty" dynamodb:"country,omitempty"`
	Currency                   string                 `json:"currency" dynamodb:"currency"`
	EmailSentAmount            int                    `json:"email_sent_amount" dynamodb:"email_sent_amount"`
	SMSSentAmount              int                    `json:"sms_sent_amount" dynamodb:"sms_sent_amount"`
	Src                        *string                `json:"src,omitempty" dynamodb:"src,omitempty"`
	UTMSource                  *string                `json:"utm_source,omitempty" dynamodb:"utm_source,omitempty"`
	UTMMedium                  *string                `json:"utm_medium,omitempty" dynamodb:"utm_medium,omitempty"`
	UTMCampaign                *string                `json:"utm_campaign,omitempty" dynamodb:"utm_campaign,omitempty"`
	UTMTerm                    *string                `json:"utm_term,omitempty" dynamodb:"utm_term,omitempty"`
	UTMContent                 *string                `json:"utm_content,omitempty" dynamodb:"utm_content,omitempty"`
	OSVersion                  *string                `json:"os_version,omitempty" dynamodb:"os_version,omitempty"`
	MercadoPagoDeviceSessionID *string                `json:"mercado_pago_device_session_id,omitempty" dynamodb:"mercado_pago_device_session_id,omitempty"`
	PixelData                  map[string]interface{} `json:"pixel_data,omitempty" dynamodb:"pixel_data,omitempty"`
	OriginalURL                *string                `json:"original_url,omitempty" dynamodb:"original_url,omitempty"`
	CreatedAt                  time.Time              `json:"created_at" dynamodb:"created_at"`
	UpdatedAt                  time.Time              `json:"updated_at" dynamodb:"updated_at"`
}

type CheckoutProps struct {
	Code                       *string
	OfferID                    *int
	ProductID                  int
	AffiliateID                *int
	UserAgent                  *string
	OS                         *string
	Browser                    *string
	BrowserVersion             *string
	IsMobile                   bool
	IP                         *string
	City                       *string
	State                      *string
	Lat                        *string
	Lon                        *string
	Country                    *string
	Currency                   string
	Src                        *string
	UTMSource                  *string
	UTMMedium                  *string
	UTMCampaign                *string
	UTMTerm                    *string
	UTMContent                 *string
	OSVersion                  *string
	MercadoPagoDeviceSessionID *string
	PixelData                  map[string]interface{}
	OriginalURL                *string
}

func NewCheckout(props CheckoutProps) *Checkout {
	now := time.Now()
	uuid := valueobjects.NewRandomUUID()

	return &Checkout{
		UUID:                       uuid.String(),
		Code:                       props.Code,
		OfferID:                    props.OfferID,
		ProductID:                  props.ProductID,
		AffiliateID:                props.AffiliateID,
		Status:                     CheckoutStatusAccessed,
		UserAgent:                  props.UserAgent,
		OS:                         props.OS,
		Browser:                    props.Browser,
		BrowserVersion:             props.BrowserVersion,
		IsMobile:                   props.IsMobile,
		IP:                         props.IP,
		City:                       props.City,
		State:                      props.State,
		Lat:                        props.Lat,
		Lon:                        props.Lon,
		Country:                    props.Country,
		Currency:                   props.Currency,
		EmailSentAmount:            0,
		SMSSentAmount:              0,
		Src:                        props.Src,
		UTMSource:                  props.UTMSource,
		UTMMedium:                  props.UTMMedium,
		UTMCampaign:                props.UTMCampaign,
		UTMTerm:                    props.UTMTerm,
		UTMContent:                 props.UTMContent,
		OSVersion:                  props.OSVersion,
		MercadoPagoDeviceSessionID: props.MercadoPagoDeviceSessionID,
		PixelData:                  props.PixelData,
		OriginalURL:                props.OriginalURL,
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}
}

// GetID returns the checkout ID
func (c *Checkout) GetID() *int {
	return c.ID
}

// GetUUID returns the checkout UUID
func (c *Checkout) GetUUID() string {
	return c.UUID
}

// SetID sets the checkout ID
func (c *Checkout) SetID(id int) {
	c.ID = &id
}

// UpdateTimestamp updates the UpdatedAt field
func (c *Checkout) UpdateTimestamp() {
	c.UpdatedAt = time.Now()
}

// IsAccessedStatus checks if checkout status is ACCESSED
func (c *Checkout) IsAccessedStatus() bool {
	return c.Status == CheckoutStatusAccessed
}

// IsAbandonedCartStatus checks if checkout status is ABANDONED_CART
func (c *Checkout) IsAbandonedCartStatus() bool {
	return c.Status == CheckoutStatusAbandonedCart
}

// IsRecoveredStatus checks if checkout status is RECOVERED
func (c *Checkout) IsRecoveredStatus() bool {
	return c.Status == CheckoutStatusRecovered
}

// IsSaleFinalizedStatus checks if checkout status is SALE_FINALIZED
func (c *Checkout) IsSaleFinalizedStatus() bool {
	return c.Status == CheckoutStatusSaleFinalized
}

// HasPixelData checks if checkout has pixel data
func (c *Checkout) HasPixelData() bool {
	return c.PixelData != nil && len(c.PixelData) > 0
}

// GetPixelValue gets a specific pixel data value by key
func (c *Checkout) GetPixelValue(key string) (interface{}, bool) {
	if c.PixelData == nil {
		return nil, false
	}
	val, exists := c.PixelData[key]
	return val, exists
}

// SetPixelValue sets a specific pixel data value
func (c *Checkout) SetPixelValue(key string, value interface{}) {
	if c.PixelData == nil {
		c.PixelData = make(map[string]interface{})
	}
	c.PixelData[key] = value
}
