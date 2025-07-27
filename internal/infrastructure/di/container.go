package di

import (
	"fmt"

	"checkout-go/internal/config"
	"checkout-go/internal/infrastructure/aws"
	"checkout-go/internal/infrastructure/dynamodb"
	"checkout-go/internal/repositories"
	"checkout-go/internal/usecases/showcheckout"
)

// Container holds all dependencies
type Container struct {
	// Configuration
	config *config.Config

	// Repositories
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

	// Use Cases
	showCheckoutUseCase *showcheckout.UseCase
}

// NewContainer creates and configures a new dependency injection container
func NewContainer() (*Container, error) {
	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize AWS configuration
	awsConfig, err := aws.NewConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AWS config: %w", err)
	}

	// Initialize DynamoDB client
	dynamoClient, err := dynamodb.NewClient(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DynamoDB client: %w", err)
	}

	// Initialize repositories
	offersRepo := dynamodb.NewOffersRepository(dynamoClient, cfg)
	productsRepo := dynamodb.NewProductsRepository(dynamoClient, cfg)
	usersRepo := dynamodb.NewUsersRepository(dynamoClient, cfg)
	companiesRepo := dynamodb.NewCompaniesRepository(dynamoClient, cfg)
	formatsRepo := dynamodb.NewFormatsRepository(dynamoClient, cfg)
	checkoutConfigsRepo := dynamodb.NewCheckoutConfigsRepository(dynamoClient, cfg)
	affiliatesRepo := dynamodb.NewAffiliatesRepository(dynamoClient, cfg)
	productAffiliateSettingsRepo := dynamodb.NewProductAffiliateSettingsRepository(dynamoClient, cfg)
	checkoutsRepo := dynamodb.NewCheckoutsRepository(dynamoClient, cfg)
	orderBumpsRepo := dynamodb.NewOrderBumpsRepository(dynamoClient, cfg)
	reviewsRepo := dynamodb.NewReviewsRepository(dynamoClient, cfg)
	pixelsRepo := dynamodb.NewPixelsRepository(dynamoClient, cfg)
	plansRepo := dynamodb.NewPlansRepository(dynamoClient, cfg)
	discountsRepo := dynamodb.NewDiscountsRepository(dynamoClient, cfg)

	// Initialize file driver (S3-based) with configuration
	fileDriver := aws.NewS3FileDriver(cfg)

	// Initialize use cases
	showCheckoutUseCase := showcheckout.NewUseCase(
		offersRepo,
		productsRepo,
		usersRepo,
		companiesRepo,
		formatsRepo,
		checkoutConfigsRepo,
		affiliatesRepo,
		productAffiliateSettingsRepo,
		checkoutsRepo,
		orderBumpsRepo,
		reviewsRepo,
		pixelsRepo,
		plansRepo,
		discountsRepo,
		fileDriver,
	)

	return &Container{
		config:                       cfg,
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
		showCheckoutUseCase:          showCheckoutUseCase,
	}, nil
}

// GetConfig returns the application configuration
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// Repository getters
func (c *Container) GetOffersRepository() repositories.OffersRepository {
	return c.offersRepo
}

func (c *Container) GetProductsRepository() repositories.ProductsRepository {
	return c.productsRepo
}

func (c *Container) GetUsersRepository() repositories.UsersRepository {
	return c.usersRepo
}

func (c *Container) GetCompaniesRepository() repositories.CompaniesRepository {
	return c.companiesRepo
}

func (c *Container) GetFormatsRepository() repositories.FormatsRepository {
	return c.formatsRepo
}

func (c *Container) GetCheckoutConfigsRepository() repositories.CheckoutConfigsRepository {
	return c.checkoutConfigsRepo
}

func (c *Container) GetAffiliatesRepository() repositories.AffiliatesRepository {
	return c.affiliatesRepo
}

func (c *Container) GetProductAffiliateSettingsRepository() repositories.ProductAffiliateSettingsRepository {
	return c.productAffiliateSettingsRepo
}

func (c *Container) GetCheckoutsRepository() repositories.CheckoutsRepository {
	return c.checkoutsRepo
}

func (c *Container) GetOrderBumpsRepository() repositories.OrderBumpsRepository {
	return c.orderBumpsRepo
}

func (c *Container) GetReviewsRepository() repositories.ReviewsRepository {
	return c.reviewsRepo
}

func (c *Container) GetPixelsRepository() repositories.PixelsRepository {
	return c.pixelsRepo
}

func (c *Container) GetPlansRepository() repositories.PlansRepository {
	return c.plansRepo
}

func (c *Container) GetDiscountsRepository() repositories.DiscountsRepository {
	return c.discountsRepo
}

func (c *Container) GetFileDriver() repositories.FileDriver {
	return c.fileDriver
}

// Use case getters
func (c *Container) GetShowCheckoutUseCase() *showcheckout.UseCase {
	return c.showCheckoutUseCase
}
