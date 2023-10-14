package services

import "auth-service/handlers"

func InitServices(h handlers.Handler) {
	InitAuthService(h)
}
