package core

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Health ...
type Health struct {
	checks []func(*Health) error
}

// HealthOption ...
type HealthOption func(*Health)

// NewHealth ...
func NewHealth(options ...HealthOption) *Health {
	h := new(Health)
	h.checks = []func(*Health) error{}

	for _, o := range options {
		o(h)
	}

	return h
}

// HTTPHealth ...
func HTTPHealth(options ...HealthOption) gin.HandlerFunc {
	return NewHealth(options...).HTTP()
}

// HTTP ...
func (h *Health) HTTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := h.Health(); err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusServiceUnavailable, "service unavailable")
			return
		}

		ctx.JSON(http.StatusOK, "healthy")
	}
}

// Health ...
func (h *Health) Health() error {
	wg := new(sync.WaitGroup)

	errCh := make(chan error, len(h.checks))
	doneCh := make(chan bool, len(h.checks))

	for _, check := range h.checks {
		wg.Add(1)
		go func(c func(*Health) error) {
			defer wg.Done()
			if err := c(h); err != nil {
				errCh <- err
			}
		}(check)
	}

	go func() {
		wg.Wait()
		doneCh <- true
	}()

	<-doneCh

	close(errCh)
	close(doneCh)

	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}
