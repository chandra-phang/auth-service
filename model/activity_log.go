package model

import (
	"auth-service/lib"
	"time"

	"github.com/labstack/echo/v4"
)

type ActivityLog struct {
	ID        string
	UserID    string
	SourceUri string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a ActivityLog) Initialize(userID string, sourceUri string) *ActivityLog {
	return &ActivityLog{
		ID:        lib.GenerateUUID(),
		UserID:    userID,
		SourceUri: sourceUri,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type IActivityLogRepository interface {
	Create(ctx echo.Context, activityLog ActivityLog) error
}
