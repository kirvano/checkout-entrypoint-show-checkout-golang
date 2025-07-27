package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"checkout-go/internal/config"
	"checkout-go/internal/core/entities"
	"checkout-go/internal/infrastructure/aws"
	"checkout-go/internal/repositories"
)

// Base repository functionality

// BaseRepository provides common DynamoDB operations
type BaseRepository struct {
	client *Client
}

func NewBaseRepository(client *Client) *BaseRepository {
	return &BaseRepository{client: client}
}

// OffersRepository implementation
type OffersRepository struct {
	*BaseRepository
	tableName string
}

func NewOffersRepository(client *Client, cfg *config.Config) *OffersRepository {
	return &OffersRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "offers"),
	}
}

func (r *OffersRepository) FindByUUID(ctx context.Context, uuid string) (*repositories.Offer, error) {
	// Query the UuidIndex GSI instead of using GetItem on primary key
	input := &dynamodb.QueryInput{
		TableName: &r.tableName,
		IndexName: stringPtr("UuidIndex"),
		KeyConditionExpression: stringPtr("#uuid = :uuid"),
		ExpressionAttributeNames: map[string]string{
			"#uuid": "uuid",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uuid": &types.AttributeValueMemberS{Value: uuid},
		},
		Limit: int32Ptr(1), // We only expect one result
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer by UUID: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var offer repositories.Offer
	if err := attributevalue.UnmarshalMap(result.Items[0], &offer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal offer: %w", err)
	}

	return &offer, nil
}

func (r *OffersRepository) Find(ctx context.Context, id int) (*repositories.Offer, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	}

	result, err := r.client.GetDynamoDB().GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer by ID: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var offer repositories.Offer
	if err := attributevalue.UnmarshalMap(result.Item, &offer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal offer: %w", err)
	}

	return &offer, nil
}

func (r *OffersRepository) IncrementCheckoutCount(ctx context.Context, uuid string) error {
	// First, find the offer by UUID to get its id (primary key)
	offer, err := r.FindByUUID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("failed to find offer by UUID: %w", err)
	}
	if offer == nil {
		return fmt.Errorf("offer not found")
	}

	// Now update using the primary key (id)
	input := &dynamodb.UpdateItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", offer.ID)},
		},
		UpdateExpression: stringPtr("ADD checkout_count :inc"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	}

	_, err = r.client.GetDynamoDB().UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to increment checkout count: %w", err)
	}

	return nil
}

// CheckoutsRepository implementation
type CheckoutsRepository struct {
	*BaseRepository
	tableName string
}

func NewCheckoutsRepository(client *Client, cfg *config.Config) *CheckoutsRepository {
	return &CheckoutsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "checkouts"),
	}
}

func (r *CheckoutsRepository) Create(ctx context.Context, checkout *entities.Checkout) error {
	item, err := attributevalue.MarshalMap(checkout)
	if err != nil {
		return fmt.Errorf("failed to marshal checkout: %w", err)
	}

	// Fix case sensitivity issue: rename UUID to uuid if it exists
	if uuidAttr, exists := item["UUID"]; exists {
		item["uuid"] = uuidAttr
		delete(item, "UUID")
		fmt.Printf("DEBUG: Fixed UUID case - moved from 'UUID' to 'uuid'\n")
	}

	// Debug: Check if uuid key exists in marshaled item
	if uuidAttr, exists := item["uuid"]; exists {
		fmt.Printf("DEBUG: Marshaled item has uuid key with value: %v\n", uuidAttr)
	} else {
		fmt.Printf("DEBUG: ERROR - Marshaled item is still missing uuid key! Available keys: %v\n", getMapKeys(item))
	}

	input := &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	}

	_, err = r.client.GetDynamoDB().PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create checkout: %w", err)
	}

	return nil
}

// Helper function to get map keys for debugging
func getMapKeys(m map[string]types.AttributeValue) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (r *CheckoutsRepository) FindByUUID(ctx context.Context, uuid string) (*entities.Checkout, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"uuid": &types.AttributeValueMemberS{Value: uuid},
		},
	}

	result, err := r.client.GetDynamoDB().GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkout by UUID: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var checkout entities.Checkout
	if err := attributevalue.UnmarshalMap(result.Item, &checkout); err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkout: %w", err)
	}

	return &checkout, nil
}

func (r *CheckoutsRepository) Update(ctx context.Context, checkout *entities.Checkout) error {
	checkout.UpdateTimestamp()

	item, err := attributevalue.MarshalMap(checkout)
	if err != nil {
		return fmt.Errorf("failed to marshal checkout: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	}

	_, err = r.client.GetDynamoDB().PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update checkout: %w", err)
	}

	return nil
}

// Simplified implementations for other repositories
// In a real application, you would implement all of these fully

func NewProductsRepository(client *Client, cfg *config.Config) repositories.ProductsRepository {
	return &ProductsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "products"),
	}
}

func NewUsersRepository(client *Client, cfg *config.Config) repositories.UsersRepository {
	return &UsersRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "users"),
	}
}

func NewCompaniesRepository(client *Client, cfg *config.Config) repositories.CompaniesRepository {
	return &CompaniesRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "companies"),
	}
}

func NewFormatsRepository(client *Client, cfg *config.Config) repositories.FormatsRepository {
	return &FormatsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "formats"),
	}
}

func NewCheckoutConfigsRepository(client *Client, cfg *config.Config) repositories.CheckoutConfigsRepository {
	return &CheckoutConfigsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "checkout_configs"),
	}
}

func NewAffiliatesRepository(client *Client, cfg *config.Config) repositories.AffiliatesRepository {
	return &AffiliatesRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "affiliates"),
	}
}

func NewProductAffiliateSettingsRepository(client *Client, cfg *config.Config) repositories.ProductAffiliateSettingsRepository {
	return &ProductAffiliateSettingsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "product_affiliate_settings"),
	}
}

func NewOrderBumpsRepository(client *Client, cfg *config.Config) repositories.OrderBumpsRepository {
	return &OrderBumpsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "order_bumps"),
	}
}

func NewReviewsRepository(client *Client, cfg *config.Config) repositories.ReviewsRepository {
	return &ReviewsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "reviews"),
	}
}

func NewPixelsRepository(client *Client, cfg *config.Config) repositories.PixelsRepository {
	return &PixelsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "pixels"),
	}
}

func NewPlansRepository(client *Client, cfg *config.Config) repositories.PlansRepository {
	return &PlansRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "plans"),
	}
}

func NewDiscountsRepository(client *Client, cfg *config.Config) repositories.DiscountsRepository {
	return &DiscountsRepository{
		BaseRepository: NewBaseRepository(client),
		tableName:      aws.GetTableName(cfg, "discounts"),
	}
}

// Stub implementations for simplified repositories
type ProductsRepository struct {
	*BaseRepository
	tableName string
}

func (r *ProductsRepository) Find(ctx context.Context, id int) (*repositories.Product, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	}

	result, err := r.client.GetDynamoDB().GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var product repositories.Product
	if err := attributevalue.UnmarshalMap(result.Item, &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return &product, nil
}

type UsersRepository struct {
	*BaseRepository
	tableName string
}

func (r *UsersRepository) Find(ctx context.Context, id int) (*repositories.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	}

	result, err := r.client.GetDynamoDB().GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var user repositories.User
	if err := attributevalue.UnmarshalMap(result.Item, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

type CompaniesRepository struct {
	*BaseRepository
	tableName string
}

func (r *CompaniesRepository) Find(ctx context.Context, id int) (*repositories.Company, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	}

	result, err := r.client.GetDynamoDB().GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get company by ID: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var company repositories.Company
	if err := attributevalue.UnmarshalMap(result.Item, &company); err != nil {
		return nil, fmt.Errorf("failed to unmarshal company: %w", err)
	}

	return &company, nil
}

type FormatsRepository struct {
	*BaseRepository
	tableName string
}

func (r *FormatsRepository) Find(ctx context.Context, id int) (*repositories.Format, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	}

	result, err := r.client.GetDynamoDB().GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get format by ID: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var format repositories.Format
	if err := attributevalue.UnmarshalMap(result.Item, &format); err != nil {
		return nil, fmt.Errorf("failed to unmarshal format: %w", err)
	}

	return &format, nil
}

type CheckoutConfigsRepository struct {
	*BaseRepository
	tableName string
}

func (r *CheckoutConfigsRepository) Find(ctx context.Context, id int) (*repositories.CheckoutConfig, error) {
	input := &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)},
		},
	}

	result, err := r.client.GetDynamoDB().GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get checkout config by ID: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var checkoutConfig repositories.CheckoutConfig
	if err := attributevalue.UnmarshalMap(result.Item, &checkoutConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkout config: %w", err)
	}

	return &checkoutConfig, nil
}

type AffiliatesRepository struct {
	*BaseRepository
	tableName string
}

func (r *AffiliatesRepository) FindByUUID(ctx context.Context, uuid string) (*repositories.Affiliate, error) {
	return &repositories.Affiliate{}, nil // Simplified for demo
}

type ProductAffiliateSettingsRepository struct {
	*BaseRepository
	tableName string
}

func (r *ProductAffiliateSettingsRepository) FindByProduct(ctx context.Context, productID int) (*repositories.ProductAffiliateSettings, error) {
	return &repositories.ProductAffiliateSettings{}, nil // Simplified for demo
}

type OrderBumpsRepository struct {
	*BaseRepository
	tableName string
}

func (r *OrderBumpsRepository) FindAllByOffer(ctx context.Context, offerID int) ([]*repositories.OrderBump, error) {
	return []*repositories.OrderBump{}, nil // Simplified for demo
}

type ReviewsRepository struct {
	*BaseRepository
	tableName string
}

func (r *ReviewsRepository) FindByCheckoutConfig(ctx context.Context, checkoutConfigID int) ([]*repositories.Review, error) {
	return []*repositories.Review{}, nil // Simplified for demo
}

type PixelsRepository struct {
	*BaseRepository
	tableName string
}

func (r *PixelsRepository) FindAllByUserAndProduct(ctx context.Context, userID, productID int) ([]*repositories.Pixel, error) {
	return []*repositories.Pixel{}, nil // Simplified for demo
}

type PlansRepository struct {
	*BaseRepository
	tableName string
}

func (r *PlansRepository) FindByOffer(ctx context.Context, offerID int) ([]*repositories.Plan, error) {
	return []*repositories.Plan{}, nil // Simplified for demo
}

type DiscountsRepository struct {
	*BaseRepository
	tableName string
}

func (r *DiscountsRepository) CheckHasDiscounts(ctx context.Context, productID int) (bool, error) {
	return false, nil // Simplified for demo
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
