package repositories

import (
	"auth-service/apperrors"
	"auth-service/model"
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepositoryInstance(db *sql.DB) model.IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) CreateUser(ctx echo.Context, user model.User) error {
	sqlStatement := `
		INSERT INTO users
			(id, name, email, external_id, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	params := []interface{}{
		user.ID,
		user.Name,
		user.Email,
		user.ExternalID,
		user.CreatedAt,
		user.UpdatedAt,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetUserByID(ctx echo.Context, externalID string) (*model.User, error) {
	sqlStatement := `
		SELECT
			id,
			name,
			email,
			external_id,
			created_at,
			updated_at
		FROM users
		WHERE id = ?
	`

	results, err := r.db.Query(sqlStatement, externalID)
	if err != nil {
		return nil, err
	}

	var user model.User
	for results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("failed to scan", err)
			return nil, err
		}
	}

	if user.ID == "" {
		return nil, apperrors.ErrUserNotFound
	}

	return &user, nil
}

func (r UserRepository) GetUserByEmailAndExternalID(ctx echo.Context, email string, externalID string) (*model.User, error) {
	sqlStatement := `
		SELECT
			id,
			name,
			email,
			external_id,
			created_at,
			updated_at
		FROM users
		WHERE email = ? AND external_id = ?
	`

	params := []interface{}{email, externalID}
	results, err := r.db.Query(sqlStatement, params...)
	if err != nil {
		return nil, err
	}

	var user model.User
	for results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.ExternalID, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println("failed to scan", err)
			return nil, err
		}
	}

	if user.ID == "" {
		return nil, apperrors.ErrUserNotFound
	}

	return &user, nil
}
