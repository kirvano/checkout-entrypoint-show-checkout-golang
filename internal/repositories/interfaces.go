package repositories

import (
	"checkout-go/internal/core/entities"
	"context"
)

// OffersRepository defines the interface for offer data access
type OffersRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*Offer, error)
	Find(ctx context.Context, id int) (*Offer, error)
	IncrementCheckoutCount(ctx context.Context, uuid string) error
}

// ProductsRepository defines the interface for product data access
type ProductsRepository interface {
	Find(ctx context.Context, id int) (*Product, error)
}

// UsersRepository defines the interface for user data access
type UsersRepository interface {
	Find(ctx context.Context, id int) (*User, error)
}

// CompaniesRepository defines the interface for company data access
type CompaniesRepository interface {
	Find(ctx context.Context, id int) (*Company, error)
}

// FormatsRepository defines the interface for format data access
type FormatsRepository interface {
	Find(ctx context.Context, id int) (*Format, error)
}

// CheckoutConfigsRepository defines the interface for checkout config data access
type CheckoutConfigsRepository interface {
	Find(ctx context.Context, id int) (*CheckoutConfig, error)
}

// AffiliatesRepository defines the interface for affiliate data access
type AffiliatesRepository interface {
	FindByUUID(ctx context.Context, uuid string) (*Affiliate, error)
}

// ProductAffiliateSettingsRepository defines the interface for product affiliate settings data access
type ProductAffiliateSettingsRepository interface {
	FindByProduct(ctx context.Context, productID int) (*ProductAffiliateSettings, error)
}

// CheckoutsRepository defines the interface for checkout data access
type CheckoutsRepository interface {
	Create(ctx context.Context, checkout *entities.Checkout) error
	FindByUUID(ctx context.Context, uuid string) (*entities.Checkout, error)
	Update(ctx context.Context, checkout *entities.Checkout) error
}

// OrderBumpsRepository defines the interface for order bump data access
type OrderBumpsRepository interface {
	FindAllByOffer(ctx context.Context, offerID int) ([]*OrderBump, error)
}

// ReviewsRepository defines the interface for review data access
type ReviewsRepository interface {
	FindByCheckoutConfig(ctx context.Context, checkoutConfigID int) ([]*Review, error)
}

// PixelsRepository defines the interface for pixel data access
type PixelsRepository interface {
	FindAllByUserAndProduct(ctx context.Context, userID, productID int) ([]*Pixel, error)
}

// PlansRepository defines the interface for plan data access
type PlansRepository interface {
	FindByOffer(ctx context.Context, offerID int) ([]*Plan, error)
}

// DiscountsRepository defines the interface for discount data access
type DiscountsRepository interface {
	CheckHasDiscounts(ctx context.Context, productID int) (bool, error)
}

// FileDriver defines the interface for file operations
type FileDriver interface {
	GetBasePath() string
	GetFullPath(relativePath string) string
}

// Domain models for repositories
type Offer struct {
	ID                     int    `json:"id" dynamodb:"id"`
	UUID                   string `json:"uuid" dynamodb:"uuid"`
	ProductID              int    `json:"product_id" dynamodb:"product_id"`
	CheckoutConfigID       int    `json:"checkout_config_id" dynamodb:"checkout_config_id"`
	Status                 string `json:"status" dynamodb:"status"`
	IsTemporary            bool   `json:"is_temporary" dynamodb:"is_temporary"`
	Price                  int64  `json:"price" dynamodb:"price"` // stored as cents
	BillingType            string `json:"billing_type" dynamodb:"billing_type"`
	IsFree                 bool   `json:"is_free" dynamodb:"is_free"`
	BackRedirectURL        string `json:"back_redirect_url" dynamodb:"back_redirect_url"`
	BackRedirectURLEnabled bool   `json:"back_redirect_url_enabled" dynamodb:"back_redirect_url_enabled"`
	OrderBumpsEnabled      bool   `json:"order_bumps_enabled" dynamodb:"order_bumps_enabled"`
}

type Product struct {
	ID               int    `json:"id" dynamodb:"id"`
	UUID             string `json:"uuid" dynamodb:"uuid"`
	Name             string `json:"name" dynamodb:"name"`
	UserID           int    `json:"user_id" dynamodb:"user_id"`
	CompanyID        int    `json:"company_id" dynamodb:"company_id"`
	FormatID         int    `json:"format_id" dynamodb:"format_id"`
	Status           string `json:"status" dynamodb:"status"`
	EvaluationStatus string `json:"evaluation_status" dynamodb:"evaluation_status"`
	Currency         string `json:"currency" dynamodb:"currency"`
	PhotoURL         string `json:"photo_url" dynamodb:"photo_url"`
	SellerName       string `json:"seller_name" dynamodb:"seller_name"`
}

type User struct {
	ID            int    `json:"id" dynamodb:"id"`
	UUID          string `json:"uuid" dynamodb:"uuid"`
	Status        string `json:"status" dynamodb:"status"`
	BlockCheckout string `json:"block_checkout" dynamodb:"block_checkout"`
}

type Company struct {
	ID            int    `json:"id" dynamodb:"id"`
	Type          string `json:"type" dynamodb:"type"`
	MovingpayEcID string `json:"movingpay_ec_id" dynamodb:"movingpay_ec_id"`
}

type Format struct {
	ID   int    `json:"id" dynamodb:"id"`
	Slug string `json:"slug" dynamodb:"slug"`
}

type CheckoutConfig struct {
	ID                          int     `json:"id" dynamodb:"id"`
	UUID                        string  `json:"uuid" dynamodb:"uuid"`
	ShowCompanyInfo             bool    `json:"show_company_info" dynamodb:"show_company_info"`
	LogoURL                     string  `json:"logo_url" dynamodb:"logo_url"`
	BannerURL                   string  `json:"banner_url" dynamodb:"banner_url"`
	FaviconEnabled              bool    `json:"favicon_enabled" dynamodb:"favicon_enabled"`
	FaviconType                 string  `json:"favicon_type" dynamodb:"favicon_type"`
	FaviconURL                  string  `json:"favicon_url" dynamodb:"favicon_url"`
	LogoEnabled                 bool    `json:"logo_enabled" dynamodb:"logo_enabled"`
	LogoPosition                string  `json:"logo_position" dynamodb:"logo_position"`
	BannerEnabled               bool    `json:"banner_enabled" dynamodb:"banner_enabled"`
	BackgroundType              string  `json:"background_type" dynamodb:"background_type"`
	BackgroundColor             string  `json:"background_color" dynamodb:"background_color"`
	ColorPrimary                string  `json:"color_primary" dynamodb:"color_primary"`
	ColorSecondary              string  `json:"color_secondary" dynamodb:"color_secondary"`
	ColorBuyButton              string  `json:"color_buy_button" dynamodb:"color_buy_button"`
	AdsTextEnabled              bool    `json:"ads_text_enabled" dynamodb:"ads_text_enabled"`
	AdsText                     string  `json:"ads_text" dynamodb:"ads_text"`
	CPFEnabled                  bool    `json:"cpf_enabled" dynamodb:"cpf_enabled"`
	CNPJEnabled                 bool    `json:"cnpj_enabled" dynamodb:"cnpj_enabled"`
	BankSlipEnabled             bool    `json:"bank_slip_enabled" dynamodb:"bank_slip_enabled"`
	CreditCardEnabled           bool    `json:"credit_card_enabled" dynamodb:"credit_card_enabled"`
	PixEnabled                  bool    `json:"pix_enabled" dynamodb:"pix_enabled"`
	NupayEnabled                bool    `json:"nupay_enabled" dynamodb:"nupay_enabled"`
	PicpayEnabled               bool    `json:"picpay_enabled" dynamodb:"picpay_enabled"`
	ApplePayEnabled             bool    `json:"apple_pay_enabled" dynamodb:"apple_pay_enabled"`
	AutomaticDiscountBankSlip   float64 `json:"automatic_discount_bank_slip" dynamodb:"automatic_discount_bank_slip"`
	AutomaticDiscountCreditCard float64 `json:"automatic_discount_credit_card" dynamodb:"automatic_discount_credit_card"`
	AutomaticDiscountPix        float64 `json:"automatic_discount_pix" dynamodb:"automatic_discount_pix"`
	AutomaticDiscountNupay      float64 `json:"automatic_discount_nupay" dynamodb:"automatic_discount_nupay"`
	AutomaticDiscountPicpay     float64 `json:"automatic_discount_picpay" dynamodb:"automatic_discount_picpay"`
	AutomaticDiscountApplePay   float64 `json:"automatic_discount_apple_pay" dynamodb:"automatic_discount_apple_pay"`
	InstallmentsLimit           int     `json:"installments_limit" dynamodb:"installments_limit"`
	PreselectedInstallment      int     `json:"preselected_installment" dynamodb:"preselected_installment"`
	InterestFreeInstallments    int     `json:"interest_free_installments" dynamodb:"interest_free_installments"`
	ShowWebsiteAddress          bool    `json:"show_website_address" dynamodb:"show_website_address"`
	AddressRequired             bool    `json:"address_required" dynamodb:"address_required"`
	WhatsappEnabled             bool    `json:"whatsapp_enabled" dynamodb:"whatsapp_enabled"`
	SupportPhone                string  `json:"support_phone" dynamodb:"support_phone"`
	SupportPhoneVerified        bool    `json:"support_phone_verified" dynamodb:"support_phone_verified"`
	CountdownEnabled            bool    `json:"countdown_enabled" dynamodb:"countdown_enabled"`
	CountdownTime               int     `json:"countdown_time" dynamodb:"countdown_time"`
	CountdownFinishMessage      string  `json:"countdown_finish_message" dynamodb:"countdown_finish_message"`
	NotificationsEnabled        bool    `json:"notifications_enabled" dynamodb:"notifications_enabled"`
	SocialProofEnabled          bool    `json:"social_proof_enabled" dynamodb:"social_proof_enabled"`
	ReviewsEnabled              bool    `json:"reviews_enabled" dynamodb:"reviews_enabled"`
}

type Affiliate struct {
	ID     int    `json:"id" dynamodb:"id"`
	UUID   string `json:"uuid" dynamodb:"uuid"`
	UserID int    `json:"user_id" dynamodb:"user_id"`
}

type ProductAffiliateSettings struct {
	ID                   int      `json:"id" dynamodb:"id"`
	ProductID            int      `json:"product_id" dynamodb:"product_id"`
	CommissionPreference string   `json:"commission_preference" dynamodb:"commission_preference"`
	CookieLifetime       int      `json:"cookie_lifetime" dynamodb:"cookie_lifetime"`
	LastOffers           []string `json:"last_offers" dynamodb:"last_offers"`
}

func (pas *ProductAffiliateSettings) GetCookieLifetimeInDays() int {
	return pas.CookieLifetime
}

type OrderBump struct {
	ID             int    `json:"id" dynamodb:"id"`
	OfferID        int    `json:"offer_id" dynamodb:"offer_id"`
	OfferedOfferID int    `json:"offered_offer_id" dynamodb:"offered_offer_id"`
	Name           string `json:"name" dynamodb:"name"`
	Tag            string `json:"tag" dynamodb:"tag"`
	Description    string `json:"description" dynamodb:"description"`
	Order          int    `json:"order" dynamodb:"order"`
}

type Review struct {
	ID                 int    `json:"id" dynamodb:"id"`
	Name               string `json:"name" dynamodb:"name"`
	Description        string `json:"description" dynamodb:"description"`
	PhotoURL           string `json:"photo_url" dynamodb:"photo_url"`
	Stars              int    `json:"stars" dynamodb:"stars"`
	CheckoutConfigUUID string `json:"checkout_config_uuid" dynamodb:"checkout_config_uuid"`
}

type Pixel struct {
	ID                               int     `json:"id" dynamodb:"id"`
	UUID                             string  `json:"uuid" dynamodb:"uuid"`
	Events                           string  `json:"events" dynamodb:"events"`
	Platform                         string  `json:"platform" dynamodb:"platform"`
	Code                             string  `json:"code" dynamodb:"code"`
	Status                           bool    `json:"status" dynamodb:"status"`
	IsAPI                            bool    `json:"is_api" dynamodb:"is_api"`
	EnableBankslipPurchasePercentage bool    `json:"enable_bankslip_purchase_percentage" dynamodb:"enable_bankslip_purchase_percentage"`
	EnablePixPurchasePercentage      bool    `json:"enable_pix_purchase_percentage" dynamodb:"enable_pix_purchase_percentage"`
	BankSlipPurchasePercentage       float64 `json:"bank_slip_purchase_percentage" dynamodb:"bank_slip_purchase_percentage"`
	PixPurchasePercentage            float64 `json:"pix_purchase_percentage" dynamodb:"pix_purchase_percentage"`
	GoogleAdsConversionLabel         string  `json:"google_ads_conversion_label" dynamodb:"google_ads_conversion_label"`
}

type Plan struct {
	ID                      int    `json:"id" dynamodb:"id"`
	UUID                    string `json:"uuid" dynamodb:"uuid"`
	Title                   string `json:"title" dynamodb:"title"`
	Tag                     string `json:"tag" dynamodb:"tag"`
	Price                   int64  `json:"price" dynamodb:"price"`                         // stored as cents
	PromotionalPrice        int64  `json:"promotional_price" dynamodb:"promotional_price"` // stored as cents
	FirstChargePriceEnabled bool   `json:"first_charge_price_enabled" dynamodb:"first_charge_price_enabled"`
	FirstChargePrice        int64  `json:"first_charge_price" dynamodb:"first_charge_price"` // stored as cents
	ChargeFrequency         string `json:"charge_frequency" dynamodb:"charge_frequency"`
	IsDefault               bool   `json:"is_default" dynamodb:"is_default"`
}

// Constants for status and type enumerations
const (
	OfferStatusActive              = "ACTIVE"
	ProductStatusActive            = "ACTIVE"
	UserStatusActive               = "ACTIVE"
	UserBlockCheckoutActive        = "ACTIVE"
	ProductEvaluationStatusRefused = "REFUSED"
	CompanyTypeLegalPerson         = "LEGAL_PERSON"
	CheckoutConfigFaviconTypeFile  = "FILE"
	OfferBillingTypeOneTime        = "ONE_TIME"
)
