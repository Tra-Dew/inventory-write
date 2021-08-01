package mock

import (
	"context"

	"github.com/d-leme/tradew-inventory-write/pkg/inventory"
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
func (r *RepositoryMock) UpdateBulk(ctx context.Context, items []*inventory.Item) error {
	args := r.Mock.Called(items)

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

// Get ...
func (r *RepositoryMock) Get(ctx context.Context, userID *string, ids []string) ([]*inventory.Item, error) {
	args := r.Mock.Called(ids)

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.([]*inventory.Item), nil
	}

	arg1 := args.Get(1)

	return nil, arg1.(error)
}

// GetByStatus ...
func (r *RepositoryMock) GetByStatus(ctx context.Context, status inventory.ItemStatus) ([]*inventory.Item, error) {
	args := r.Mock.Called()

	arg0 := args.Get(0)
	if arg0 != nil {
		return arg0.([]*inventory.Item), nil
	}

	arg1 := args.Get(1)

	return nil, arg1.(error)
}
