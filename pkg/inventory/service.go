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

	if err := s.repository.UpdateBulk(ctx, items); err != nil {
		return err
	}

	return nil
}

// LockItems ...
func (s *service) LockItems(ctx context.Context, req *LockItemsRequest) error {

	wantedIDs := make([]string, len(req.WantedItems))
	for i, item := range req.WantedItems {
		wantedIDs[i] = item.ID
	}

	// getting wanted items and validating whether
	// all items exists and they all belong to the same user
	wantedItems, err := s.repository.Get(ctx, &req.WantedItemsOwnerID, wantedIDs)

	if err != nil {
		return err
	}

	if len(wantedItems) != len(req.WantedItems) {
		return core.ErrInvalidWantedItems
	}

	ids := make([]string, len(req.OfferedItems))
	itemsToLock := make(map[string]*LockItemModel, len(req.OfferedItems))

	for i, item := range req.OfferedItems {
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

	if err := s.repository.UpdateBulk(ctx, items); err != nil {
		return err
	}

	return nil
}

// TradeItems ...
func (s *service) TradeItems(ctx context.Context, req *TradeItemsRequest) error {
	// TODO: refactor this method and send this to trade service
	// use the inventory api only as a crud like service

	// Getting all offered items
	offeredIDs := make([]string, len(req.OfferedItems))
	for i, item := range req.OfferedItems {
		offeredIDs[i] = item.ID
	}

	offeredItems, err := s.repository.Get(ctx, &req.OwnerID, offeredIDs)
	if err != nil {
		return err
	}

	// Getting all wanted items
	wantedIDs := make([]string, len(req.WantedItems))
	wantedQuantities := make(map[string]int64, len(req.WantedItems))
	for i, item := range req.WantedItems {
		wantedIDs[i] = item.ID
		wantedQuantities[item.ID] = item.Quantity
	}

	wantedItems, err := s.repository.Get(ctx, &req.WantedItemsOwnerID, wantedIDs)
	if err != nil {
		return err
	}

	var itemsToAdd []*Item
	var itemsToUpdate []*Item
	var itemsToDelete []string

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

		newItem, err := NewItem(
			uuid.NewString(), req.WantedItemsOwnerID, string(item.Name),
			(*string)(item.Description), int64(offeredQuantity), ItemPendingUpdateDispatch,
		)

		if err != nil {
			return err
		}

		itemsToAdd = append(itemsToAdd, newItem)

		if item.TotalQuantity > 0 {
			itemsToUpdate = append(itemsToUpdate, item)
		} else {
			itemsToDelete = append(itemsToDelete, item.ID)
		}
	}

	for _, item := range wantedItems {
		wantedQuantity := ItemQuantity(wantedQuantities[item.ID])

		if item.TotalQuantity < wantedQuantity {
			return core.ErrValidationFailed
		}

		item.TotalQuantity = item.TotalQuantity - wantedQuantity

		newItem, err := NewItem(
			uuid.NewString(),
			req.OwnerID,
			string(item.Name),
			(*string)(item.Description),
			int64(wantedQuantity),
			ItemPendingUpdateDispatch,
		)

		if err != nil {
			return err
		}

		itemsToAdd = append(itemsToAdd, newItem)

		if item.TotalQuantity > 0 {
			itemsToUpdate = append(itemsToUpdate, item)
		} else {
			itemsToDelete = append(itemsToDelete, item.ID)
		}
	}

	//TODO: make this operationg a single transaction
	if err := s.repository.UpdateBulk(ctx, itemsToUpdate); err != nil {
		return err
	}

	if err := s.repository.InsertBulk(ctx, itemsToAdd); err != nil {
		return err
	}

	if err := s.repository.DeleteBulk(ctx, itemsToDelete); err != nil {
		return err
	}

	return nil
}

// DeleteItems ...
func (s *service) DeleteItems(ctx context.Context, userID, correlationID string, req *DeleteItemsRequest) error {

	items, err := s.repository.Get(ctx, &userID, req.IDs)

	if err != nil {
		return err
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}

	if err := s.repository.DeleteBulk(ctx, ids); err != nil {
		return err
	}

	return nil
}
