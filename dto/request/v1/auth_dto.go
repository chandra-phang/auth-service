package v1

import (
	"auth-service/apperrors"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthenticateDTO struct {
	SourceUri string `json:"sourceUri" validate:"required"`
}

func (dto AuthenticateDTO) Validate(ctx echo.Context) error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		vErr := apperrors.TryTranslateValidationErrors(err)
		return errors.New(vErr)
	}

	return nil
}
