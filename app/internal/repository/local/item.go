package local

import (
	"context"
	"sync"

	"github.com/TrevorEdris/hypernova-bot/app/model/item"
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
func (r *ItemRepo) Get(ctx context.Context, id string) (item.Model, error) {
	tmp, exists := r.storage.Load(id)
	it := tmp.(item.Model)
	if !exists {
		return item.Model{}, item.ErrItemNotFound
	}
	return it, nil
}

// Create creates a new item with the properties of the given item model.
func (r *ItemRepo) Create(ctx context.Context, it item.Model) (item.Model, error) {
	id := "some_unique_id"
	it.ID = id
	r.storage.Store(id, it)
	return it, nil
}

// Create updates the fields of the item identified by id to match the fields of the given item model.
func (r *ItemRepo) Update(ctx context.Context, id string, updates item.Model) (item.Model, error) {
	tmp, exists := r.storage.Load(id)
	it := tmp.(item.Model)
	if !exists {
		return item.Model{}, item.ErrItemNotFound
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
		return item.ErrItemNotFound
	}
	return nil
}
