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
	ID        string    `json:"id"`
	LockedBy  string    `json:"locked_by"`
	Quantity  int64     `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ItemsLockCompletedEvent ...
type ItemsLockCompletedEvent struct {
	Items []*ItemLockCompletedEvent `json:"items"`
}

// ItemCreatedEvent ...
type ItemCreatedEvent struct {
	ID            string    `json:"id"`
	OwnerID       string    `json:"owner_id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description"`
	TotalQuantity int64     `json:"total_quantity"`
	CreatedAt     time.Time `json:"created_at"`
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

// TradeOfferAcceptedItemEvent ...
type TradeOfferAcceptedItemEvent struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// TradeOfferAcceptedEvent ...
type TradeOfferAcceptedEvent struct {
	ID           string                         `json:"id"`
	OwnerID      string                         `json:"owner_id"`
	OfferedItems []*TradeOfferAcceptedItemEvent `json:"offered_items"`
	WantedItems  []*TradeOfferAcceptedItemEvent `json:"wanted_items"`
}

// ItemsTradeCompletedEvent ...
type ItemsTradeCompletedEvent struct {
	ID string `json:"id"`
}

// ParseItemsToItemsLockCompletedEvent ...
func ParseItemsToItemsLockCompletedEvent(s []*Item) *ItemsLockCompletedEvent {

	items := []*ItemLockCompletedEvent{}

	for _, item := range s {
		for _, lock := range item.Locks {
			items = append(items, &ItemLockCompletedEvent{
				ID:        item.ID,
				LockedBy:  lock.LockedBy,
				Quantity:  int64(lock.Quantity),
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	return &ItemsLockCompletedEvent{Items: items}
}

// ParseItemsToItemsCreatedEvent ...
func ParseItemsToItemsCreatedEvent(s []*Item) *ItemsCreatedEvent {

	items := make([]*ItemCreatedEvent, len(s))

	for i, item := range s {
		items[i] = &ItemCreatedEvent{
			ID:            item.ID,
			OwnerID:       item.OwnerID,
			Name:          string(item.Name),
			Description:   (*string)(item.Description),
			TotalQuantity: int64(item.TotalQuantity),
			CreatedAt:     item.CreatedAt,
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
