package inventory

import (
	"context"

	"github.com/d-leme/tradew-inventory-write/pkg/core"
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

	items, err := s.repository.Get(ctx, &userID, ids)

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
func (s *service) LockItems(ctx context.Context, req *LockItemsRequest) error {

	itemsToLock := make(map[string]*LockItemModel, len(req.Items))
	ids := make([]string, len(req.Items))

	for i, item := range req.Items {
		ids[i] = item.ID
		itemsToLock[item.ID] = item
	}

	items, err := s.repository.Get(ctx, &req.OwnerID, ids)

	if err != nil {
		return err
	}

	for _, item := range items {
		itemToUpdate := itemsToLock[item.ID]

		if err := item.Lock(req.LockedBy, itemToUpdate.Quantity); err != nil {
			return err
		}
	}

	if err := s.repository.UpdateBulk(ctx, &req.OwnerID, items); err != nil {
		return err
	}

	return nil
}

// TradeItems ...
func (s *service) TradeItems(ctx context.Context, req *TradeItemsRequest) error {

	offeredIDs := make([]string, len(req.OfferedItems))
	for i, item := range req.OfferedItems {
		offeredIDs[i] = item.ID
	}

	wantedQuantities := make(map[string]int64, len(req.WantedItems))
	wantedIDs := make([]string, len(req.WantedItems))
	for i, item := range req.WantedItems {
		wantedIDs[i] = item.ID
		wantedQuantities[item.ID] = item.Quantity
	}

	offeredItems, err := s.repository.Get(ctx, &req.OwnerID, offeredIDs)
	if err != nil {
		return err
	}

	wantedItems, err := s.repository.Get(ctx, &req.WantedItemsOwnerID, wantedIDs)
	if err != nil {
		return err
	}

	var itemsToAdd []*Item

	for _, item := range offeredItems {
		var offeredQuantity ItemQuantity
		newLocks := make([]*ItemLock, len(item.Locks)-1)
		for _, lock := range item.Locks {
			if lock.LockedBy == req.TradeID {
				offeredQuantity = lock.Quantity
			} else {
				newLocks = append(newLocks, lock)
			}
		}

		item.TotalQuantity = item.TotalQuantity - offeredQuantity
		item.Locks = newLocks

		newItem, err := NewItem(uuid.NewString(), wantedItems[0].OwnerID, string(item.Name), (*string)(item.Description), int64(offeredQuantity), ItemPendingUpdateDispatch)
		if err != nil {
			return err
		}

		itemsToAdd = append(itemsToAdd, newItem)
	}

	for _, item := range wantedItems {
		wantedQuantity := ItemQuantity(wantedQuantities[item.ID])

		if item.TotalQuantity < wantedQuantity {
			return core.ErrValidationFailed
		}

		item.TotalQuantity = item.TotalQuantity - wantedQuantity

		newItem, err := NewItem(
			uuid.NewString(),
			offeredItems[0].OwnerID,
			string(item.Name),
			(*string)(item.Description),
			int64(wantedQuantity),
			ItemPendingUpdateDispatch,
		)

		if err != nil {
			return err
		}

		itemsToAdd = append(itemsToAdd, newItem)
	}

	allItems := append(append(offeredItems, wantedItems...), itemsToAdd...)

	if err := s.repository.UpdateBulk(ctx, nil, allItems); err != nil {
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
