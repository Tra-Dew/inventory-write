package inventory

import (
	"net/http"

	"github.com/d-leme/tradew-inventory-write/pkg/core"
	"github.com/gin-gonic/gin"
)

// Controller ...
type Controller struct {
	authenticate *core.Authenticate
	settings     *core.Settings
	service      Service
}

// NewController ...
func NewController(settings *core.Settings, authenticate *core.Authenticate, service Service) Controller {
	return Controller{
		settings:     settings,
		authenticate: authenticate,
		service:      service,
	}
}

// RegisterRoutes ...
func (c *Controller) RegisterRoutes(r *gin.RouterGroup) {
	inventory := r.Group("/inventory-write")
	{
		inventory.Use(
			c.authenticate.Middleware(),
		)

		inventory.POST("", c.post)
		inventory.PUT("", c.put)
		inventory.DELETE("", c.delete)
	}
}

func (c *Controller) post(ctx *gin.Context) {
	req := new(CreateItemsRequest)
	correlationID := ctx.GetString("X-Correlation-ID")
	userID := ctx.GetString("user_id")

	if err := ctx.ShouldBindJSON(req); err != nil {
		core.HandleRestError(ctx, core.ErrMalformedJSON)
		return
	}

	if err := c.service.CreateItems(ctx, userID, correlationID, req); err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) put(ctx *gin.Context) {
	req := new(UpdateItemsRequest)
	correlationID := ctx.GetString("X-Correlation-ID")
	userID := ctx.GetString("user_id")

	if err := ctx.ShouldBindJSON(req); err != nil {
		core.HandleRestError(ctx, core.ErrMalformedJSON)
		return
	}

	if err := c.service.UpdateItems(ctx, userID, correlationID, req); err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) delete(ctx *gin.Context) {
	req := new(DeleteItemsRequest)
	correlationID := ctx.GetString("X-Correlation-ID")
	userID := ctx.GetString("user_id")

	if err := ctx.ShouldBindJSON(req); err != nil {
		core.HandleRestError(ctx, core.ErrMalformedJSON)
		return
	}

	if err := c.service.DeleteItems(ctx, userID, correlationID, req); err != nil {
		core.HandleRestError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
