package core

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	// CorrelationIDHeader key of correlation id header
	CorrelationIDHeader = "X-Correlation-ID"
)

// InternalErrorRecovery actions to perform after a panic from the server
func InternalErrorRecovery() gin.RecoveryFunc {
	return func(c *gin.Context, err interface{}) {
		logrus.
			WithFields(logrus.Fields{
				"err":   err,
				"stack": debug.Stack(),
			}).
			Errorf("server panic")

		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// CorrelationIDMiddleware sets the correlation id header
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := ctx.Request.Header.Get(CorrelationIDHeader)

		id, err := uuid.Parse(h)
		if err != nil {
			id = uuid.New()
		}

		ctx.Set(CorrelationIDHeader, id.String())
	}
}

// LogMiddleware sets up and logs requests
func LogMiddleware(timeFormat string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		path := ctx.Request.URL.Path
		ctx.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)
		uuid, _ := ctx.Get(CorrelationIDHeader)

		entry := logrus.WithFields(logrus.Fields{
			"status":         ctx.Writer.Status(),
			"method":         ctx.Request.Method,
			"path":           path,
			"ip":             ctx.ClientIP(),
			"latency":        latency,
			"user-agent":     ctx.Request.UserAgent(),
			"time":           end.Format(timeFormat),
			"correlation_id": uuid,
		})

		if len(ctx.Errors) > 0 {
			entry.Error(ctx.Errors.String())
		} else {
			entry.Info()
		}
	}
}
