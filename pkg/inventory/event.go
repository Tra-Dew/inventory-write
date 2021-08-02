package inventory

import "time"

// ItemUpdatedEvent ...
type ItemUpdatedEvent struct {
	ID             string    `json:"id"`
	OwnerID        string    `json:"owner_id"`
	Name           string    `json:"name"`
	Description    *string   `json:"description"`
	TotalQuantity  int64     `json:"total_quantity"`
	LockedQuantity int64     `json:"locked_quantity"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ItemsUpdatedEvent ...
type ItemsUpdatedEvent struct {
	Items []*ItemUpdatedEvent `json:"items"`
}

// ParseItemsToItemsUpdatedEvent ...
func ParseItemsToItemsUpdatedEvent(s []*Item) *ItemsUpdatedEvent {

	items := make([]*ItemUpdatedEvent, len(s))

	for i, item := range s {

		var lockedQuantity int64
		for _, lock := range item.Locks {
			lockedQuantity = lockedQuantity + int64(lock.Quantity)
		}

		items[i] = &ItemUpdatedEvent{
			ID:             item.ID,
			OwnerID:        item.OwnerID,
			Name:           string(item.Name),
			Description:    (*string)(item.Description),
			TotalQuantity:  int64(item.TotalQuantity),
			LockedQuantity: lockedQuantity,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		}
	}

	return &ItemsUpdatedEvent{Items: items}
}
