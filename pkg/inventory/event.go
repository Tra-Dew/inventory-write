package inventory

import "time"

// TradeItemCreatedEvent ...
type TradeItemCreatedEvent struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// TradeCreatedEvent ...
type TradeCreatedEvent struct {
	ID                 string                   `json:"id"`
	OwnerID            string                   `json:"owner_id"`
	WantedItemsOwnerID string                   `json:"wanted_items_owner_id"`
	OfferedItems       []*TradeItemCreatedEvent `json:"offered_items"`
	WantedItems        []*TradeItemCreatedEvent `json:"wanted_items"`
	CreatedAt          time.Time                `json:"created_at"`
}

// ItemLockCompletedEvent ...
type ItemLockCompletedEvent struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

// ItemsLockCompletedEvent ...
type ItemsLockCompletedEvent struct {
	LockedBy           string                    `json:"locked_by"`
	WantedItemsOwnerID string                    `json:"wanted_items_owner_id"`
	Items              []*ItemLockCompletedEvent `json:"items"`
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
	ID                 string                         `json:"id"`
	OwnerID            string                         `json:"owner_id"`
	WantedItemsOwnerID string                         `json:"wanted_items_owner_id"`
	OfferedItems       []*TradeOfferAcceptedItemEvent `json:"offered_items"`
	WantedItems        []*TradeOfferAcceptedItemEvent `json:"wanted_items"`
}

// ItemsTradeCompletedEvent ...
type ItemsTradeCompletedEvent struct {
	ID string `json:"id"`
}

// ParseItemsToItemsLockCompletedEvent ...
func ParseItemsToItemsLockCompletedEvent(lockedBy string, s []*TradeItemCreatedEvent) *ItemsLockCompletedEvent {

	items := []*ItemLockCompletedEvent{}

	for _, item := range s {
		items = append(items, &ItemLockCompletedEvent{
			ID:       item.ID,
			Quantity: int64(item.Quantity),
		})
	}

	return &ItemsLockCompletedEvent{
		LockedBy: lockedBy,
		Items:    items,
	}
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
