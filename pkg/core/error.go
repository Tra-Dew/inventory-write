package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error used as a wrapper for all application errors
type Error struct {
	Key string
}

func newError(key string) *Error {
	return &Error{Key: key}
}

func (e *Error) Error() string {
	return e.Key
}

var (
	// ErrValidationFailed returned when an entity has a invalid field
	ErrValidationFailed = newError("validation-failed")

	// ErrMalformedJSON returned when a json request could not be parsed
	ErrMalformedJSON = newError("malformed-json")

	// ErrNotFound returned when an entity is not found
	ErrNotFound = newError("not-found")

	// ErrInvalidCredentials returned when the password or email is invalid
	ErrInvalidCredentials = newError("invalid-credentials")
)

// RestError used as a Rest api call error
type RestError struct {
	Key string `json:"key"`
}

// ErrorStatusMap mapping between application erros and status codes
var ErrorStatusMap = map[string]int{
	ErrValidationFailed.Key:   http.StatusUnprocessableEntity,
	ErrMalformedJSON.Key:      http.StatusUnprocessableEntity,
	ErrInvalidCredentials.Key: http.StatusBadRequest,
	ErrNotFound.Key:           http.StatusNotFound,
}

// HandleRestError handles applications errors using ErrorStatusMap
func HandleRestError(ctx *gin.Context, err error) {

	if ierr, ok := err.(*Error); ok {
		if s, exists := ErrorStatusMap[ierr.Key]; exists {
			ctx.JSON(s, &RestError{Key: ierr.Key})
			return
		}
	}

	ctx.JSON(http.StatusInternalServerError, &RestError{Key: "internal-server-error"})
	return
}
