package app

import (
	"auth-service/api"
	"auth-service/config"
	"auth-service/db"
	"auth-service/handlers"
	"auth-service/services"
)

type Application struct {
}

// Returns a new instance of the application
func NewApplication() Application {
	return Application{}
}

func (a Application) InitApplication() {
	config.InitConfig()

	database := db.InitConnection()
	h := handlers.New(database)
	services.InitServices(h)

	api.InitRoutes()

	db.CloseConnection(database)
}
