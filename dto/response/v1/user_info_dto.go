package v1

type UserInfoDTO struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	ExternalID string `json:"sub"`
}
