package services

import (
	"auth-service/apperrors"
	"auth-service/config"
	"auth-service/dto/response"
	"auth-service/handlers"
	"auth-service/lib"
	"auth-service/model"
	"auth-service/repositories"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type IAuthService interface {
	Login(ctx echo.Context) string
	LoginCallback(ctx echo.Context, code string) (string, error)
	Logout(ctx echo.Context) error
	Authenticate(ctx echo.Context) error
}

type authSvc struct {
	config          config.Config
	dbCon           *sql.DB
	userRepo        model.IUserRepository
	accessTokenRepo model.IAccessTokenRepository
}

var authSvcSingleton IAuthService

func InitAuthService(h handlers.Handler) {
	authSvcSingleton = authSvc{
		config:          *config.GetConfig(),
		dbCon:           h.DB,
		userRepo:        repositories.NewUserRepositoryInstance(h.DB),
		accessTokenRepo: repositories.NewAccessTokenRepositoryInstance(h.DB),
	}
}

func GetAuthService() IAuthService {
	return authSvcSingleton
}

func (svc authSvc) Login(ctx echo.Context) string {
	return svc.config.OauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (svc authSvc) LoginCallback(ctx echo.Context, code string) (string, error) {
	token, err := svc.config.OauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", err
	}

	userInfo, err := svc.fetchUserInfo(token.AccessToken)
	if err != nil {
		return "", err
	}

	// add DB transaction
	tx, _ := svc.dbCon.Begin()
	defer tx.Rollback()

	user, err := svc.userRepo.GetUserByEmailAndExternalID(ctx, userInfo.Email, userInfo.ExternalID)
	if err != nil && err != apperrors.ErrUserNotFound {
		return "", err
	}

	// create user if it's not exist
	if user == nil {
		user = &model.User{
			ID:         lib.GenerateUUID(),
			Name:       userInfo.Name,
			Email:      userInfo.Email,
			ExternalID: userInfo.ExternalID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err = svc.userRepo.CreateUser(ctx, *user)
		if err != nil {
			return "", err
		}
	}

	accessTokens, err := svc.accessTokenRepo.GetActiveAccessTokensByUserID(ctx, user.ID)
	if err != nil && err != apperrors.ErrAccessTokenNotFound {
		return "", err
	}
	for _, accessToken := range accessTokens {
		err := svc.accessTokenRepo.RevokeAccessTokenByTokenString(ctx, accessToken.TokenString)
		if err != nil {
			return "", nil
		}
	}

	// create a new accessToken with new tokenString
	accessToken := &model.AccessToken{
		ID:          lib.GenerateUUID(),
		UserID:      user.ID,
		TokenString: token.AccessToken,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = svc.accessTokenRepo.CreateAccessToken(ctx, *accessToken)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func (svc authSvc) Logout(ctx echo.Context) error {
	tokenString := ctx.Get("accessToken").(string)
	if tokenString == "" {
		return apperrors.ErrAccessTokenIsEmpty
	}

	accessToken, err := svc.accessTokenRepo.GetAccessTokenByTokenString(ctx, tokenString)
	if err != nil {
		return apperrors.ErrAccessTokenNotFound
	}

	err = svc.accessTokenRepo.RevokeAccessTokenByTokenString(ctx, accessToken.TokenString)
	if err != nil {
		return err
	}

	return nil
}

func (svc authSvc) Authenticate(ctx echo.Context) error {
	tokenString := ctx.Get("accessToken").(string)
	if tokenString == "" {
		return apperrors.ErrAccessTokenIsEmpty
	}

	accessToken, err := svc.accessTokenRepo.GetAccessTokenByTokenString(ctx, tokenString)
	if err != nil {
		return apperrors.ErrAccessTokenNotFound
	}

	if accessToken.ExpiredAt != nil && accessToken.ExpiredAt.Before(time.Now()) {
		return apperrors.ErrAccessTokenIsExpired
	}

	return nil
}

func (svc authSvc) fetchUserInfo(accessToken string) (*response.UserInfoDTO, error) {
	// Make a request to the UserInfo endpoint to fetch user details
	queryParams := fmt.Sprintf("?access_token=%s", accessToken)
	resp, err := http.Get(svc.config.UserInfoURL + queryParams)
	if err != nil {
		return nil, err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo response.UserInfoDTO
	if err := json.Unmarshal(respBytes, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
