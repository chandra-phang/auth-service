package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// contextKey is a custom type for putting auth-related values into a context.
type contextKey string

// userContextKey is a key for saving a User object into a context.
const AccessTokenKey contextKey = "accessToken"

func AccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the access token from the request
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			authParts := strings.Fields(authHeader)
			if len(authParts) == 2 && authParts[0] == "Bearer" {
				accessToken := authParts[1]
				echo.Context.Set(c, string(AccessTokenKey), accessToken)
			}

			// Clear the Authorization header to ensure it doesn't get processed further
			c.Request().Header.Del("Authorization")
		}
		return next(c)
	}
}
