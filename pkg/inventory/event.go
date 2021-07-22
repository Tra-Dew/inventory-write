package inventory

import "time"

// ItemLockRequestedEvent ...
type ItemLockRequestedEvent struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// ItemsLockRequestedEvent ...
type ItemsLockRequestedEvent struct {
	Items         []*ItemLockRequestedEvent `json:"items"`
	OwnerID       string                    `json:"owner_id"`
	CorrelationID string                    `json:"correlation_id"`
}

// ItemLockCompletedEvent ...
type ItemLockCompletedEvent struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// ItemsLockCompletedEvent ...
type ItemsLockCompletedEvent struct {
	Items []*ItemLockCompletedEvent `json:"items"`
}

// ItemCreatedEvent ...
type ItemCreatedEvent struct {
	ID             string    `json:"id"`
	OwnerID        string    `json:"owner_id"`
	Name           string    `json:"name"`
	Description    *string   `json:"description"`
	TotalQuantity  int64     `json:"total_quantity"`
	LockedQuantity int64     `json:"locked_quantity"`
	CreatedAt      time.Time `json:"created_at"`
}

// ItemsCreatedEvent ...
type ItemsCreatedEvent struct {
	Items []*ItemCreatedEvent `json:"items"`
}

// ItemUpdatedEvent ...
type ItemUpdatedEvent struct {
	ID            string    `json:"id"`
	OwnerID       string    `json:"owner_id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description"`
	TotalQuantity int64     `json:"total_quantity"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ItemsUpdatedEvent ...
type ItemsUpdatedEvent struct {
	Items []*ItemUpdatedEvent `json:"items"`
}

// ParseItemsToItemsLockCompletedEvent ...
func ParseItemsToItemsLockCompletedEvent(s []*Item) *ItemsLockCompletedEvent {

	items := make([]*ItemLockCompletedEvent, len(s))

	for i, item := range s {
		items[i] = &ItemLockCompletedEvent{
			ID:       item.ID,
			Quantity: int64(item.LockedQuantity),
		}
	}

	return &ItemsLockCompletedEvent{Items: items}
}

// ParseItemsToItemsCreatedEvent ...
func ParseItemsToItemsCreatedEvent(s []*Item) *ItemsCreatedEvent {

	items := make([]*ItemCreatedEvent, len(s))

	for i, item := range s {
		items[i] = &ItemCreatedEvent{
			ID:             item.ID,
			OwnerID:        item.OwnerID,
			Name:           string(item.Name),
			Description:    (*string)(item.Description),
			TotalQuantity:  int64(item.TotalQuantity),
			LockedQuantity: int64(item.LockedQuantity),
			CreatedAt:      item.CreatedAt,
		}
	}

	return &ItemsCreatedEvent{Items: items}
}

// ParseItemsToItemsUpdatedEvent ...
func ParseItemsToItemsUpdatedEvent(s []*Item) *ItemsUpdatedEvent {

	items := make([]*ItemUpdatedEvent, len(s))

	for i, item := range s {
		items[i] = &ItemUpdatedEvent{
			ID:            item.ID,
			OwnerID:       item.OwnerID,
			Name:          string(item.Name),
			Description:   (*string)(item.Description),
			TotalQuantity: int64(item.TotalQuantity),
			UpdatedAt:     item.UpdatedAt,
		}
	}

	return &ItemsUpdatedEvent{Items: items}
}
