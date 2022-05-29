package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/TrevorEdris/hypernova-bot/app/config"
	"github.com/TrevorEdris/hypernova-bot/app/model/item"
)

type (
	ItemRepo struct {
		storage DynamodbClient
		table   string
	}

	ddbItemModel struct {
		Id          string  `dynamodbav:"id"`
		Name        string  `dynamodbav:"title"`
		Description string  `dynamodbav:"description"`
		Price       float64 `dynamodbav:"price"`
	}
)

func NewItemRepo(cfg *config.Config, driver DynamodbClient) *ItemRepo {
	return &ItemRepo{
		storage: driver,
		table:   cfg.DynamoDB.ItemTable,
	}
}

// Get retrieves the item identified by the specified id.
func (r *ItemRepo) Get(ctx context.Context, id string) (item.Model, error) {
	params := &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
		TableName: aws.String(r.table),
	}
	result, err := r.storage.GetItem(ctx, params)
	if err != nil {
		return item.Model{}, fmt.Errorf("failed to GetItem: %w", err)
	}
	if result.Item == nil {
		return item.Model{}, item.ErrItemNotFound
	}
	var rec ddbItemModel
	if err = attributevalue.UnmarshalMap(result.Item, &rec); err != nil {
		return item.Model{}, fmt.Errorf("failed to UnmarshalMap: %w", err)
	}
	return r.ddbToModel(rec), nil
}

// Create creates a new item with the properties of the given item model.
func (r *ItemRepo) Create(ctx context.Context, it item.Model) (item.Model, error) {
	it.ID = "some_unique_key"
	marshaled, err := attributevalue.MarshalMap(r.modelToDDB(it))
	if err != nil {
		return item.Model{}, fmt.Errorf("failed to Marshal item: %w", err)
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String(r.table),
		Item:      marshaled,
	}
	_, err = r.storage.PutItem(ctx, params)
	if err != nil {
		return item.Model{}, fmt.Errorf("failed to PutItem: %w", err)
	}
	return it, nil
}

// Create updates the fields of the item identified by id to match the fields of the given item model.
func (r *ItemRepo) Update(ctx context.Context, id string, updates item.Model) (item.Model, error) {
	baseItem, err := r.Get(ctx, id)
	if err != nil {
		return item.Model{}, fmt.Errorf("failed to retrieve item for updates: %w", err)
	}
	baseItem.Update(updates)

	params := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":title":       &types.AttributeValueMemberS{Value: baseItem.Name},
			":description": &types.AttributeValueMemberS{Value: baseItem.Description},
			":price":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", baseItem.Price)},
		},
		UpdateExpression: aws.String("SET title = :title, description = :description, price = :price"),
	}
	_, err = r.storage.UpdateItem(ctx, params)
	if err != nil {
		return item.Model{}, fmt.Errorf("failed to UpdateItem: %w", err)
	}
	return updates, nil
}

// Delete removes the item identified by the specified id.
func (r *ItemRepo) Delete(ctx context.Context, id string) error {
	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.table),
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: id}},
	}
	_, err := r.storage.DeleteItem(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (r *ItemRepo) ddbToModel(rec ddbItemModel) item.Model {
	return item.New(rec.Id, rec.Name, rec.Description, rec.Price)
}

func (r *ItemRepo) modelToDDB(it item.Model) ddbItemModel {
	return ddbItemModel{
		Id:          it.ID,
		Name:        it.Name,
		Description: it.Description,
		Price:       it.Price,
	}
}
