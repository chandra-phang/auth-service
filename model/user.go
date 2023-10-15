package model

import (
	"auth-service/lib"
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID         string
	Name       string
	Email      string
	ExternalID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (u User) Initialize(name string, email string, externalID string) *User {
	return &User{
		ID:         lib.GenerateUUID(),
		Name:       name,
		Email:      email,
		ExternalID: externalID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

type IUserRepository interface {
	CreateUser(ctx echo.Context, user User) error
	GetUserByID(ctx echo.Context, ID string) (*User, error)
	GetUserByEmailAndExternalID(ctx echo.Context, email string, externalID string) (*User, error)
}
