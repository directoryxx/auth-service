package api

import (
	"auth/infrastructure"
	"auth/internal/controller"
	"auth/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

var UserController controller.UserController

type errorresponse struct {
	Message string
}

func Routes(
	router *echo.Echo,
	UserController controller.UserController,
) {

	router.POST("/register", UserController.Register)
	router.POST("/login", UserController.Login)

	router.Use(authMiddleware)
	router.GET("/profile", UserController.Profile)

}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		whitelistUrl := []string{"/login", "/register"}
		redisConn := infrastructure.OpenRedis()

		defer redisConn.Close()

		ctx := c.Request().Context()

		// Register public route
		if slices.Contains(whitelistUrl, c.Request().RequestURI) {
			return next(c)
		}

		// Check header Authorization is empty
		if c.Request().Header.Get("Authorization") == "" {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		// Get value Header Authorization
		token := c.Request().Header.Get("Authorization")

		// Check value Header contain bearer
		if !strings.Contains(token, "Bearer ") {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		// Delete bearer and left only token
		tokenFix := strings.Replace(token, "Bearer ", "", 1)

		// Parse token & verify secret
		claims := jwt.MapClaims{}
		tokenParse, errFix := jwt.ParseWithClaims(tokenFix, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		// Return if cant parse and validation
		if errFix != nil {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		// Convert parse to claim
		claim := tokenParse.Claims.(jwt.MapClaims)
		// convert from interface to string
		uuid := fmt.Sprintf("%v", claim["uuid"])

		resUuid, _ := redisConn.Get(ctx, uuid).Result()

		// check uuid if uuid exist pass it
		if resUuid != "" {
			// Parse return to struct
			user := domain.User{}
			json.Unmarshal([]byte(resUuid), &user)
			// set to context
			c.Set("user", user)
			return next(c)
		} else {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}
	}
}
