package inventory

// CreateItemsRequest ...
type CreateItemsRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Quantity    int64   `json:"quantity"`
}

// CreateItemsResponse ...
type CreateItemsResponse struct {
	IDs []string `json:"ids"`
}

// DeleteItemsRequest ...
type DeleteItemsRequest struct {
	IDs []string `json:"ids"`
}
