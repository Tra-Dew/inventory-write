package inventory

import (
	"context"

	"github.com/d-leme/tradew-inventory-write/pkg/core"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
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

	fields := logrus.Fields{
		"user_id":        userID,
		"correlation_id": correlationID,
	}

	items := make([]*Item, len(req.Items))

	for i, it := range req.Items {
		item, err := NewItem(
			uuid.NewString(),
			userID,
			it.Name,
			it.Description,
			it.Quantity,
			ItemPendingUpdateDispatch,
		)

		if err != nil {
			logrus.WithError(err).WithFields(fields).Error("error creating new item")
			return err
		}

		items[i] = item
	}

	if err := s.repository.InsertBulk(ctx, items); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while inserting new items")
		return err
	}

	logrus.WithFields(fields).Info("created new items successfully")

	return nil
}

// UpdateItems ...
func (s *service) UpdateItems(ctx context.Context, userID, correlationID string, req *UpdateItemsRequest) error {

	fields := logrus.Fields{
		"user_id":        userID,
		"correlation_id": correlationID,
	}

	itemsToUpdate := make(map[string]*UpdateItemModel, len(req.Items))
	ids := make([]string, len(req.Items))

	for i, item := range req.Items {
		ids[i] = item.ID
		itemsToUpdate[item.ID] = item
	}

	fields["ids"] = ids

	items, err := s.repository.Get(ctx, &userID, ids)

	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while getting items")
		return err
	}

	for _, item := range items {
		itemToUpdate := itemsToUpdate[item.ID]

		err := item.Update(itemToUpdate.Name, itemToUpdate.Description, itemToUpdate.Quantity)

		if err != nil {
			logrus.WithError(err).WithFields(fields).Error("validation error on item")
			return err
		}
	}

	if err := s.repository.UpdateBulk(ctx, items); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while updating items")
		return err
	}

	logrus.WithFields(fields).Info("updated all items succefully")

	return nil
}

// LockItems ...
func (s *service) LockItems(ctx context.Context, req *LockItemsRequest) error {

	fields := logrus.Fields{
		"locked_by":             req.LockedBy,
		"owner_id":              req.OwnerID,
		"wanted_items_owner_id": req.WantedItemsOwnerID,
	}

	wantedIDs := make([]string, len(req.WantedItems))
	for i, item := range req.WantedItems {
		wantedIDs[i] = item.ID
	}

	// getting wanted items and validating whether
	// all items exists and they all belong to the same user
	wantedItems, err := s.repository.Get(ctx, &req.WantedItemsOwnerID, wantedIDs)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while getting wanted items")
		return err
	}

	if len(wantedItems) != len(req.WantedItems) {
		logrus.WithError(core.ErrInvalidWantedItems).WithFields(fields).Error("tried to select invalid wanted items")
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
		logrus.WithError(err).WithFields(fields).Error("error while getting offered items")
		return err
	}

	for _, item := range items {
		itemToUpdate := itemsToLock[item.ID]

		if err := item.Lock(req.LockedBy, itemToUpdate.Quantity); err != nil {
			logrus.WithError(err).WithFields(fields).Error("error item lock failed")
			return err
		}
	}

	if err := s.repository.UpdateBulk(ctx, items); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while updating items")
		return err
	}

	logrus.WithFields(fields).Info("all items updated successfully")

	return nil
}

// TradeItems ...
func (s *service) TradeItems(ctx context.Context, req *TradeItemsRequest) error {
	// TODO: refactor this method and send this to trade service
	// use the inventory api only as a crud like service

	fields := logrus.Fields{
		"locked_by":             req.TradeID,
		"owner_id":              req.OwnerID,
		"wanted_items_owner_id": req.WantedItemsOwnerID,
	}

	// Getting all offered items
	offeredIDs := make([]string, len(req.OfferedItems))
	for i, item := range req.OfferedItems {
		offeredIDs[i] = item.ID
	}

	offeredItems, err := s.repository.Get(ctx, &req.OwnerID, offeredIDs)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while getting offered items")
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
		logrus.WithError(err).WithFields(fields).Error("error while getting wanted items")
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
			logrus.WithError(err).WithFields(fields).Error("error while creating new item to add")
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
			logrus.
				WithError(core.ErrValidationFailed).
				WithFields(fields).
				Error("wanted quantity bigger than total quantity")

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
			logrus.WithError(err).WithFields(fields).Error("error while creating new item to add")
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
		logrus.WithError(err).WithFields(fields).Error("error while updating items")
		return err
	}

	logrus.WithFields(fields).Error("items updated")

	if err := s.repository.InsertBulk(ctx, itemsToAdd); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while inserting items")
		return err
	}

	logrus.WithFields(fields).Error("items inserted")

	if err := s.repository.DeleteBulk(ctx, itemsToDelete); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while deleting items")
		return err
	}

	logrus.WithFields(fields).Error("items deleted")

	return nil
}

// DeleteItems ...
func (s *service) DeleteItems(ctx context.Context, userID, correlationID string, req *DeleteItemsRequest) error {

	fields := logrus.Fields{
		"ids":            req.IDs,
		"owner_id":       userID,
		"correlation_id": correlationID,
	}

	items, err := s.repository.Get(ctx, &userID, req.IDs)
	if err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while getting items")
		return err
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}

	if err := s.repository.DeleteBulk(ctx, ids); err != nil {
		logrus.WithError(err).WithFields(fields).Error("error while deleting items")
		return err
	}

	logrus.WithFields(fields).Info("deleted all items sucessfully")

	return nil
}
