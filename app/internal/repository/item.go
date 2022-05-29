package repository

import (
	"context"

	"github.com/TrevorEdris/hypernova-bot/app/config"
	"github.com/TrevorEdris/hypernova-bot/app/internal/repository/dynamodb"
	"github.com/TrevorEdris/hypernova-bot/app/internal/repository/local"
	"github.com/TrevorEdris/hypernova-bot/app/model/item"
)

type (
	// ItemRepo defines the interface by which services can interact with a storage medium
	// that stores the model for an Item.
	ItemRepo interface {
		Get(ctx context.Context, id string) (item.Model, error)
		Create(ctx context.Context, it item.Model) (item.Model, error)
		Update(ctx context.Context, id string, updates item.Model) (item.Model, error)
		Delete(ctx context.Context, id string) error
	}
)

func NewItemRepoLocal() *local.ItemRepo {
	return local.NewItemRepo()
}

func NewItemRepoDynamoDB(cfg *config.Config, driver dynamodb.DynamodbClient) *dynamodb.ItemRepo {
	return dynamodb.NewItemRepo(cfg, driver)
}
