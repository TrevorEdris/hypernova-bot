package local

import (
	"context"
	"sync"

	"github.com/TrevorEdris/hypernova-bot/api/domain"
)

type (
	// ItemRepo stores item models in local memory.
	ItemRepo struct {
		storage sync.Map
	}
)

// NewItemRepo creates a new ItemRepo using local memory as the storage medium.
func NewItemRepo() *ItemRepo {
	return &ItemRepo{
		storage: sync.Map{},
	}
}

// Get retrieves the item identified by the specified id.
func (r *ItemRepo) Get(ctx context.Context, id string) (domain.Item, error) {
	tmp, exists := r.storage.Load(id)
	it := tmp.(domain.Item)
	if !exists {
		return domain.Item{}, domain.ErrItemNotFound
	}
	return it, nil
}

// Create creates a new item with the properties of the given item model.
func (r *ItemRepo) Create(ctx context.Context, it domain.Item) (domain.Item, error) {
	id := "some_unique_id"
	it.ID = id
	r.storage.Store(id, it)
	return it, nil
}

// Create updates the fields of the item identified by id to match the fields of the given item model.
func (r *ItemRepo) Update(ctx context.Context, id string, updates domain.Item) (domain.Item, error) {
	tmp, exists := r.storage.Load(id)
	it := tmp.(domain.Item)
	if !exists {
		return domain.Item{}, domain.ErrItemNotFound
	}

	// Apply updates to the item here
	if updates.Name != "" {
		it.Name = updates.Name
	}
	if updates.Description != "" {
		it.Description = updates.Description
	}
	if updates.Price >= float64(0) {
		it.Price = updates.Price
	}

	r.storage.Store(id, it)
	return it, nil
}

// Delete removes the item identified by the specified id.
func (r *ItemRepo) Delete(ctx context.Context, id string) error {
	_, exists := r.storage.LoadAndDelete(id)
	if !exists {
		return domain.ErrItemNotFound
	}
	return nil
}
