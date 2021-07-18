package core

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Authenticate ...
type Authenticate struct {
	Secret string
}

// NewAuthenticate ...
func NewAuthenticate(secret string) *Authenticate {
	return &Authenticate{
		Secret: secret,
	}
}

// Middleware ...
func (a *Authenticate) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearToken := ctx.GetHeader("Authorization")
		strArr := strings.Split(bearToken, "Bearer ")

		if len(strArr) != 2 {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		tokenString := strArr[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.Secret), nil
		})

		if err != nil {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		userID, ok := claims["user_id"]
		if !ok {
			ctx.Status(http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
	}
}
