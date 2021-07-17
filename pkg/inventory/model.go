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
	Items []UpdateItemModel `json:"items"`
}
