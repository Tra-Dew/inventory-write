package inventory_test

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/d-leme/tradew-inventory-write/pkg/core"
	"github.com/d-leme/tradew-inventory-write/pkg/inventory"
	"github.com/d-leme/tradew-inventory-write/pkg/inventory/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	ctx        context.Context
	repository *mock.RepositoryMock
	service    inventory.Service
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (s *serviceTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
	s.ctx = context.Background()
}

func (s *serviceTestSuite) SetupTest() {
	s.repository = mock.NewRepository().(*mock.RepositoryMock)
	s.service = inventory.NewService(s.repository)
}

func (s *serviceTestSuite) TestCreateItems() {

	s.repository.On("InsertBulk").Return(nil)

	correlationID := uuid.NewString()
	userID := uuid.NewString()

	req := &inventory.CreateItemsRequest{
		Items: []*inventory.CreateItemModel{
			createItem(),
			createItem(),
			createItem(),
		},
	}

	err := s.service.CreateItems(s.ctx, userID, correlationID, req)

	s.assert.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "InsertBulk", 1)
}

func (s *serviceTestSuite) TestCreateItemsInvalidItem() {

	s.repository.On("InsertBulk").Return(nil)

	correlationID := uuid.NewString()
	userID := uuid.NewString()

	invalidItem := createItem()
	invalidItem.Name = ""

	req := &inventory.CreateItemsRequest{
		Items: []*inventory.CreateItemModel{
			invalidItem,
			createItem(),
			createItem(),
			createItem(),
		},
	}

	err := s.service.CreateItems(s.ctx, userID, correlationID, req)

	s.assert.ErrorIs(err, core.ErrValidationFailed)
	s.repository.AssertNumberOfCalls(s.T(), "InsertBulk", 0)
}

func (s *serviceTestSuite) TestUpdateItems() {

	correlationID := uuid.NewString()
	userID := uuid.NewString()

	items := []*inventory.Item{
		{
			ID:            uuid.NewString(),
			OwnerID:       userID,
			Name:          inventory.ItemName(faker.Name()),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
		},
		{
			ID:            uuid.NewString(),
			OwnerID:       userID,
			Name:          inventory.ItemName(faker.Name()),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
		},
	}

	s.repository.On("Get").Return(items, nil)
	s.repository.On("UpdateBulk").Return(nil)

	itemModels := make([]*inventory.UpdateItemModel, len(items))

	for i, item := range items {
		itemModels[i] = &inventory.UpdateItemModel{
			ID:       item.ID,
			Name:     string(item.Name),
			Quantity: 10,
		}
	}

	req := &inventory.UpdateItemsRequest{Items: itemModels}

	err := s.service.UpdateItems(s.ctx, userID, correlationID, req)

	s.assert.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "Get", 1)
	s.repository.AssertNumberOfCalls(s.T(), "UpdateBulk", 1)
}

func (s *serviceTestSuite) TestUpdateItemsInvalidItem() {

	correlationID := uuid.NewString()
	userID := uuid.NewString()

	items := []*inventory.Item{
		{
			ID:            uuid.NewString(),
			OwnerID:       userID,
			Name:          inventory.ItemName(faker.Name()),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
		},
		{
			ID:            uuid.NewString(),
			OwnerID:       userID,
			Name:          inventory.ItemName(faker.Name()),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
		},
	}

	s.repository.On("Get").Return(items, nil)
	s.repository.On("UpdateBulk").Return(nil)

	itemModels := make([]*inventory.UpdateItemModel, len(items))

	for i, item := range items {
		itemModels[i] = &inventory.UpdateItemModel{
			ID:       item.ID,
			Name:     string(item.Name),
			Quantity: -1,
		}
	}

	req := &inventory.UpdateItemsRequest{Items: itemModels}

	err := s.service.UpdateItems(s.ctx, userID, correlationID, req)

	s.assert.ErrorIs(err, core.ErrValidationFailed)
	s.repository.AssertNumberOfCalls(s.T(), "Get", 1)
	s.repository.AssertNumberOfCalls(s.T(), "UpdateBulk", 0)
}

func (s *serviceTestSuite) TestLockItems() {

	userID := uuid.NewString()

	items := []*inventory.Item{
		{
			ID:            uuid.NewString(),
			OwnerID:       userID,
			Name:          inventory.ItemName(faker.Name()),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
		},
		{
			ID:            uuid.NewString(),
			OwnerID:       userID,
			Name:          inventory.ItemName(faker.Name()),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
		},
	}

	s.repository.On("Get").Return(items, nil)
	s.repository.On("UpdateBulk").Return(nil)

	itemModels := make([]*inventory.LockItemModel, len(items))

	for i, item := range items {
		itemModels[i] = &inventory.LockItemModel{
			ID:       item.ID,
			Quantity: 3,
		}
	}

	req := &inventory.LockItemsRequest{Items: itemModels}

	err := s.service.LockItems(s.ctx, userID, req)

	s.assert.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "Get", 1)
	s.repository.AssertNumberOfCalls(s.T(), "UpdateBulk", 1)
}

func createItem() *inventory.CreateItemModel {
	faker.SetRandomStringLength(15)

	description := faker.Sentence()
	return &inventory.CreateItemModel{
		Name:        faker.Name(),
		Description: &description,
		Quantity:    5,
	}
}
