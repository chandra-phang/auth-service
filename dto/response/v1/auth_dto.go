package v1

type LoginDTO struct {
	LoginUrl string `json:"loginUrl"`
}

type AuthenticateDTO struct {
	Message string `json:"message"`
}
