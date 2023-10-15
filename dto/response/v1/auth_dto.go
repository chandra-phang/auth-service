package v1

import "auth-service/model"

type LoginDTO struct {
	LoginUrl string `json:"loginUrl"`
}

type LoginCallbackDTO struct {
	AccessToken string `json:"accessToken"`
}

type LogoutDTO struct {
	Message string `json:"message"`
}

type AuthenticateDTO struct {
	User UserDTO `json:"user"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (dto *AuthenticateDTO) ConvertFromUser(user model.User) *AuthenticateDTO {
	return &AuthenticateDTO{
		User: UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}
}
