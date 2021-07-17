package mock

import (
	"context"

	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/stretchr/testify/mock"
)

// RepositoryMock ...
type RepositoryMock struct {
	mock.Mock
}

// NewRepository ...
func NewRepository() inventory.Repository {
	return &RepositoryMock{}
}

// InsertBulk ...
func (r *RepositoryMock) InsertBulk(ctx context.Context, items []*inventory.Item) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}

// UpdateBulk ...
func (r *RepositoryMock) UpdateBulk(ctx context.Context, userID string, items []*inventory.UpdateItem) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}

// DeleteBulk ...
func (r *RepositoryMock) DeleteBulk(ctx context.Context, userID string, ids []string) error {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.(error)
	}

	return nil
}
