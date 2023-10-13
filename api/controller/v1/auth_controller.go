package v1

import (
	"auth-service/config"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type authController struct {
	config config.Config
}

func InitAuthController() *authController {
	return &authController{
		config: *config.GetConfig(),
	}
}

func (c *authController) HandleLogin(ctx echo.Context) error {
	url := c.config.OauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return ctx.JSON(http.StatusOK, map[string]string{"redirectUrl": url})
}

func (c *authController) HandleCallback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	token, err := c.config.OauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	userInfo, err := fetchUserInfo(token.AccessToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	// Store the token and associate it with the user session

	return ctx.JSON(http.StatusOK, map[string]interface{}{"access_token": token.AccessToken, "userInfo": userInfo})
}

func (c *authController) HandleLogout(ctx echo.Context) error {
	// Revoke the user's access token and refresh token
	// Optionally: destroy the user session

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Logout successful"})
}

func (c *authController) HandleAuthenticate(ctx echo.Context) error {
	accessToken := ctx.QueryParam("access_token")
	// Validate the access token, check user permissions, etc.
	log.Println(accessToken)
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Authentication successful"})
}

func fetchUserInfo(accessToken string) (map[string]interface{}, error) {
	// Make a request to the UserInfo endpoint to fetch user details
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
