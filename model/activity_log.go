package model

import (
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

type IActivityLogRepository interface {
	Create(ctx echo.Context, activityLog ActivityLog) error
}
