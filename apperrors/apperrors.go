package apperrors

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrAccessTokenNotFound  = errors.New("accessToken not found")
	ErrAccessTokenIsEmpty   = errors.New("accessToken is empty")
	ErrAccessTokenIsExpired = errors.New("accessToken is expired")
	ErrInvalidServiceID     = errors.New("invalid Service-ID")
)
