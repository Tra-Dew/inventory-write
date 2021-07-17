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

// Item ...
type Item struct {
	ID          string     `bson:"id"`
	OwnerID     string     `bson:"owner_id"`
	Name        string     `bson:"name"`
	Status      ItemStatus `bson:"status"`
	Description *string    `bson:"description"`
	CreatedAt   time.Time  `bson:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at"`
}

// Repository ...
type Repository interface {
	InsertMany(ctx context.Context, items []*Item) error
	DeleteMany(ctx context.Context, userID string, ids []string) error
}

// Service ...
type Service interface {
	CreateItems(ctx context.Context, userID, correlationID string, req *CreateItemsRequest) (*CreateItemsResponse, error)
	DeleteMany(ctx context.Context, userID, correlationID string, req *DeleteItemsRequest) error
}

// NewItem ...
func NewItem(id, ownerID, name string, description *string, status ItemStatus) (*Item, error) {

	if id == "" {
		return nil, core.ErrValidationFailed
	}

	if ownerID == "" {
		return nil, core.ErrValidationFailed
	}

	name = strings.TrimSpace(name)
	if len(name) < 3 {
		return nil, core.ErrValidationFailed
	}

	if description != nil {
		fixDescription := strings.TrimSpace(*description)
		description = &fixDescription
	}

	if status == "" {
		return nil, core.ErrValidationFailed
	}

	return &Item{
		ID:          id,
		OwnerID:     ownerID,
		Name:        name,
		Status:      status,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
