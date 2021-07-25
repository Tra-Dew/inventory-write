package inventory_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/d-leme/tradew-inventory-write/pkg/core"
	"github.com/d-leme/tradew-inventory-write/pkg/inventory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type domainTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(domainTestSuite))
}

func (s *domainTestSuite) SetupSuite() {
	s.assert = assert.New(s.T())
}

func (s *domainTestSuite) TestNewItem() {
	description := faker.Sentence()

	item, err := inventory.NewItem(
		uuid.NewString(),
		uuid.NewString(),
		faker.Name(),
		&description,
		5,
		inventory.ItemAvailable,
	)

	s.assert.NoError(err)
	s.assert.NotNil(item)
}

func (s *domainTestSuite) TestNewItemInvalidName() {
	description := faker.Sentence()

	item, err := inventory.NewItem(
		uuid.NewString(),
		uuid.NewString(),
		"x",
		&description,
		5,
		inventory.ItemAvailable,
	)

	s.assert.Error(err)
	s.assert.Nil(item)
}

func (s *domainTestSuite) TestNewItemInvalidQuantity() {
	description := faker.Sentence()

	item, err := inventory.NewItem(
		uuid.NewString(),
		uuid.NewString(),
		faker.Name(),
		&description,
		0,
		inventory.ItemAvailable,
	)

	s.assert.Error(err)
	s.assert.Nil(item)
}

func (s *domainTestSuite) TestUpdate() {
	name := faker.Name()
	description := faker.Sentence()
	quantity := int64(5)

	item, err := inventory.NewItem(
		uuid.NewString(),
		uuid.NewString(),
		name,
		&description,
		quantity,
		inventory.ItemAvailable,
	)

	s.assert.NoError(err)
	s.assert.NotNil(item)

	name = "new name"
	description = "new description"
	quantity = 7
	err = item.Update("new name", &description, quantity)

	s.assert.NoError(err)
	s.assert.Equal(name, string(item.Name))
	s.assert.Equal(description, string(*item.Description))
	s.assert.Equal(quantity, int64(item.TotalQuantity))
}

func (s *domainTestSuite) TestLock() {
	description := faker.Sentence()
	quantity := int64(5)

	item, err := inventory.NewItem(
		uuid.NewString(),
		uuid.NewString(),
		faker.Name(),
		&description,
		quantity,
		inventory.ItemAvailable,
	)

	s.assert.NoError(err)
	s.assert.NotNil(item)

	lockQuantity := int64(3)
	err = item.Lock(uuid.NewString(), lockQuantity)

	s.assert.NoError(err)
	s.assert.Equal(quantity, int64(item.TotalQuantity))
	s.assert.Equal(lockQuantity, int64(item.GetLockedQuantity()))
}

func (s *domainTestSuite) TestLockNotEnoughtItemsToLock() {
	description := faker.Sentence()
	quantity := int64(5)

	item, err := inventory.NewItem(
		uuid.NewString(),
		uuid.NewString(),
		faker.Name(),
		&description,
		quantity,
		inventory.ItemAvailable,
	)

	s.assert.NoError(err)
	s.assert.NotNil(item)

	err = item.Lock(uuid.NewString(), 7)

	s.assert.ErrorIs(err, core.ErrNotEnoughtItemsToLock)
	s.assert.Equal(quantity, int64(item.TotalQuantity))
	s.assert.Equal(int64(0), int64(item.GetLockedQuantity()))
}
