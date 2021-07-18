package inventory_test

import (
	"testing"

	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/bxcodec/faker/v3"
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
