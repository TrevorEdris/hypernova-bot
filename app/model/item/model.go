package item

import (
	"errors"
	"fmt"
)

var (
	// ErrItemNotFound is the error returned when the specified item cannot be found.
	ErrItemNotFound = errors.New("item not found")
)

type (
	// Model defines the domain model of an Item. This is the struct that will be
	// used for all business logic. The struct has no annotations because the
	// way the model is rendered to some medium such as an HTTP response, database,
	// queue, etc. should be completely independent of the actual struct used by
	// the business logic functions. The way the item is rendered should be left
	// up to the package performing the rendering.
	Model struct {
		ID              string
		Name            string
		Description     string
		Price           float64
		calculatedField string
	}
)

// New creates a new Model.
func New(id, name, description string, price float64) Model {
	return Model{
		ID:              id,
		Name:            name,
		Description:     description,
		Price:           price,
		calculatedField: fmt.Sprintf("%02f:%s", price, id),
	}
}

func (m *Model) Update(updates Model) {
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
