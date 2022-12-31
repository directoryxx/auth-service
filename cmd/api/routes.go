package api

import (
	"auth/internal/controller"

	"github.com/labstack/echo/v4"
)

var UserController controller.UserController

func Routes(
	router *echo.Echo,
	UserController controller.UserController,
) {

	router.POST("/register", UserController.Register)
	router.POST("/login", UserController.Login)

	router.GET("/profile", UserController.Profile)

}

func PrivateRoutes(
	router *echo.Echo,
	UserController controller.UserController,
) {

	router.GET("/profile", UserController.Profile)

}
