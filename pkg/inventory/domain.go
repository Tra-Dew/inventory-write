package inventory

import (
	"context"
	"strings"
	"time"

	"github.com/d-leme/tradew-inventory-write/pkg/core"
)

// ItemStatus ...
type ItemStatus string

const (
	// ItemAvailable is set when an item is available to be traded
	ItemAvailable ItemStatus = "Available"

	// ItemPendingCreateDispatch is set when an item has been created
	// but has not yet dispached an event
	ItemPendingCreateDispatch ItemStatus = "PendingCreateDispatch"

	// ItemPendingUpdateDispatch is set when an item has been updated
	// but has not yet dispached an event
	ItemPendingUpdateDispatch ItemStatus = "PendingUpdateDispatch"

	// ItemPendingLockDispatch is set when an item has been locked
	// but has not yet dispached an event
	ItemPendingLockDispatch ItemStatus = "PendingLockDispatch"
)

// ItemName ...
type ItemName string

// ItemDescription ...
type ItemDescription string

// ItemQuantity ...
type ItemQuantity int64

// ItemLock ...
type ItemLock struct {
	LockedBy string
	Quantity ItemQuantity
}

// Item ...
type Item struct {
	ID            string
	OwnerID       string
	Name          ItemName
	Status        ItemStatus
	Description   *ItemDescription
	TotalQuantity ItemQuantity
	Locks         []*ItemLock
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Repository ...
type Repository interface {
	InsertBulk(ctx context.Context, items []*Item) error
	UpdateBulk(ctx context.Context, items []*Item) error
	DeleteBulk(ctx context.Context, userID string, ids []string) error
	Get(ctx context.Context, userID *string, ids []string) ([]*Item, error)
	GetByStatus(ctx context.Context, status ItemStatus) ([]*Item, error)
}

// Service ...
type Service interface {
	CreateItems(ctx context.Context, userID, correlationID string, req *CreateItemsRequest) error
	UpdateItems(ctx context.Context, userID, correlationID string, req *UpdateItemsRequest) error
	LockItems(ctx context.Context, req *LockItemsRequest) error
	TradeItems(ctx context.Context, req *TradeItemsRequest) error
	DeleteItems(ctx context.Context, userID, correlationID string, req *DeleteItemsRequest) error
}

// NewItemName ...
func NewItemName(name string) (ItemName, error) {
	name = strings.TrimSpace(name)
	if len(name) < 3 {
		return "", core.ErrValidationFailed
	}

	return ItemName(name), nil
}

// NewItemDescription ...
func NewItemDescription(description *string) *ItemDescription {
	if description == nil || *description == "" {
		return nil
	}

	itemDescription := ItemDescription(strings.TrimSpace(*description))

	return &itemDescription
}

// NewItemQuantity ...
func NewItemQuantity(quantity int64) (ItemQuantity, error) {
	if quantity <= 0 {
		return 0, core.ErrValidationFailed
	}

	return ItemQuantity(quantity), nil
}

// NewItem ...
func NewItem(id, ownerID, name string, description *string, quantity int64, status ItemStatus) (*Item, error) {

	if id == "" {
		return nil, core.ErrValidationFailed
	}

	if ownerID == "" {
		return nil, core.ErrValidationFailed
	}

	itemName, err := NewItemName(name)
	if err != nil {
		return nil, err
	}

	itemQuantity, err := NewItemQuantity(quantity)
	if err != nil {
		return nil, err
	}

	if status == "" {
		return nil, core.ErrValidationFailed
	}

	return &Item{
		ID:            id,
		OwnerID:       ownerID,
		Name:          itemName,
		Status:        status,
		Description:   NewItemDescription(description),
		TotalQuantity: itemQuantity,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// GetLockedQuantity ...
func (item *Item) GetLockedQuantity() ItemQuantity {
	var locksQuantity ItemQuantity
	for _, lock := range item.Locks {
		locksQuantity = locksQuantity + lock.Quantity
	}
	return locksQuantity
}

// Update ...
func (item *Item) Update(name string, description *string, quantity int64) error {
	itemName, err := NewItemName(name)
	if err != nil {
		return err
	}

	itemDescription := NewItemDescription(description)
	if err != nil {
		return err
	}

	itemQuantity, err := NewItemQuantity(quantity)
	if err != nil {
		return err
	}

	if item.GetLockedQuantity() > itemQuantity {
		return core.ErrNotEnoughtItemsToLock
	}

	item.Name = itemName
	item.Description = itemDescription
	item.TotalQuantity = itemQuantity
	item.Status = ItemPendingUpdateDispatch
	item.UpdatedAt = time.Now()

	return nil
}

// UpdateStatus ...
func (item *Item) UpdateStatus(status ItemStatus) {
	item.Status = status
	item.UpdatedAt = time.Now()
}

// Lock ...
func (item *Item) Lock(lockedBy string, quantity int64) error {

	itemQuantity, err := NewItemQuantity(quantity)
	if err != nil {
		return err
	}

	lockedQuantity := item.GetLockedQuantity() + itemQuantity

	if lockedQuantity > item.TotalQuantity {
		return core.ErrNotEnoughtItemsToLock
	}

	item.Locks = append(item.Locks, &ItemLock{
		LockedBy: lockedBy,
		Quantity: itemQuantity,
	})

	item.Status = ItemPendingLockDispatch
	item.UpdatedAt = time.Now()

	return nil
}
