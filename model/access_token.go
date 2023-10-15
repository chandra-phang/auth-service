package model

import (
	"auth-service/lib"
	"time"

	"github.com/labstack/echo/v4"
)

type AccessToken struct {
	ID          string
	UserID      string
	TokenString string
	ExpiredAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (a AccessToken) Initialize(userID string, token string) *AccessToken {
	return &AccessToken{
		ID:          lib.GenerateUUID(),
		UserID:      userID,
		TokenString: token,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type IAccessTokenRepository interface {
	CreateAccessToken(ctx echo.Context, user AccessToken) error
	GetAccessTokenByTokenString(ctx echo.Context, tokenString string) (*AccessToken, error)
	GetActiveAccessTokensByUserID(ctx echo.Context, userID string) ([]AccessToken, error)
	RevokeAccessTokenByTokenString(ctx echo.Context, tokenString string) error
}
