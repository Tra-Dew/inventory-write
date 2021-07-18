package memory

import (
	"context"

	"github.com/Tra-Dew/inventory-write/pkg/inventory"
)

type repositoryInMemory struct {
	data map[string]*inventory.Item
}

// NewRepository ...
func NewRepository() inventory.Repository {
	return &repositoryInMemory{
		data: make(map[string]*inventory.Item),
	}
}

// InsertBulk ...
func (r *repositoryInMemory) InsertBulk(ctx context.Context, items []*inventory.Item) error {
	for _, item := range items {
		r.data[item.ID] = item
	}

	return nil
}

// UpdateBulk ...
func (r *repositoryInMemory) UpdateBulk(ctx context.Context, userID string, items []*inventory.UpdateItem) error {
	for _, item := range items {
		currentItem := r.data[item.ID]
		if currentItem != nil && currentItem.OwnerID == userID {
			currentItem.Name = item.Name
			currentItem.Description = item.Description
			currentItem.Quantity = item.Quantity
			currentItem.UpdatedAt = item.UpdatedAt
		}
	}

	return nil
}

// DeleteBulk ...
func (r *repositoryInMemory) DeleteBulk(ctx context.Context, userID string, ids []string) error {
	for _, id := range ids {
		item := r.data[id]
		if item != nil && item.OwnerID == userID {
			delete(r.data, id)
		}
	}

	return nil
}
