package v1

import (
	"auth-service/api/controller"
	"auth-service/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authController struct {
	authService services.IAuthService
}

func InitAuthController() *authController {
	return &authController{
		authService: services.GetAuthService(),
	}
}

func (c *authController) Login(ctx echo.Context) error {
	url := c.authService.Login(ctx)
	return controller.WriteSuccess(ctx, http.StatusOK, url)
}

func (c *authController) Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	accessToken, err := c.authService.LoginCallback(ctx, code)
	if err != nil {
		return controller.WriteError(ctx, http.StatusInternalServerError, err)
	}

	resp := map[string]interface{}{"accessToken": accessToken}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *authController) Logout(ctx echo.Context) error {
	err := c.authService.Logout(ctx)
	if err != nil {
		return controller.WriteError(ctx, http.StatusInternalServerError, err)
	}

	resp := map[string]interface{}{"message": "Logout successful"}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *authController) Authenticate(ctx echo.Context) error {
	err := c.authService.Authenticate(ctx)
	if err != nil {
		return controller.WriteError(ctx, http.StatusUnauthorized, err)
	}

	resp := map[string]interface{}{"message": "Authentication successful"}
	return ctx.JSON(http.StatusOK, resp)
}
