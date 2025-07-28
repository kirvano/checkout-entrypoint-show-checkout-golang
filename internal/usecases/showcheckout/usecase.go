package showcheckout

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"checkout-go/internal/core/entities"
	"checkout-go/internal/core/errors"
	"checkout-go/internal/core/valueobjects"
	"checkout-go/internal/repositories"
)

// UseCase implements the ShowCheckout business logic
type UseCase struct {
	offersRepo                   repositories.OffersRepository
	productsRepo                 repositories.ProductsRepository
	usersRepo                    repositories.UsersRepository
	companiesRepo                repositories.CompaniesRepository
	formatsRepo                  repositories.FormatsRepository
	checkoutConfigsRepo          repositories.CheckoutConfigsRepository
	affiliatesRepo               repositories.AffiliatesRepository
	productAffiliateSettingsRepo repositories.ProductAffiliateSettingsRepository
	checkoutsRepo                repositories.CheckoutsRepository
	orderBumpsRepo               repositories.OrderBumpsRepository
	reviewsRepo                  repositories.ReviewsRepository
	pixelsRepo                   repositories.PixelsRepository
	plansRepo                    repositories.PlansRepository
	discountsRepo                repositories.DiscountsRepository
	fileDriver                   repositories.FileDriver
}

// NewUseCase creates a new ShowCheckout use case
func NewUseCase(
	offersRepo repositories.OffersRepository,
	productsRepo repositories.ProductsRepository,
	usersRepo repositories.UsersRepository,
	companiesRepo repositories.CompaniesRepository,
	formatsRepo repositories.FormatsRepository,
	checkoutConfigsRepo repositories.CheckoutConfigsRepository,
	affiliatesRepo repositories.AffiliatesRepository,
	productAffiliateSettingsRepo repositories.ProductAffiliateSettingsRepository,
	checkoutsRepo repositories.CheckoutsRepository,
	orderBumpsRepo repositories.OrderBumpsRepository,
	reviewsRepo repositories.ReviewsRepository,
	pixelsRepo repositories.PixelsRepository,
	plansRepo repositories.PlansRepository,
	discountsRepo repositories.DiscountsRepository,
	fileDriver repositories.FileDriver,
) *UseCase {
	return &UseCase{
		offersRepo:                   offersRepo,
		productsRepo:                 productsRepo,
		usersRepo:                    usersRepo,
		companiesRepo:                companiesRepo,
		formatsRepo:                  formatsRepo,
		checkoutConfigsRepo:          checkoutConfigsRepo,
		affiliatesRepo:               affiliatesRepo,
		productAffiliateSettingsRepo: productAffiliateSettingsRepo,
		checkoutsRepo:                checkoutsRepo,
		orderBumpsRepo:               orderBumpsRepo,
		reviewsRepo:                  reviewsRepo,
		pixelsRepo:                   pixelsRepo,
		plansRepo:                    plansRepo,
		discountsRepo:                discountsRepo,
		fileDriver:                   fileDriver,
	}
}

// Execute performs the ShowCheckout use case
func (uc *UseCase) Execute(ctx context.Context, req *ShowCheckoutRequest) (*ShowCheckoutResponse, error) {
	// Validate UUID
	if !valueobjects.IsValidUUID(req.OfferUUID) {
		return nil, errors.NewDontWorryError(StringPtr("A oferta informada é inválida"))
	}

	// Get offer
	offer, err := uc.offersRepo.FindByUUID(ctx, req.OfferUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to find offer: %w", err)
	}

	if offer == nil || offer.Status != repositories.OfferStatusActive || offer.IsTemporary {
		return nil, errors.NewDontWorryError(StringPtr("Oferta não encontrada"))
	}

	// Get product
	product, err := uc.productsRepo.Find(ctx, offer.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	if product == nil || product.Status != repositories.ProductStatusActive || product.EvaluationStatus == repositories.ProductEvaluationStatusRefused {
		return nil, errors.NewDontWorryError(StringPtr("Produto não encontrado"))
	}

	// Get user
	user, err := uc.usersRepo.Find(ctx, product.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil || user.Status != repositories.UserStatusActive {
		return nil, errors.NewDontWorryError(StringPtr("Usuário não encontrado"))
	}

	// Check if user is blocked from checkout (empty means not blocked)
	if user.BlockCheckout != "" && user.BlockCheckout != repositories.UserBlockCheckoutActive {
		return nil, errors.NewDontWorryError(StringPtr("Usuário não encontrado"))
	}

	// Get company
	company, err := uc.companiesRepo.Find(ctx, product.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to find company: %w", err)
	}

	if company == nil {
		return nil, errors.NewDontWorryError(StringPtr("Empresa não encontrada"))
	}

	// Get product format
	productFormat, err := uc.formatsRepo.Find(ctx, product.FormatID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product format: %w", err)
	}

	if productFormat == nil {
		return nil, errors.NewDontWorryError(StringPtr("Formato do produto não encontrado"))
	}

	// Get checkout config
	checkoutConfig, err := uc.checkoutConfigsRepo.Find(ctx, offer.CheckoutConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to find checkout config: %w", err)
	}

	if checkoutConfig == nil {
		return nil, errors.NewDontWorryError(StringPtr("Configuração de checkout não encontrada"))
	}

	// Handle affiliate logic
	userID := product.UserID
	var affiliateID *int
	var productAffiliateSettings *repositories.ProductAffiliateSettings

	affiliateFromQueryString := req.Aff
	affiliateCookieID := fmt.Sprintf("aff.%s", product.UUID)
	affiliateUUIDFromCookie := uc.getCookie(affiliateCookieID, req.Cookie)

	var affiliateUUID *string
	if affiliateFromQueryString != nil {
		affiliateUUID = affiliateFromQueryString
	} else if affiliateUUIDFromCookie != nil {
		affiliateUUID = affiliateUUIDFromCookie
	}

	if affiliateUUID != nil {
		affiliate, err := uc.affiliatesRepo.FindByUUID(ctx, *affiliateUUID)
		if err != nil {
			log.Printf("Failed to find affiliate: %v", err)
		} else if affiliate != nil {
			userAffiliate, err := uc.usersRepo.Find(ctx, affiliate.UserID)
			if err != nil {
				return nil, fmt.Errorf("failed to find affiliate user: %w", err)
			}

			if userAffiliate == nil || userAffiliate.Status != repositories.UserStatusActive || userAffiliate.BlockCheckout != repositories.UserBlockCheckoutActive {
				return nil, errors.NewDontWorryError(StringPtr("Afiliado não encontrado"))
			}

			productAffiliateSettings, err = uc.productAffiliateSettingsRepo.FindByProduct(ctx, product.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to find product affiliate settings: %w", err)
			}

			if productAffiliateSettings == nil {
				return nil, errors.NewDontWorryError(StringPtr("Configurações de afiliação não encontradas"))
			}

			if productAffiliateSettings.LastOffers == nil || !uc.contains(productAffiliateSettings.LastOffers, req.OfferUUID) {
				return nil, errors.NewDontWorryError(StringPtr("Afiliado não autorizado"))
			}

			affiliateID = &affiliate.ID
			userID = affiliate.UserID
		}
	}

	// Extract pixel data
	pixelData := uc.extractPixelData(req)

	// Create checkout
	checkout := entities.NewCheckout(entities.CheckoutProps{
		OfferID:        &offer.ID,
		ProductID:      product.ID,
		AffiliateID:    affiliateID,
		Currency:       product.Currency,
		UserAgent:      req.ClientInfo.UserAgent,
		OS:             req.ClientInfo.OS,
		Browser:        req.ClientInfo.Browser,
		BrowserVersion: req.ClientInfo.BrowserVersion,
		IsMobile:       req.ClientInfo.IsMobile,
		IP:             req.ClientInfo.IP,
		City:           req.ClientInfo.City,
		State:          req.ClientInfo.State,
		Lat:            req.ClientInfo.Lat,
		Lon:            req.ClientInfo.Lon,
		Country:        req.ClientInfo.Country,
		Src:            req.UTMInfo.Src,
		UTMSource:      req.UTMInfo.UTMSource,
		UTMMedium:      req.UTMInfo.UTMMedium,
		UTMCampaign:    req.UTMInfo.UTMCampaign,
		UTMTerm:        req.UTMInfo.UTMTerm,
		UTMContent:     req.UTMInfo.UTMContent,
		PixelData:      pixelData,
		OriginalURL:    req.OriginalURL,
	})

	// Debug: Log checkout details before saving
	fmt.Printf("DEBUG: Creating checkout with UUID: '%s', ProductID: %d\n", checkout.UUID, checkout.ProductID)
	
	// Save checkout
	if err := uc.checkoutsRepo.Create(ctx, checkout); err != nil {
		return nil, fmt.Errorf("failed to create checkout: %w", err)
	}

	// Increment checkout count
	if err := uc.offersRepo.IncrementCheckoutCount(ctx, req.OfferUUID); err != nil {
		log.Printf("Failed to increment checkout count: %v", err)
	}

	// Build order bumps
	responseOrderBumps, err := uc.buildOrderBumps(ctx, offer)
	if err != nil {
		log.Printf("Failed to build order bumps: %v", err)
		responseOrderBumps = []ResponseOrderBump{}
	}

	// Build reviews
	responseReviews, err := uc.buildReviews(ctx, checkoutConfig)
	if err != nil {
		log.Printf("Failed to build reviews: %v", err)
		responseReviews = []ResponseReview{}
	}

	// Build pixels
	responsePixels, err := uc.buildPixels(ctx, userID, product.ID)
	if err != nil {
		log.Printf("Failed to build pixels: %v", err)
		responsePixels = []ResponsePixel{}
	}

	// Build plans
	responsePlans, err := uc.buildPlans(ctx, offer.ID)
	if err != nil {
		log.Printf("Failed to build plans: %v", err)
		responsePlans = []ResponsePlan{}
	}

	// Build affiliate settings
	var affiliateSettings *ResponseAffiliateSettings
	if productAffiliateSettings != nil {
		affiliateSettings = &ResponseAffiliateSettings{
			CommissionPreference: productAffiliateSettings.CommissionPreference,
			CookieLifetime:       productAffiliateSettings.GetCookieLifetimeInDays(),
		}
	}

	// Build company response
	var responseCompany *ResponseCompany
	if checkoutConfig.ShowCompanyInfo && company.Type == repositories.CompanyTypeLegalPerson && product.SellerName != "" {
		responseCompany = &ResponseCompany{
			FantasyName: product.SellerName,
		}
	}

	// Build configuration URLs
	var logo *string
	if checkoutConfig.LogoURL != "" {
		logoURL := uc.fileDriver.GetFullPath(checkoutConfig.LogoURL)
		logo = &logoURL
	}

	var banner *string
	if checkoutConfig.BannerURL != "" {
		bannerURL := uc.fileDriver.GetFullPath(checkoutConfig.BannerURL)
		banner = &bannerURL
	}

	var favicon *string
	if checkoutConfig.FaviconEnabled {
		if checkoutConfig.FaviconType == repositories.CheckoutConfigFaviconTypeFile {
			if checkoutConfig.FaviconURL != "" {
				faviconURL := uc.fileDriver.GetFullPath(checkoutConfig.FaviconURL)
				favicon = &faviconURL
			}
		} else if checkoutConfig.LogoEnabled && logo != nil {
			favicon = logo
		}
	}

	// Check for discounts
	hasDiscount, err := uc.discountsRepo.CheckHasDiscounts(ctx, product.ID)
	if err != nil {
		log.Printf("Failed to check discounts: %v", err)
		hasDiscount = false
	}

	// Determine payment options
	creditCardEnabled := checkoutConfig.CreditCardEnabled && (!uc.isProduction() || company.MovingpayEcID != "")
	applePayEnabled := checkoutConfig.ApplePayEnabled && offer.BillingType == repositories.OfferBillingTypeOneTime
	googlePayEnabled := checkoutConfig.GooglePayEnabled && offer.BillingType == repositories.OfferBillingTypeOneTime

	// Build response
	response := &ShowCheckoutResponse{
		BillingType:         offer.BillingType,
		IsFree:              offer.IsFree,
		BackRedirectURL:     uc.getBackRedirectURL(offer),
		GooglePayMerchantID: uc.getGooglePayMerchantID(checkoutConfig),
		Config: CheckoutConfig{
			CheckoutUUID:                checkout.GetUUID(),
			CheckoutDate:                checkout.CreatedAt.Format(time.RFC3339),
			HasDiscount:                 hasDiscount,
			Favicon:                     favicon,
			LogoEnabled:                 checkoutConfig.LogoEnabled,
			Logo:                        logo,
			LogoPosition:                checkoutConfig.LogoPosition,
			BannerEnabled:               checkoutConfig.BannerEnabled,
			Banner:                      banner,
			BackgroundType:              checkoutConfig.BackgroundType,
			BackgroundColor:             checkoutConfig.BackgroundColor,
			ColorPrimary:                checkoutConfig.ColorPrimary,
			ColorSecondary:              checkoutConfig.ColorSecondary,
			ColorBuyButton:              checkoutConfig.ColorBuyButton,
			AdsTextEnabled:              checkoutConfig.AdsTextEnabled,
			AdsText:                     uc.getAdsText(checkoutConfig),
			CPFEnabled:                  checkoutConfig.CPFEnabled,
			CNPJEnabled:                 checkoutConfig.CNPJEnabled,
			BankSlipEnabled:             checkoutConfig.BankSlipEnabled && offer.BillingType == repositories.OfferBillingTypeOneTime,
			CreditCardEnabled:           creditCardEnabled,
			PixEnabled:                  checkoutConfig.PixEnabled,
			NupayEnabled:                false, // TODO: Enable Nupay
			PicpayEnabled:               checkoutConfig.PicpayEnabled && offer.BillingType == repositories.OfferBillingTypeOneTime,
			ApplePayEnabled:             applePayEnabled,
			GooglePayEnabled:            googlePayEnabled,
			AutomaticDiscountBankSlip:   checkoutConfig.AutomaticDiscountBankSlip,
			AutomaticDiscountCreditCard: checkoutConfig.AutomaticDiscountCreditCard,
			AutomaticDiscountPix:        checkoutConfig.AutomaticDiscountPix,
			AutomaticDiscountNupay:      checkoutConfig.AutomaticDiscountNupay,
			AutomaticDiscountPicpay:     checkoutConfig.AutomaticDiscountPicpay,
			AutomaticDiscountApplePay:   checkoutConfig.AutomaticDiscountApplePay,
			AutomaticDiscountGooglePay:  checkoutConfig.AutomaticDiscountGooglePay,
			InstallmentsLimit:           checkoutConfig.InstallmentsLimit,
			PreselectedInstallment:      checkoutConfig.PreselectedInstallment,
			InterestFreeInstallments:    checkoutConfig.InterestFreeInstallments,
			ShowWebsiteAddress:          checkoutConfig.ShowWebsiteAddress,
			ShowCompanyInfo:             checkoutConfig.ShowCompanyInfo,
			AddressRequired:             checkoutConfig.AddressRequired,
			WhatsappEnabled:             checkoutConfig.WhatsappEnabled,
			SupportPhone:                uc.getSupportPhone(checkoutConfig),
			SupportPhoneVerified:        checkoutConfig.SupportPhoneVerified,
			CountdownEnabled:            checkoutConfig.CountdownEnabled,
			CountdownTime:               checkoutConfig.CountdownTime,
			CountdownFinishMessage:      uc.getCountdownFinishMessage(checkoutConfig),
			NotificationsEnabled:        checkoutConfig.NotificationsEnabled,
			SocialProofEnabled:          checkoutConfig.SocialProofEnabled,
			ReviewsEnabled:              checkoutConfig.ReviewsEnabled,
		},
		OrderBumps:        responseOrderBumps,
		Product: ResponseProduct{
			UUID:   product.UUID,
			Name:   product.Name,
			Price:  uc.databaseToFloat(offer.Price),
			Photo:  uc.getProductPhoto(product),
			Format: productFormat.Slug,
		},
		Reviews:           responseReviews,
		Pixels:            responsePixels,
		Company:           responseCompany,
		AffiliateSettings: affiliateSettings,
		Customer:          nil, // Not implemented in original
		Plans:             responsePlans,
	}

	return response, nil
}

// Helper methods

func (uc *UseCase) extractPixelData(req *ShowCheckoutRequest) map[string]interface{} {
	pixelData := make(map[string]interface{})

	// Extract fbclid
	var fbclid *string
	if req.Fbclid != nil {
		fbclid = req.Fbclid
	} else {
		fbclid = uc.getCookie("_fbc", req.Cookie)
	}
	if fbclid != nil {
		pixelData["fbclid"] = *fbclid
	}

	// Extract fbp
	if fbp := uc.getCookie("_fbp", req.Cookie); fbp != nil {
		pixelData["fbp"] = *fbp
	}

	// Extract gclid
	var gclid *string
	if req.Gclid != nil {
		gclid = req.Gclid
	} else {
		gclid = uc.getCookie("_gcl_au", req.Cookie)
	}
	if gclid != nil {
		pixelData["gclid"] = *gclid
	}

	// Extract ttclid
	var ttclid *string
	if req.Ttclid != nil {
		ttclid = req.Ttclid
	} else {
		ttclid = uc.getCookie("ttclid", req.Cookie)
	}
	if ttclid != nil {
		pixelData["ttclid"] = *ttclid
	}

	// Extract ttp
	if ttp := uc.getCookie("_ttp", req.Cookie); ttp != nil {
		pixelData["ttp"] = *ttp
	}

	// Extract clickId
	if req.ClickID != nil {
		pixelData["click_id"] = *req.ClickID
	}

	if len(pixelData) == 0 {
		return nil
	}

	return pixelData
}

func (uc *UseCase) getCookie(name string, cookie *string) *string {
	if cookie == nil || *cookie == "" {
		return nil
	}

	// Simple cookie parsing - in production you might want a more robust implementation
	parts := strings.Split(*cookie, ";")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if strings.HasPrefix(trimmed, name+"=") {
			value := strings.TrimPrefix(trimmed, name+"=")
			if value != "" {
				return &value
			}
		}
	}
	return nil
}

func (uc *UseCase) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (uc *UseCase) databaseToFloat(value int64) float64 {
	return float64(value) / 100.0 // Convert cents to dollars
}

func (uc *UseCase) getGooglePayMerchantID(checkoutConfig *repositories.CheckoutConfig) *string {
	if checkoutConfig.GooglePayMerchantID != "" {
		return &checkoutConfig.GooglePayMerchantID
	}
	return nil
}

func (uc *UseCase) isProduction() bool {
	// TODO: Implement environment check
	return false
}

func (uc *UseCase) getBackRedirectURL(offer *repositories.Offer) *string {
	if offer.BackRedirectURLEnabled && offer.BackRedirectURL != "" {
		return &offer.BackRedirectURL
	}
	return nil
}

func (uc *UseCase) getAdsText(config *repositories.CheckoutConfig) *string {
	if config.AdsText != "" {
		return &config.AdsText
	}
	return nil
}

func (uc *UseCase) getSupportPhone(config *repositories.CheckoutConfig) *string {
	if config.SupportPhone != "" {
		return &config.SupportPhone
	}
	return nil
}

func (uc *UseCase) getCountdownFinishMessage(config *repositories.CheckoutConfig) *string {
	if config.CountdownFinishMessage != "" {
		return &config.CountdownFinishMessage
	}
	return nil
}

func (uc *UseCase) getProductPhoto(product *repositories.Product) *string {
	if product.PhotoURL != "" {
		photoURL := uc.fileDriver.GetFullPath(product.PhotoURL)
		return &photoURL
	}
	return nil
}

// Build methods for complex data structures

func (uc *UseCase) buildOrderBumps(ctx context.Context, offer *repositories.Offer) ([]ResponseOrderBump, error) {
	var responseOrderBumps []ResponseOrderBump

	if !offer.OrderBumpsEnabled {
		return responseOrderBumps, nil
	}

	orderBumps, err := uc.orderBumpsRepo.FindAllByOffer(ctx, offer.ID)
	if err != nil {
		return nil, err
	}

	for _, orderBump := range orderBumps {
		offeredOffer, err := uc.offersRepo.Find(ctx, orderBump.OfferedOfferID)
		if err != nil || offeredOffer == nil || offeredOffer.Status != repositories.OfferStatusActive {
			continue
		}

		product, err := uc.productsRepo.Find(ctx, offeredOffer.ProductID)
		if err != nil || product == nil || product.Status != repositories.ProductStatusActive || product.EvaluationStatus == repositories.ProductEvaluationStatusRefused {
			continue
		}

		format, err := uc.formatsRepo.Find(ctx, product.FormatID)
		if err != nil || format == nil {
			continue
		}

		price := uc.databaseToFloat(offeredOffer.Price)
		photo := uc.getProductPhoto(product)

		responseOrderBumps = append(responseOrderBumps, ResponseOrderBump{
			UUID:        offeredOffer.UUID,
			ProductName: product.Name,
			Name:        orderBump.Name,
			Tag:         orderBump.Tag,
			Description: orderBump.Description,
			Price:       price,
			Photo:       photo,
			Format:      format.Slug,
			Order:       orderBump.Order,
		})
	}

	return responseOrderBumps, nil
}

func (uc *UseCase) buildReviews(ctx context.Context, checkoutConfig *repositories.CheckoutConfig) ([]ResponseReview, error) {
	reviews, err := uc.reviewsRepo.FindByCheckoutConfig(ctx, checkoutConfig.ID)
	if err != nil {
		return nil, err
	}

	var responseReviews []ResponseReview
	for _, review := range reviews {
		var photo *string
		if review.PhotoURL != "" {
			photoURL := uc.fileDriver.GetFullPath(review.PhotoURL)
			photo = &photoURL
		}

		var description *string
		if review.Description != "" {
			description = &review.Description
		}

		responseReviews = append(responseReviews, ResponseReview{
			Name:        review.Name,
			Description: description,
			Photo:       photo,
			Stars:       review.Stars,
		})
	}

	return responseReviews, nil
}

func (uc *UseCase) buildPixels(ctx context.Context, userID, productID int) ([]ResponsePixel, error) {
	pixels, err := uc.pixelsRepo.FindAllByUserAndProduct(ctx, userID, productID)
	if err != nil {
		return nil, err
	}

	var responsePixels []ResponsePixel
	for _, pixel := range pixels {
		if !pixel.Status {
			continue
		}

		var googleAdsLabel *string
		if pixel.GoogleAdsConversionLabel != "" {
			googleAdsLabel = &pixel.GoogleAdsConversionLabel
		}

		responsePixels = append(responsePixels, ResponsePixel{
			UUID:                             pixel.UUID,
			Events:                           pixel.Events,
			Platform:                         pixel.Platform,
			Code:                             pixel.Code,
			IsAPI:                            pixel.IsAPI,
			EnableBankslipPurchasePercentage: pixel.EnableBankslipPurchasePercentage,
			EnablePixPurchasePercentage:      pixel.EnablePixPurchasePercentage,
			BankSlipPurchasePercentage:       pixel.BankSlipPurchasePercentage,
			PixPurchasePercentage:            pixel.PixPurchasePercentage,
			GoogleAdsConversionLabel:         googleAdsLabel,
		})
	}

	return responsePixels, nil
}

func (uc *UseCase) buildPlans(ctx context.Context, offerID int) ([]ResponsePlan, error) {
	plans, err := uc.plansRepo.FindByOffer(ctx, offerID)
	if err != nil {
		return nil, err
	}

	var responsePlans []ResponsePlan
	for _, plan := range plans {
		var tag *string
		if plan.Tag != "" && plan.Tag != "Nenhum" {
			tag = &plan.Tag
		}

		responsePlans = append(responsePlans, ResponsePlan{
			UUID:                    plan.UUID,
			Title:                   plan.Title,
			Tag:                     tag,
			Price:                   uc.databaseToFloat(plan.Price),
			PromotionalPrice:        uc.databaseToFloat(plan.PromotionalPrice),
			FirstChargePriceEnabled: plan.FirstChargePriceEnabled,
			FirstChargePrice:        uc.databaseToFloat(plan.FirstChargePrice),
			ChargeFrequency:         plan.ChargeFrequency,
			IsDefault:               plan.IsDefault,
		})
	}

	return responsePlans, nil
}

