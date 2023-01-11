package controller

import (
	"auth/internal/domain"
	"auth/internal/usecase"
	"auth/internal/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type registerresponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Token   any    `json:"token"`
}

type profileresponse struct {
	Error   bool         `json:"error"`
	Message string       `json:"message"`
	Profile *domain.User `json:"profile"`
}

type loginresponse struct {
	Error   bool `json:"error"`
	Message any  `json:"message"`
	Data    any  `json:"data"`
	Token   any  `json:"token"`
}

type errorresponse struct {
	Error   bool `json:"error"`
	Message any  `json:"message"`
}

// interface
type UserController interface {
	Login(ec echo.Context) error
	Register(ec echo.Context) error
	Profile(ec echo.Context) error
}

// implement interface
type UserControllerImpl struct {
	UserUsecase usecase.UserUseCase
}

func NewUserController(userUsecase usecase.UserUseCase) UserController {
	return &UserControllerImpl{
		UserUsecase: userUsecase,
	}
}

func (uc *UserControllerImpl) Register(c echo.Context) error {
	// Convert Echo Context
	con := c.Request().Context()
	ctx, cancel := context.WithTimeout(con, 10000*time.Second)
	defer cancel()

	// Validation
	u := new(domain.RegisterValidation)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(u); err != nil {
		fmt.Println(err)
		return err
	}

	// Registering User
	user, uuid, err := uc.UserUsecase.Register(ctx, u)

	if err != nil {
		response := errorresponse{
			Error:   true,
			Message: err.Error(),
		}
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	// Generate Token
	token, _ := utils.GenerateToken(user, uuid)

	response := registerresponse{
		Error:   false,
		Message: "Berhasil mendaftar",
		Data:    user,
		Token:   token,
	}

	return c.JSON(http.StatusOK, response)
}

func (uc *UserControllerImpl) Login(c echo.Context) error {
	// Convert echo context
	con := c.Request().Context()
	ctx, cancel := context.WithTimeout(con, 10000*time.Second)
	defer cancel()

	// Validation
	u := new(domain.LoginValidation)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(u); err != nil {
		return err
	}

	// Check credentials
	login, uuidGen, err := uc.UserUsecase.Login(ctx, u)

	if err != nil {
		response := errorresponse{
			Error:   true,
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Generate Token
	token, _ := utils.GenerateToken(login, uuidGen)

	response := loginresponse{
		Error:   false,
		Message: "Berhasil login",
		Data:    login,
		Token:   token,
	}

	return c.JSON(http.StatusOK, response)
}

func (uc *UserControllerImpl) Profile(c echo.Context) error {
	// Convert echo context
	con := c.Request().Context()
	_, cancel := context.WithTimeout(con, 10000*time.Second)
	defer cancel()

	// Get JWT Content
	user := c.Get("user").(domain.User)

	response := &profileresponse{
		Error:   false,
		Message: "Berhasil mengambil data",
		Profile: &user,
	}

	return c.JSON(http.StatusOK, response)
}
