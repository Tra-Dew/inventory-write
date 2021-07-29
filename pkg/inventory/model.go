package inventory

// CreateItemModel ...
type CreateItemModel struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Quantity    int64   `json:"quantity"`
}

// CreateItemsRequest ...
type CreateItemsRequest struct {
	Items []*CreateItemModel `json:"items"`
}

// UpdateItemModel ...
type UpdateItemModel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Quantity    int64   `json:"quantity"`
}

// UpdateItemsRequest ...
type UpdateItemsRequest struct {
	Items []*UpdateItemModel `json:"items"`
}

// LockItemModel ...
type LockItemModel struct {
	ID       string
	Quantity int64
}

// LockItemsRequest ...
type LockItemsRequest struct {
	LockedBy string
	Items    []*LockItemModel
}

// TradeItemModel ...
type TradeItemModel struct {
	ID       string
	Quantity int64
}

// TradeItemsRequest ...
type TradeItemsRequest struct {
	TradeID            string
	OwnerID            string
	WantedItemsOwnerID string
	OfferedItems       []*TradeItemModel
	WantedItems        []*TradeItemModel
}

// DeleteItemsRequest ...
type DeleteItemsRequest struct {
	IDs []string `json:"ids"`
}
