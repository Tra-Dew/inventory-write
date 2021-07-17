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
func (s *service) CreateItems(ctx context.Context, userID, correlationID string, req *CreateItemsRequest) (*CreateItemsResponse, error) {
	items := make([]*Item, req.Quantity)
	ids := make([]string, req.Quantity)

	for i := 0; i < int(req.Quantity); i++ {
		item, err := NewItem(uuid.NewString(), userID, req.Name, req.Description, ItemAvailable)
		if err != nil {
			return nil, err
		}

		ids[i] = item.ID
		items[i] = item
	}

	if err := s.repository.InsertMany(ctx, items); err != nil {
		return nil, err
	}

	return &CreateItemsResponse{IDs: ids}, nil
}

//TODO: make this entire operation asyncronous, since mongo bulk delete is not atomic
// DeleteMany ...
func (s *service) DeleteMany(ctx context.Context, userID, correlationID string, req *DeleteItemsRequest) error {

	if err := s.repository.DeleteMany(ctx, userID, req.IDs); err != nil {
		return err
	}

	return nil
}
