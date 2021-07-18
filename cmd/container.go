package cmd

import (
	"github.com/Tra-Dew/inventory-write/pkg/core"
	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/Tra-Dew/inventory-write/pkg/inventory/memory"
)

// Container contains all depencies from our api
type Container struct {
	Settings *core.Settings

	Authenticate *core.Authenticate

	InventoryRepository inventory.Repository
	InventoryService    inventory.Service
	InventoryController inventory.Controller
}

// NewContainer creates new instace of Container
func NewContainer(settings *core.Settings) *Container {

	container := new(Container)

	container.Settings = settings

	container.Authenticate = core.NewAuthenticate(settings.JWT.Secret)

	container.InventoryRepository = memory.NewRepository()
	container.InventoryService = inventory.NewService(container.InventoryRepository)
	container.InventoryController = inventory.NewController(settings, container.Authenticate, container.InventoryService)

	return container
}

// Controllers maps all routes and exposes them
func (c *Container) Controllers() []core.Controller {
	return []core.Controller{
		&c.InventoryController,
	}
}

// Close terminates every opened resource
func (c *Container) Close() {}
