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
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("UuidIndex"),
		KeyConditionExpression: stringPtr("#uuid = :uuid"),
		ExpressionAttributeNames: map[string]string{
			"#uuid": "uuid",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uuid": &types.AttributeValueMemberS{Value: uuid},
		},
		Limit: int32Ptr(1),
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get affiliate by UUID: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var affiliate repositories.Affiliate
	if err := attributevalue.UnmarshalMap(result.Items[0], &affiliate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal affiliate: %w", err)
	}

	return &affiliate, nil
}

type ProductAffiliateSettingsRepository struct {
	*BaseRepository
	tableName string
}

func (r *ProductAffiliateSettingsRepository) FindByProduct(ctx context.Context, productID int) (*repositories.ProductAffiliateSettings, error) {
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("ProductIdIndex"),
		KeyConditionExpression: stringPtr("product_id = :product_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":product_id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", productID)},
		},
		Limit: int32Ptr(1),
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get product affiliate settings by product ID: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var settings repositories.ProductAffiliateSettings
	if err := attributevalue.UnmarshalMap(result.Items[0], &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product affiliate settings: %w", err)
	}

	return &settings, nil
}

type OrderBumpsRepository struct {
	*BaseRepository
	tableName string
}

func (r *OrderBumpsRepository) FindAllByOffer(ctx context.Context, offerID int) ([]*repositories.OrderBump, error) {
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("OfferIdIndex"),
		KeyConditionExpression: stringPtr("offer_id = :offer_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":offer_id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", offerID)},
		},
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query order bumps by offer: %w", err)
	}

	var orderBumps []*repositories.OrderBump
	for _, item := range result.Items {
		var orderBump repositories.OrderBump
		if err := attributevalue.UnmarshalMap(item, &orderBump); err != nil {
			return nil, fmt.Errorf("failed to unmarshal order bump: %w", err)
		}
		orderBumps = append(orderBumps, &orderBump)
	}

	return orderBumps, nil
}

type ReviewsRepository struct {
	*BaseRepository
	tableName string
}

func (r *ReviewsRepository) FindByCheckoutConfig(ctx context.Context, checkoutConfigID int) ([]*repositories.Review, error) {
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("checkoutConfigId-index"), // Fixed to match TypeScript
		KeyConditionExpression: stringPtr("#checkoutConfigId = :checkoutConfigId"),
		ExpressionAttributeNames: map[string]string{
			"#checkoutConfigId": "checkoutConfigId", // Fixed to match TypeScript
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":checkoutConfigId": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", checkoutConfigID)},
		},
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query reviews by checkout config: %w", err)
	}

	var reviews []*repositories.Review
	for _, item := range result.Items {
		var review repositories.Review
		if err := attributevalue.UnmarshalMap(item, &review); err != nil {
			return nil, fmt.Errorf("failed to unmarshal review: %w", err)
		}
		
		// Filter by status = ACTIVE (same as TypeScript)
		if review.Status == repositories.ReviewStatusActive {
			reviews = append(reviews, &review)
		}
	}

	return reviews, nil
}

type PixelsRepository struct {
	*BaseRepository
	tableName string
}

func (r *PixelsRepository) FindAllByUserAndProduct(ctx context.Context, userID, productID int) ([]*repositories.Pixel, error) {
	// Use composite key query to match TypeScript implementation
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("productId-userId-index"), // Fixed to match TypeScript
		KeyConditionExpression: stringPtr("#productId = :productId AND #userId = :userId"), // Composite key like TypeScript
		ExpressionAttributeNames: map[string]string{
			"#productId": "productId", // Fixed field names to match TypeScript
			"#userId":    "userId",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":productId": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", productID)},
			":userId":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", userID)},
		},
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query pixels by user and product: %w", err)
	}

	var pixels []*repositories.Pixel
	for _, item := range result.Items {
		var pixel repositories.Pixel
		if err := attributevalue.UnmarshalMap(item, &pixel); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pixel: %w", err)
		}
		pixels = append(pixels, &pixel)
	}

	return pixels, nil
}

type PlansRepository struct {
	*BaseRepository
	tableName string
}

func (r *PlansRepository) FindByUuid(ctx context.Context, uuid string) (*repositories.Plan, error) {
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("UuidIndex"),
		KeyConditionExpression: stringPtr("#uuid = :uuid"),
		ExpressionAttributeNames: map[string]string{
			"#uuid": "uuid",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":uuid": &types.AttributeValueMemberS{Value: uuid},
		},
		Limit: int32Ptr(1),
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query plan by UUID: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var plan repositories.Plan
	if err := attributevalue.UnmarshalMap(result.Items[0], &plan); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
	}

	return &plan, nil
}

func (r *PlansRepository) FindByOffer(ctx context.Context, offerID int) ([]*repositories.Plan, error) {
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("offerId-index"), // Fixed to match TypeScript
		KeyConditionExpression: stringPtr("#offerId = :offerId"),
		ExpressionAttributeNames: map[string]string{
			"#offerId": "offerId", // Fixed to match TypeScript
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":offerId": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", offerID)},
		},
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query plans by offer: %w", err)
	}

	var plans []*repositories.Plan
	for _, item := range result.Items {
		var plan repositories.Plan
		if err := attributevalue.UnmarshalMap(item, &plan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
		}
		plans = append(plans, &plan)
	}

	return plans, nil
}

type DiscountsRepository struct {
	*BaseRepository
	tableName string
}

func (r *DiscountsRepository) CheckHasDiscounts(ctx context.Context, productID int) (bool, error) {
	input := &dynamodb.QueryInput{
		TableName:              &r.tableName,
		IndexName:              stringPtr("ProductIdIndex"),
		KeyConditionExpression: stringPtr("product_id = :product_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":product_id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", productID)},
		},
		Limit: int32Ptr(1), // We only need to know if at least one exists
	}

	result, err := r.client.GetDynamoDB().Query(ctx, input)
	if err != nil {
		return false, fmt.Errorf("failed to query discounts by product: %w", err)
	}

	return len(result.Items) > 0, nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}
