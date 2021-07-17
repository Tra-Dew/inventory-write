package core

import "github.com/gin-gonic/gin"

// Controller ...
type Controller interface {
	RegisterRoutes(r *gin.RouterGroup)
}
