package repositories

import (
	"auth-service/model"
	"database/sql"

	"github.com/labstack/echo/v4"
)

type ActivityLogRepository struct {
	db *sql.DB
}

func NewActivityLogRepositoryInstance(db *sql.DB) model.IActivityLogRepository {
	return &ActivityLogRepository{
		db: db,
	}
}

func (r ActivityLogRepository) Create(ctx echo.Context, activityLog model.ActivityLog) error {
	sqlStatement := `
		INSERT INTO activity_logs
			(id, user_id, source_uri, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?)
	`

	params := []interface{}{
		activityLog.ID,
		activityLog.UserID,
		activityLog.SourceUri,
		activityLog.CreatedAt,
		activityLog.UpdatedAt,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}
