package inventory

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type service struct {
	repository Repository
	pool       *pgxpool.Pool
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

	for i, it := range req.Items {
		item, err := NewItem(
			uuid.NewString(),
			userID,
			it.Name,
			it.Description,
			it.Quantity,
			ItemPendingCreateDispatch,
		)

		if err != nil {
			return err
		}

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

	itemsToUpdate := make(map[string]*UpdateItemModel, len(req.Items))
	ids := make([]string, len(req.Items))

	for i, item := range req.Items {
		ids[i] = item.ID
		itemsToUpdate[item.ID] = item
	}

	items, err := s.repository.Get(ctx, userID, ids)

	if err != nil {
		return err
	}

	for _, item := range items {
		itemToUpdate := itemsToUpdate[item.ID]

		err := item.Update(itemToUpdate.Name, itemToUpdate.Description, itemToUpdate.Quantity)

		if err != nil {
			return err
		}
	}

	if err := s.repository.UpdateBulk(ctx, &userID, items); err != nil {
		return err
	}

	return nil
}

// LockItems ...
func (s *service) LockItems(ctx context.Context, userID, correlationID string, req *LockItemsRequest) error {

	itemsToLock := make(map[string]*LockItemModel, len(req.Items))
	ids := make([]string, len(req.Items))

	for i, item := range req.Items {
		ids[i] = item.ID
		itemsToLock[item.ID] = item
	}

	items, err := s.repository.Get(ctx, userID, ids)

	if err != nil {
		return err
	}

	for _, item := range items {
		itemToUpdate := itemsToLock[item.ID]

		if err := item.Lock(itemToUpdate.Quantity); err != nil {
			return err
		}
	}

	if err := s.repository.UpdateBulk(ctx, &userID, items); err != nil {
		return err
	}

	return nil
}

// DeleteItems ...
func (s *service) DeleteItems(ctx context.Context, userID, correlationID string, req *DeleteItemsRequest) error {

	err := s.repository.DeleteBulk(ctx, userID, req.IDs)

	if err != nil {
		return err
	}

	return nil
}
