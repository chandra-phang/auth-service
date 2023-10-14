package repositories

import (
	"auth-service/apperrors"
	"auth-service/model"
	"database/sql"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

type AccessTokenRepository struct {
	db *sql.DB
}

func NewAccessTokenRepositoryInstance(db *sql.DB) model.IAccessTokenRepository {
	return &AccessTokenRepository{
		db: db,
	}
}

func (r AccessTokenRepository) CreateAccessToken(ctx echo.Context, accessToken model.AccessToken) error {
	sqlStatement := `
		INSERT INTO access_tokens
			(id, user_id, token_string, expired_at, created_at, updated_at)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	params := []interface{}{
		accessToken.ID,
		accessToken.UserID,
		accessToken.TokenString,
		accessToken.ExpiredAt,
		accessToken.CreatedAt,
		accessToken.UpdatedAt,
	}

	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}

func (r AccessTokenRepository) GetAccessTokenByTokenString(ctx echo.Context, tokenString string) (*model.AccessToken, error) {
	sqlStatement := `
		SELECT
			id,
			user_id,
			token_string,
			expired_at,
			created_at,
			updated_at
		FROM access_tokens
		WHERE token_string = ?
	`

	results, err := r.db.Query(sqlStatement, tokenString)
	if err != nil {
		return nil, err
	}

	var accessToken model.AccessToken
	for results.Next() {
		err = results.Scan(&accessToken.ID, &accessToken.UserID, &accessToken.TokenString, &accessToken.ExpiredAt, &accessToken.CreatedAt, &accessToken.UpdatedAt)
		if err != nil {
			log.Println("failed to scan", err)
			return nil, err
		}
	}

	if accessToken.ID == "" {
		return nil, apperrors.ErrAccessTokenNotFound
	}

	return &accessToken, nil
}

func (r AccessTokenRepository) GetActiveAccessTokensByUserID(ctx echo.Context, userID string) ([]model.AccessToken, error) {
	sqlStatement := `
		SELECT
			id,
			user_id,
			token_string,
			expired_at,
			created_at,
			updated_at
		FROM access_tokens
		WHERE user_id = ? AND (expired_at IS NULL OR expired_at < ?)
	`

	params := []interface{}{userID, time.Now()}
	results, err := r.db.Query(sqlStatement, params...)
	if err != nil {
		return nil, err
	}

	var accessTokens = make([]model.AccessToken, 0)
	for results.Next() {
		var accessToken model.AccessToken
		err = results.Scan(&accessToken.ID, &accessToken.UserID, &accessToken.TokenString, &accessToken.ExpiredAt, &accessToken.CreatedAt, &accessToken.UpdatedAt)
		if err != nil {
			log.Println("failed to scan", err)
			return nil, err
		}

		accessTokens = append(accessTokens, accessToken)
	}

	return accessTokens, nil
}

func (r AccessTokenRepository) RevokeAccessTokenByTokenString(ctx echo.Context, tokenString string) error {
	sqlStatement := `
		UPDATE access_tokens
		SET
			expired_at = ?
		WHERE token_string = ?
	`

	params := []interface{}{time.Now(), tokenString}
	_, err := r.db.Exec(sqlStatement, params...)
	if err != nil {
		return err
	}

	return nil
}
