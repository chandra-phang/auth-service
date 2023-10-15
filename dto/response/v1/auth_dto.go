package v1

type LoginDTO struct {
	LoginUrl string `json:"loginUrl"`
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
