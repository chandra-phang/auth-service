package v1

import (
	"auth-service/api/controllers"
	v1request "auth-service/dto/request/v1"
	v1response "auth-service/dto/response/v1"
	"auth-service/services"
	"encoding/json"
	"io/ioutil"
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
	dto := v1response.LoginDTO{LoginUrl: url}
	return controllers.WriteSuccess(ctx, http.StatusOK, dto)
}

func (c *authController) Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	accessToken, err := c.authService.LoginCallback(ctx, code)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusInternalServerError, err)
	}

	resp := map[string]interface{}{"accessToken": accessToken}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *authController) Logout(ctx echo.Context) error {
	err := c.authService.Logout(ctx)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusInternalServerError, err)
	}

	dto := v1response.AuthenticateDTO{Message: "Logout success!"}
	return controllers.WriteSuccess(ctx, http.StatusOK, dto)
}

func (c *authController) Authenticate(ctx echo.Context) error {
	reqBody, _ := ioutil.ReadAll(ctx.Request().Body)
	dto := v1request.AuthenticateDTO{}

	if err := json.Unmarshal(reqBody, &dto); err != nil {
		return controllers.WriteError(ctx, http.StatusBadRequest, err)
	}

	err := dto.Validate(ctx)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusBadRequest, err)
	}

	err = c.authService.Authenticate(ctx, dto)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusUnauthorized, err)
	}

	responseDTO := v1response.AuthenticateDTO{Message: "Authentication success!"}
	return controllers.WriteSuccess(ctx, http.StatusOK, responseDTO)
}
