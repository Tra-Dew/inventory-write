package inventory

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
