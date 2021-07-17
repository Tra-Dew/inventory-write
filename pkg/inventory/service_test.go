package inventory_test

import (
	"context"
	"testing"

	"github.com/Tra-Dew/inventory-write/pkg/core"
	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/Tra-Dew/inventory-write/pkg/inventory/mock"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type serviceTestSuite struct {
	suite.Suite
	assert     *assert.Assertions
	repository *mock.RepositoryMock
	service    inventory.Service
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (s *serviceTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
}

func (s *serviceTestSuite) SetupTest() {
	s.repository = mock.NewRepository().(*mock.RepositoryMock)
	s.service = inventory.NewService(s.repository)
}

func (s *serviceTestSuite) TestCreate() {

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

	err := s.service.CreateItems(context.TODO(), userID, correlationID, req)

	s.assert.NoError(err)
	s.repository.AssertNumberOfCalls(s.T(), "InsertBulk", 1)
}

func (s *serviceTestSuite) TestCreateInvalidItem() {

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

	err := s.service.CreateItems(context.TODO(), userID, correlationID, req)

	s.assert.ErrorIs(err, core.ErrValidationFailed)
	s.repository.AssertNumberOfCalls(s.T(), "InsertBulk", 0)
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
