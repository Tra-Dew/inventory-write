package cmd

import (
	"fmt"

	"github.com/Tra-Dew/inventory-write/pkg/core"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Server is a cmd to setup api server
func Server(command *cobra.Command, args []string) {

	settings := new(core.Settings)

	if err := core.FromYAML(command.Flag("settings").Value.String(), settings); err != nil {
		logrus.
			WithError(err).
			Fatal("unable to parse settings, shutting down...")
		return
	}

	container := NewContainer(settings)

	defer container.Close()

	engine := configureAPI(container, settings)

	logrus.WithField("port", settings.Port).Info("starting server")

	err := engine.Run(fmt.Sprintf(":%d", settings.Port))
	logrus.Fatal(err)
}

func configureAPI(container *Container, settings *core.Settings) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// middleware
	engine.Use(gin.CustomRecovery(core.InternalErrorRecovery()))
	engine.Use(core.CorrelationIDMiddleware())
	engine.Use(core.LogMiddleware("2006-01-02T15:04:05Z07:00"))

	// helth check
	engine.GET("/health", core.HTTPHealth())

	// routes
	rg := engine.Group("/api/v1")

	for _, c := range container.Controllers() {
		c.RegisterRoutes(rg)
	}

	return engine
}
