package app

import (
	"auth-service/api"
	"auth-service/config"
)

type Application struct {
}

// Returns a new instance of the application
func NewApplication() Application {
	return Application{}
}

func (a Application) InitApplication() {
	config.InitConfig()
	api.InitRoutes()
}
