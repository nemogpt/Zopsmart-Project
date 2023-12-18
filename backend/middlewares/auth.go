package middlewares

import (
	"net/http"
	"time"

	"gofr.dev/gofr/config"
	"gofr.dev/gofr"
)

func NewAuthMiddleware(cfg *config.Config) gofr.Handler {
	return func(c gofr.Context) error {
		// Parse the token from the Authorization header.
		token, err := utils.ParseTokenFromHeader(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(gofr.Map{
				"error": "unauthorized",
			})
		}

		// Get the token payload.
		payload, ok := token.Claims.(*utils.StandardClaims)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(gofr.Map{
				"error": "malformed token",
			})
		}

		// Verify the token expiration.
		if payload.ExpiresAt < time.Now().Unix() {
			return c.Status(http.StatusUnauthorized).JSON(gofr.Map{
				"error": "token expired",
			})
		}

		// Pass the context to the next handler.
		c.Locals("user", payload)
		return c.Next()
	}
}
