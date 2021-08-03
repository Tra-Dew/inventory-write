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
	testifyMock "github.com/stretchr/testify/mock"
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

var anyItems = testifyMock.AnythingOfType("[]*inventory.Item")
var anyStrings = testifyMock.AnythingOfType("[]string")

func (s *serviceTestSuite) TestCreateItems() {

	s.repository.On("InsertBulk", anyItems).Return(nil)

	correlationID := uuid.NewString()
	userID := uuid.NewString()

	req := &inventory.CreateItemsRequest{
		Items: []*inventory.CreateItemModel{
			createItemModel(),
			createItemModel(),
			createItemModel(),
		},
	}

	err := s.service.CreateItems(s.ctx, userID, correlationID, req)

	s.assert.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "InsertBulk", 1)
}

func (s *serviceTestSuite) TestCreateItemsInvalidItem() {

	s.repository.On("InsertBulk", anyItems).Return(nil)

	correlationID := uuid.NewString()
	userID := uuid.NewString()

	invalidItem := createItemModel()
	invalidItem.Name = ""

	req := &inventory.CreateItemsRequest{
		Items: []*inventory.CreateItemModel{
			invalidItem,
			createItemModel(),
			createItemModel(),
			createItemModel(),
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

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}

	s.repository.On("Get", ids).Return(items, nil)
	s.repository.On("UpdateBulk", anyItems).Return(nil)

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

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}

	s.repository.On("Get", ids).Return(items, nil)
	s.repository.On("UpdateBulk", anyItems).Return(nil)

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

func (s *serviceTestSuite) TestLockItemsInvalidWantedItem() {

	lockedBy := uuid.NewString()
	ownerID := uuid.NewString()
	wantedItemsOwnerID := uuid.NewString()

	// creating offered items
	offeredItems := createItems(2, ownerID)
	offeredIDs := make([]string, len(offeredItems))
	offeredItemModels := make([]*inventory.LockItemModel, len(offeredItems))

	for i, item := range offeredItems {
		offeredIDs[i] = item.ID
		offeredItemModels[i] = &inventory.LockItemModel{
			ID:       item.ID,
			Quantity: 3,
		}
	}

	// creating wanted items
	wantedItems := createItems(1, wantedItemsOwnerID)
	wantedIDs := make([]string, len(wantedItems))
	wantedItemModels := make([]*inventory.LockItemModel, len(wantedItems))

	for i, item := range wantedItems {
		wantedIDs[i] = item.ID
		wantedItemModels[i] = &inventory.LockItemModel{
			ID:       item.ID,
			Quantity: 3,
		}
	}

	invalidItem := &inventory.LockItemModel{
		ID:       uuid.NewString(),
		Quantity: 3,
	}

	wantedIDs = append(wantedIDs, invalidItem.ID)
	wantedItemModels = append(wantedItemModels, invalidItem)

	s.repository.On("Get", offeredIDs).Return(offeredItems, nil)
	s.repository.On("Get", wantedIDs).Return(wantedItems, nil)
	s.repository.On("UpdateBulk", anyItems).Return(nil)

	req := &inventory.LockItemsRequest{
		LockedBy:           lockedBy,
		OwnerID:            ownerID,
		WantedItemsOwnerID: wantedItemsOwnerID,
		OfferedItems:       offeredItemModels,
		WantedItems:        wantedItemModels,
	}

	err := s.service.LockItems(s.ctx, req)

	s.assert.ErrorIs(core.ErrInvalidWantedItems, err)
	s.repository.AssertNumberOfCalls(s.T(), "Get", 1)
	s.repository.AssertNumberOfCalls(s.T(), "UpdateBulk", 0)
}

func (s *serviceTestSuite) TestLockItems() {

	lockedBy := uuid.NewString()
	ownerID := uuid.NewString()
	wantedItemsOwnerID := uuid.NewString()

	// creating offered items
	offeredItems := createItems(2, ownerID)
	offeredIDs := make([]string, len(offeredItems))
	offeredItemModels := make([]*inventory.LockItemModel, len(offeredItems))

	for i, item := range offeredItems {
		offeredIDs[i] = item.ID
		offeredItemModels[i] = &inventory.LockItemModel{
			ID:       item.ID,
			Quantity: 3,
		}
	}

	// creating wanted items
	wantedItems := createItems(1, wantedItemsOwnerID)
	wantedIDs := make([]string, len(wantedItems))
	wantedItemModels := make([]*inventory.LockItemModel, len(wantedItems))

	for i, item := range wantedItems {
		wantedIDs[i] = item.ID
		wantedItemModels[i] = &inventory.LockItemModel{
			ID:       item.ID,
			Quantity: 3,
		}
	}

	s.repository.On("Get", offeredIDs).Return(offeredItems, nil)
	s.repository.On("Get", wantedIDs).Return(wantedItems, nil)
	s.repository.On("UpdateBulk", anyItems).Return(nil)

	req := &inventory.LockItemsRequest{
		LockedBy:           lockedBy,
		OwnerID:            ownerID,
		WantedItemsOwnerID: wantedItemsOwnerID,
		OfferedItems:       offeredItemModels,
		WantedItems:        wantedItemModels,
	}

	err := s.service.LockItems(s.ctx, req)

	s.assert.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "Get", 2)
	s.repository.AssertNumberOfCalls(s.T(), "UpdateBulk", 1)
}

func (s *serviceTestSuite) TestTradeItems() {

	tradeID := uuid.NewString()
	offeredItemsUserID := "offered-id"
	wantedItemsUserID := "wanted-id"

	offeredItems := []*inventory.Item{
		{
			ID:            uuid.NewString(),
			OwnerID:       offeredItemsUserID,
			Name:          inventory.ItemName("Offered Item 1"),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
			Locks: []*inventory.ItemLock{
				{
					LockedBy: tradeID,
					Quantity: 1,
				},
			},
		},
		{
			ID:            uuid.NewString(),
			OwnerID:       offeredItemsUserID,
			Name:          inventory.ItemName("Offered Item 2"),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
			Locks: []*inventory.ItemLock{
				{
					LockedBy: tradeID,
					Quantity: 4,
				},
			},
		},
	}

	wantedItems := []*inventory.Item{
		{
			ID:            uuid.NewString(),
			OwnerID:       wantedItemsUserID,
			Name:          inventory.ItemName("Wanted Item 1"),
			TotalQuantity: inventory.ItemQuantity(2),
			Status:        inventory.ItemAvailable,
		},
		{
			ID:            uuid.NewString(),
			OwnerID:       wantedItemsUserID,
			Name:          inventory.ItemName("Wanted Item 2"),
			TotalQuantity: inventory.ItemQuantity(5),
			Status:        inventory.ItemAvailable,
		},
	}

	offeredItemsModels := make([]*inventory.TradeItemModel, len(offeredItems))
	offeredItemIDs := make([]string, len(offeredItems))

	for i, item := range offeredItems {
		offeredItemIDs[i] = item.ID
		offeredItemsModels[i] = &inventory.TradeItemModel{
			ID:       item.ID,
			Quantity: int64(item.Locks[0].Quantity),
		}
	}

	wantedItemsModels := make([]*inventory.TradeItemModel, len(wantedItems))
	wantedItemIDs := make([]string, len(wantedItems))

	for i, item := range wantedItems {
		wantedItemIDs[i] = item.ID
		wantedItemsModels[i] = &inventory.TradeItemModel{
			ID:       item.ID,
			Quantity: int64(item.TotalQuantity),
		}
	}

	s.repository.On("Get", offeredItemIDs).Return(offeredItems, nil)
	s.repository.On("Get", wantedItemIDs).Return(wantedItems, nil)
	s.repository.On("UpdateBulk", testifyMock.MatchedBy(func(items []*inventory.Item) bool {
		var (
			countOffered, countWanted int
		)

		for _, item := range items {
			if item.OwnerID == wantedItemsUserID {
				countWanted++
			}
			if item.OwnerID == offeredItemsUserID {
				countOffered++
			}
		}

		return countOffered == 2 && countWanted == 0
	})).Return(nil)
	s.repository.On("InsertBulk", testifyMock.MatchedBy(func(items []*inventory.Item) bool {
		var (
			countOffered, countWanted int
		)

		for _, item := range items {
			if item.OwnerID == wantedItemsUserID {
				countWanted++
			}
			if item.OwnerID == offeredItemsUserID {
				countOffered++
			}
		}

		return countOffered == 2 && countWanted == 2
	})).Return(nil)

	s.repository.On("DeleteBulk", testifyMock.MatchedBy(func(ids []string) bool {
		return len(ids) == len(wantedItems)
	})).Return(nil)

	req := &inventory.TradeItemsRequest{
		TradeID:            tradeID,
		OwnerID:            offeredItemsUserID,
		WantedItemsOwnerID: wantedItemsUserID,
		OfferedItems:       offeredItemsModels,
		WantedItems:        wantedItemsModels,
	}

	err := s.service.TradeItems(s.ctx, req)

	s.assert.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "Get", 2)
	s.repository.AssertNumberOfCalls(s.T(), "UpdateBulk", 1)
	s.repository.AssertNumberOfCalls(s.T(), "InsertBulk", 1)
	s.repository.AssertNumberOfCalls(s.T(), "DeleteBulk", 1)
}

func createItemModel() *inventory.CreateItemModel {
	faker.SetRandomStringLength(15)

	description := faker.Sentence()
	return &inventory.CreateItemModel{
		Name:        faker.Name(),
		Description: &description,
		Quantity:    5,
	}
}

func createItems(amount int, ownerID string) []*inventory.Item {

	faker.SetRandomStringLength(15)

	items := make([]*inventory.Item, amount)

	for i := 0; i < amount; i++ {
		description := faker.Sentence()
		ints, _ := faker.RandomInt(5, 10)

		items[i], _ = inventory.NewItem(
			uuid.NewString(),
			ownerID,
			faker.Name(),
			&description,
			int64(ints[0]),
			inventory.ItemAvailable,
		)
	}

	return items
}
