package repository

import (
	"context"

	"github.com/TrevorEdris/hypernova-bot/api/config"
	"github.com/TrevorEdris/hypernova-bot/api/domain"
	"github.com/TrevorEdris/hypernova-bot/api/internal/repository/dynamodb"
	"github.com/TrevorEdris/hypernova-bot/api/internal/repository/local"
)

type (
	// ItemRepo defines the interface by which services can interact with a storage medium
	// that stores the model for an Item.
	ItemRepo interface {
		Get(ctx context.Context, id string) (domain.Item, error)
		Create(ctx context.Context, it domain.Item) (domain.Item, error)
		Update(ctx context.Context, id string, updates domain.Item) (domain.Item, error)
		Delete(ctx context.Context, id string) error
	}
)

func NewItemRepoLocal() *local.ItemRepo {
	return local.NewItemRepo()
}

func NewItemRepoDynamoDB(cfg *config.Config, driver dynamodb.DynamodbClient) *dynamodb.ItemRepo {
	return dynamodb.NewItemRepo(cfg, driver)
}
