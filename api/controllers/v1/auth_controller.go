package v1

import (
	"auth-service/api/controllers"
	v1req "auth-service/dto/request/v1"
	v1resp "auth-service/dto/response/v1"
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
	resp := v1resp.LoginDTO{LoginUrl: url}
	return controllers.WriteSuccess(ctx, http.StatusOK, resp)
}

func (c *authController) Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	accessToken, err := c.authService.LoginCallback(ctx, code)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusInternalServerError, err)
	}

	resp := v1resp.LoginCallbackDTO{AccessToken: accessToken}
	return ctx.JSON(http.StatusOK, resp)
}

func (c *authController) Logout(ctx echo.Context) error {
	err := c.authService.Logout(ctx)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusInternalServerError, err)
	}

	resp := v1resp.LogoutDTO{Message: "Logout success!"}
	return controllers.WriteSuccess(ctx, http.StatusOK, resp)
}

func (c *authController) Authenticate(ctx echo.Context) error {
	reqBody, _ := ioutil.ReadAll(ctx.Request().Body)
	dto := v1req.AuthenticateDTO{}

	if err := json.Unmarshal(reqBody, &dto); err != nil {
		return controllers.WriteError(ctx, http.StatusBadRequest, err)
	}

	err := dto.Validate(ctx)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusBadRequest, err)
	}

	user, err := c.authService.Authenticate(ctx, dto)
	if err != nil {
		return controllers.WriteError(ctx, http.StatusUnauthorized, err)
	}

	resp := new(v1resp.AuthenticateDTO).ConvertFromUser(*user)
	return controllers.WriteSuccess(ctx, http.StatusOK, resp)
}
