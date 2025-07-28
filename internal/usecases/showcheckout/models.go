package showcheckout

// ShowCheckoutRequest represents the input for the ShowCheckout use case
type ShowCheckoutRequest struct {
	OfferUUID   string     `json:"offer_uuid" validate:"required,uuid"`
	Aff         *string    `json:"aff,omitempty"`
	Cookie      *string    `json:"cookie,omitempty"`
	ClientInfo  ClientInfo `json:"client_info" validate:"required"`
	UTMInfo     UTMInfo    `json:"utm_info"`
	OriginalURL *string    `json:"original_url,omitempty" validate:"omitempty,url"`
	Fbclid      *string    `json:"fbclid,omitempty"`
	Gclid       *string    `json:"gclid,omitempty"`
	Ttclid      *string    `json:"ttclid,omitempty"`
	ClickID     *string    `json:"click_id,omitempty"`
}

// ClientInfo contains client device and location information
type ClientInfo struct {
	IP             *string `json:"ip,omitempty" validate:"omitempty,ip"`
	UserAgent      *string `json:"user_agent,omitempty"`
	IsMobile       bool    `json:"is_mobile"`
	Browser        *string `json:"browser,omitempty" validate:"omitempty,min=1"`
	BrowserVersion *string `json:"browser_version,omitempty" validate:"omitempty,min=1"`
	OS             *string `json:"os,omitempty" validate:"omitempty,min=1"`
	OSVersion      *string `json:"os_version,omitempty" validate:"omitempty,min=1"`
	Country        *string `json:"country,omitempty" validate:"omitempty,min=2"`
	State          *string `json:"state,omitempty" validate:"omitempty,min=2"`
	City           *string `json:"city,omitempty" validate:"omitempty,min=1"`
	Lat            *string `json:"lat,omitempty"`
	Lon            *string `json:"lon,omitempty"`
}

// UTMInfo contains UTM tracking parameters
type UTMInfo struct {
	Src         *string `json:"src,omitempty" validate:"omitempty,min=1"`
	UTMSource   *string `json:"utm_source,omitempty" validate:"omitempty,min=1"`
	UTMMedium   *string `json:"utm_medium,omitempty" validate:"omitempty,min=1"`
	UTMCampaign *string `json:"utm_campaign,omitempty" validate:"omitempty,min=1"`
	UTMTerm     *string `json:"utm_term,omitempty" validate:"omitempty,min=1"`
	UTMContent  *string `json:"utm_content,omitempty" validate:"omitempty,min=1"`
}

// ShowCheckoutResponse represents the output of the ShowCheckout use case
type ShowCheckoutResponse struct {
	BillingType         string                     `json:"billing_type"`
	IsFree              bool                       `json:"is_free"`
	BackRedirectURL     *string                    `json:"back_redirect_url,omitempty"`
	Config              CheckoutConfig             `json:"config"`
	OrderBumps          []ResponseOrderBump        `json:"order_bumps"`
	Product             ResponseProduct            `json:"product"`
	Reviews             []ResponseReview           `json:"reviews"`
	Pixels              []ResponsePixel            `json:"pixels"`
	Company             *ResponseCompany           `json:"company,omitempty"`
	AffiliateSettings   *ResponseAffiliateSettings `json:"affiliate_settings,omitempty"`
	Customer            *ResponseCustomer          `json:"customer,omitempty"`
	Plans               []ResponsePlan             `json:"plans,omitempty"`
	GooglePayMerchantID *string                    `json:"google_pay_merchant_id,omitempty"`
}

// CheckoutConfig contains checkout configuration settings
type CheckoutConfig struct {
	CheckoutUUID                string  `json:"checkout_uuid"`
	CheckoutDate                string  `json:"checkout_date"`
	HasDiscount                 bool    `json:"has_discount"`
	Favicon                     *string `json:"favicon,omitempty"`
	LogoEnabled                 bool    `json:"logo_enabled"`
	Logo                        *string `json:"logo,omitempty"`
	LogoPosition                string  `json:"logo_position"`
	BannerEnabled               bool    `json:"banner_enabled"`
	Banner                      *string `json:"banner,omitempty"`
	BackgroundType              string  `json:"background_type"`
	BackgroundColor             string  `json:"background_color"`
	ColorPrimary                string  `json:"color_primary"`
	ColorSecondary              string  `json:"color_secondary"`
	ColorBuyButton              string  `json:"color_buy_button"`
	AdsTextEnabled              bool    `json:"ads_text_enabled"`
	AdsText                     *string `json:"ads_text,omitempty"`
	CPFEnabled                  bool    `json:"cpf_enabled"`
	CNPJEnabled                 bool    `json:"cnpj_enabled"`
	BankSlipEnabled             bool    `json:"bank_slip_enabled"`
	CreditCardEnabled           bool    `json:"credit_card_enabled"`
	PixEnabled                  bool    `json:"pix_enabled"`
	NupayEnabled                bool    `json:"nupay_enabled"`
	PicpayEnabled               bool    `json:"picpay_enabled"`
	ApplePayEnabled             bool    `json:"apple_pay_enabled"`
	GooglePayEnabled            bool    `json:"google_pay_enabled"`
	AutomaticDiscountBankSlip   float64 `json:"automatic_discount_bank_slip"`
	AutomaticDiscountCreditCard float64 `json:"automatic_discount_credit_card"`
	AutomaticDiscountPix        float64 `json:"automatic_discount_pix"`
	AutomaticDiscountNupay      float64 `json:"automatic_discount_nupay"`
	AutomaticDiscountPicpay     float64 `json:"automatic_discount_picpay"`
	AutomaticDiscountApplePay   float64 `json:"automatic_discount_apple_pay"`
	AutomaticDiscountGooglePay  float64 `json:"automatic_discount_google_pay"`
	InstallmentsLimit           int     `json:"installments_limit"`
	PreselectedInstallment      int     `json:"preselected_installment"`
	InterestFreeInstallments    int     `json:"interest_free_installments"`
	ShowWebsiteAddress          bool    `json:"show_website_address"`
	ShowCompanyInfo             bool    `json:"show_company_info"`
	AddressRequired             bool    `json:"address_required"`
	WhatsappEnabled             bool    `json:"whatsapp_enabled"`
	SupportPhone                *string `json:"support_phone,omitempty"`
	SupportPhoneVerified        bool    `json:"support_phone_verified"`
	CountdownEnabled            bool    `json:"countdown_enabled"`
	CountdownTime               int     `json:"countdown_time"`
	CountdownFinishMessage      *string `json:"countdown_finish_message,omitempty"`
	NotificationsEnabled        bool    `json:"notifications_enabled"`
	SocialProofEnabled          bool    `json:"social_proof_enabled"`
	ReviewsEnabled              bool    `json:"reviews_enabled"`
}

// ResponseOrderBump represents an order bump offer
type ResponseOrderBump struct {
	UUID        string  `json:"uuid"`
	ProductName string  `json:"product_name"`
	Name        string  `json:"name"`
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Photo       *string `json:"photo,omitempty"`
	Format      string  `json:"format"`
	Order       int     `json:"order"`
}

// ResponseProduct represents product information
type ResponseProduct struct {
	UUID   string  `json:"uuid"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	Photo  *string `json:"photo,omitempty"`
	Format string  `json:"format"`
}

// ResponseReview represents a customer review
type ResponseReview struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Photo       *string `json:"photo,omitempty"`
	Stars       int     `json:"stars"`
}

// ResponsePixel represents tracking pixel information
type ResponsePixel struct {
	UUID                             string  `json:"uuid"`
	Events                           string  `json:"events"`
	Platform                         string  `json:"platform"`
	Code                             string  `json:"code"`
	IsAPI                            bool    `json:"is_api"`
	EnableBankslipPurchasePercentage bool    `json:"enable_bankslip_purchase_percentage"`
	EnablePixPurchasePercentage      bool    `json:"enable_pix_purchase_percentage"`
	BankSlipPurchasePercentage       float64 `json:"bank_slip_purchase_percentage"`
	PixPurchasePercentage            float64 `json:"pix_purchase_percentage"`
	GoogleAdsConversionLabel         *string `json:"google_ads_conversion_label,omitempty"`
}

// ResponseCompany represents company information
type ResponseCompany struct {
	FantasyName string `json:"fantasy_name"`
}

// ResponseAffiliateSettings represents affiliate configuration
type ResponseAffiliateSettings struct {
	CommissionPreference string `json:"commission_preference"`
	CookieLifetime       int    `json:"cookie_lifetime"`
}

// ResponseCustomer represents customer information
type ResponseCustomer struct {
	// Add customer fields as needed
}

// ResponsePlan represents a subscription plan
type ResponsePlan struct {
	UUID                    string  `json:"uuid"`
	Title                   string  `json:"title"`
	Tag                     *string `json:"tag,omitempty"`
	Price                   float64 `json:"price"`
	PromotionalPrice        float64 `json:"promotional_price"`
	FirstChargePriceEnabled bool    `json:"first_charge_price_enabled"`
	FirstChargePrice        float64 `json:"first_charge_price"`
	ChargeFrequency         string  `json:"charge_frequency"`
	IsDefault               bool    `json:"is_default"`
}

// Helper functions for pointer conversion
func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func Float64Ptr(f float64) *float64 {
	return &f
}

func BoolPtr(b bool) *bool {
	return &b
}
