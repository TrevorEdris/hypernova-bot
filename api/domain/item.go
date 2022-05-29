package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrItemNotFound is the error returned when the specified item cannot be found.
	ErrItemNotFound = errors.New("item not found")
)

type (
	// Item defines the domain model of an Item. This is the struct that will be
	// used for all business logic. The struct has no annotations because the
	// way the model is rendered to some medium such as an HTTP response, database,
	// queue, etc. should be completely independent of the actual struct used by
	// the business logic functions. The way the item is rendered should be left
	// up to the package performing the rendering.
	Item struct {
		ID              string
		Name            string
		Description     string
		Price           float64
		calculatedField string
	}
)

// NewItem creates a new Item model.
func NewItem(id, name, description string, price float64) Item {
	return Item{
		ID:              id,
		Name:            name,
		Description:     description,
		Price:           price,
		calculatedField: fmt.Sprintf("%02f:%s", price, id),
	}
}

func (m *Item) Update(updates Item) {
	if updates.Name != "" {
		m.Name = updates.Name
	}
	if updates.Description != "" {
		m.Description = updates.Description
	}
	if updates.Price >= 0 {
		m.Price = updates.Price
		m.calculatedField = fmt.Sprintf("%02f:%s", m.Price, m.ID)
	}
}
