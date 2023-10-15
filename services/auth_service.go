package services

import (
	"auth-service/api/middleware"
	"auth-service/apperrors"
	"auth-service/config"
	reqV1 "auth-service/dto/request/v1"
	respV1 "auth-service/dto/response/v1"
	"auth-service/handlers"
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
	Authenticate(ctx echo.Context, dto reqV1.AuthenticateDTO) (*model.User, error)
}

type authSvc struct {
	config          config.Config
	dbCon           *sql.DB
	userRepo        model.IUserRepository
	accessTokenRepo model.IAccessTokenRepository
	activityLogRepo model.IActivityLogRepository
}

var authSvcSingleton IAuthService

func InitAuthService(h handlers.Handler) {
	authSvcSingleton = authSvc{
		config:          *config.GetConfig(),
		dbCon:           h.DB,
		userRepo:        repositories.NewUserRepositoryInstance(h.DB),
		accessTokenRepo: repositories.NewAccessTokenRepositoryInstance(h.DB),
		activityLogRepo: repositories.NewActivityLogRepositoryInstance(h.DB),
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

	// get user by email and externalID
	user, err := svc.userRepo.GetUserByEmailAndExternalID(ctx, userInfo.Email, userInfo.ExternalID)
	if err != nil && err != apperrors.ErrUserNotFound {
		return "", err
	}

	// create user if it's not exist
	if user == nil {
		user = new(model.User).Initialize(userInfo.Name, userInfo.Email, userInfo.ExternalID)
		err = svc.userRepo.CreateUser(ctx, *user)
		if err != nil {
			return "", err
		}
	}

	// retrieve user's active accessTokens and revoke them
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

	// create a new accessToken with the new tokenString
	accessToken := new(model.AccessToken).Initialize(user.ID, token.AccessToken)
	err = svc.accessTokenRepo.CreateAccessToken(ctx, *accessToken)
	if err != nil {
		return "", err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func (svc authSvc) Logout(ctx echo.Context) error {
	tokenString := ctx.Get(middleware.AccessTokenKey)
	if tokenString == nil {
		return apperrors.ErrAccessTokenIsEmpty
	}

	accessToken, err := svc.accessTokenRepo.GetAccessTokenByTokenString(ctx, tokenString.(string))
	if err != nil {
		return apperrors.ErrAccessTokenNotFound
	}

	err = svc.accessTokenRepo.RevokeAccessTokenByTokenString(ctx, accessToken.TokenString)
	if err != nil {
		return err
	}

	return nil
}

func (svc authSvc) Authenticate(ctx echo.Context, dto reqV1.AuthenticateDTO) (*model.User, error) {
	tokenString := ctx.Get(middleware.AccessTokenKey)
	if tokenString == nil {
		return nil, apperrors.ErrAccessTokenIsEmpty
	}

	// add DB transaction
	tx, _ := svc.dbCon.Begin()
	defer tx.Rollback()

	accessToken, err := svc.accessTokenRepo.GetAccessTokenByTokenString(ctx, tokenString.(string))
	if err != nil {
		return nil, apperrors.ErrAccessTokenNotFound
	}

	if accessToken.ExpiredAt != nil && accessToken.ExpiredAt.Before(time.Now()) {
		return nil, apperrors.ErrAccessTokenIsExpired
	}

	activityLog := new(model.ActivityLog).Initialize(accessToken.UserID, dto.SourceUri)
	err = svc.activityLogRepo.Create(ctx, *activityLog)
	if err != nil {
		return nil, err
	}

	user, err := svc.userRepo.GetUserByID(ctx, accessToken.UserID)
	if err != nil {
		return nil, err
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (svc authSvc) fetchUserInfo(accessToken string) (*respV1.UserInfoDTO, error) {
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

	var userInfo respV1.UserInfoDTO
	if err := json.Unmarshal(respBytes, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
