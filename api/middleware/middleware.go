package middleware

import (
	"auth-service/apperrors"
	"strings"

	"github.com/labstack/echo/v4"
)

func AccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the access token from the request
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			authParts := strings.Fields(authHeader)
			if len(authParts) == 2 && authParts[0] == "Bearer" {
				accessToken := authParts[1]
				echo.Context.Set(c, "accessToken", accessToken)
			}

			// Clear the Authorization header to ensure it doesn't get processed further
			c.Request().Header.Del("Authorization")
		}
		return next(c)
	}
}

func ServiceIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the access token from the request
		serviceID := c.Request().Header.Get("Service-ID")
		if serviceID == "" {
			return apperrors.ErrInvalidServiceID
		}

		echo.Context.Set(c, "serviceID", serviceID)
		return next(c)
	}
}
