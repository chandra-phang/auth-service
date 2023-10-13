package api

import (
	v1 "auth-service/api/controller/v1"
	"log"

	"github.com/labstack/echo/v4"
)

func InitRoutes() {
	e := echo.New()

	productController := v1.InitAuthController()
	v1Api := e.Group("v1")

	v1Api.GET("/login", productController.HandleLogin)
	v1Api.GET("/callback", productController.HandleCallback)
	v1Api.GET("/logout", productController.HandleLogout)
	v1Api.GET("/authenticate", productController.HandleAuthenticate)

	log.Println("Server is running at 8081 port.")
	e.Logger.Fatal(e.Start(":8081"))
}
