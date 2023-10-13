package api

import (
	v1 "auth-service/api/controller/v1"
	"log"

	"github.com/labstack/echo/v4"
)

func InitRoutes() {
	e := echo.New()

	authController := v1.InitAuthController()

	v1Api := e.Group("v1")
	v1Api.GET("/login", authController.Login)
	v1Api.GET("/callback", authController.Callback)
	v1Api.GET("/logout", authController.Logout)
	v1Api.GET("/authenticate", authController.Authenticate)

	log.Println("Server is running at 8081 port.")
	e.Logger.Fatal(e.Start(":8081"))
}
