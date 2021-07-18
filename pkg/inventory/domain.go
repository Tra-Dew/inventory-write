package inventory

import (
	"context"
	"strings"
	"time"

	"github.com/Tra-Dew/inventory-write/pkg/core"
)

// ItemStatus ...
type ItemStatus string

const (
	// ItemAvailable is set when an item is available to be traded
	ItemAvailable ItemStatus = "Available"

	// ItemLocked is set when an item is unavailable and
	// currently in a trade session
	ItemLocked ItemStatus = "Locked"
)

// ItemName ...
type ItemName string

// ItemDescription ...
type ItemDescription string

// ItemQuantity ...
type ItemQuantity int64

// Item ...
type Item struct {
	ID          string           `bson:"id"`
	OwnerID     string           `bson:"owner_id"`
	Name        ItemName         `bson:"name"`
	Status      ItemStatus       `bson:"status"`
	Description *ItemDescription `bson:"description"`
	Quantity    ItemQuantity     `bson:"quantity"`
	CreatedAt   time.Time        `bson:"created_at"`
	UpdatedAt   time.Time        `bson:"updated_at"`
}

// UpdateItem ...
type UpdateItem struct {
	ID          string           `bson:"id"`
	Name        ItemName         `bson:"name"`
	Description *ItemDescription `bson:"description"`
	Quantity    ItemQuantity     `bson:"quantity"`
	UpdatedAt   time.Time        `bson:"updated_at"`
}

// Repository ...
type Repository interface {
	InsertBulk(ctx context.Context, items []*Item) error
	UpdateBulk(ctx context.Context, userID string, items []*UpdateItem) error
	DeleteBulk(ctx context.Context, userID string, ids []string) error
}

// Service ...
type Service interface {
	CreateItems(ctx context.Context, userID, correlationID string, req *CreateItemsRequest) error
	UpdateItems(ctx context.Context, userID, correlationID string, req *UpdateItemsRequest) error
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
		ID:          id,
		OwnerID:     ownerID,
		Name:        itemName,
		Status:      status,
		Description: NewItemDescription(description),
		Quantity:    itemQuantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// NewUpdateItem ...
func NewUpdateItem(id, name string, description *string, quantity int64) (*UpdateItem, error) {

	if id == "" {
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

	return &UpdateItem{
		ID:          id,
		Name:        itemName,
		Description: NewItemDescription(description),
		Quantity:    itemQuantity,
		UpdatedAt:   time.Now(),
	}, nil
}
