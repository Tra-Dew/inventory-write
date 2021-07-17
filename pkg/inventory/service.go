package inventory

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	repository Repository
}

// NewService ...
func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

//TODO: make this entire operation asyncronous, since mongo bulk write is not atomic
// CreateItems ...
func (s *service) CreateItems(ctx context.Context, userID, correlationID string, req *CreateItemsRequest) error {
	items := make([]*Item, len(req.Items))
	ids := make([]string, len(req.Items))

	for i, it := range req.Items {
		item, err := NewItem(uuid.NewString(), userID, it.Name, it.Description, it.Quantity, ItemAvailable)

		if err != nil {
			return err
		}

		ids[i] = item.ID
		items[i] = item
	}

	if err := s.repository.InsertBulk(ctx, items); err != nil {
		return err
	}

	return nil
}

//TODO: make this entire operation asyncronous, since mongo bulk delete is not atomic
// UpdateItems ...
func (s *service) UpdateItems(ctx context.Context, userID, correlationID string, req *UpdateItemsRequest) error {

	var itemsToUpdate []*UpdateItem
	var itemsToDelete []string

	for _, item := range req.Items {

		if item.Quantity == 0 {
			itemsToDelete = append(itemsToDelete, item.ID)
		} else {

			updatedItem, err := NewUpdateItem(item.ID, item.Name, item.Description, item.Quantity)
			if err != nil {
				return err
			}

			itemsToUpdate = append(itemsToUpdate, updatedItem)
		}
	}

	if err := s.repository.UpdateBulk(ctx, userID, itemsToUpdate); err != nil {
		return err
	}

	if err := s.repository.DeleteBulk(ctx, userID, itemsToDelete); err != nil {
		return err
	}

	return nil
}
