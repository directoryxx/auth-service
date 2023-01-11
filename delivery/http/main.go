package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"auth/infrastructure"
	"auth/internal/controller"
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/internal/usecase"
	"auth/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"auth/cmd/api"
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func RunApi() {
	app := echo.New()

	app.Validator = &CustomValidator{validator: validator.New()}

	log.Println("[INFO] Starting Auth Service on port", os.Getenv("APPLICATION_PORT"))

	log.Println("[INFO] Loading Database")
	dbSQL, err := infrastructure.Open()

	if err != nil {
		log.Fatalf("Could not initialize Database connection using sqlx %s", err)
	}

	defer dbSQL.Close()

	log.Println("[INFO] Loading Redis")
	redisConnect := infrastructure.OpenRedis()

	defer redisConnect.Close()

	log.Println("[INFO] Loading Kafka Producer")
	kafkaProducer, err := infrastructure.ConnectKafka()

	if err != nil {
		log.Fatalf("Could not initialize connection to kafka producer %s", err)
	}

	defer kafkaProducer.Close()

	log.Println("[INFO] Loading Repository")
	userRepo := repository.NewUserRepository(dbSQL, redisConnect, kafkaProducer)

	log.Println("[INFO] Loading Usecase")
	userUsecase := usecase.NewUserUseCase(userRepo)

	log.Println("[INFO] Loading Controller")
	userController := controller.NewUserController(userUsecase)

	log.Println("[INFO] Loading Middleware")
	SetMiddleware(app, userRepo)

	log.Println("[INFO] Loading Routes")
	api.Routes(app, userController)

	log.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("APPLICATION_PORT"))))
}

func SetMiddleware(r *echo.Echo, userRepo repository.UserRepository) {
	// Middleware
	r.Use(middleware.Logger())
	r.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogUserAgent: true,
		LogMethod:    true,
		LogHost:      true,
		LogLatency:   true,
		LogRequestID: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			// log :=
			hostname, _ := os.Hostname()
			log := &domain.Log{
				RemoteIP:      values.RemoteIP,
				Service:       os.Getenv("APPLICATION_NAME") + " Service",
				ContainerName: hostname,
				Time:          values.StartTime.String(),
				Host:          values.Host,
				Method:        values.Method,
				Uri:           values.URI,
				UserAgent:     values.UserAgent,
				Status:        strconv.Itoa(values.Status),
				Latency:       values.Latency.String(),
				LatencyHuman:  values.Latency.String(),
				// Error:         values.Error.Error(),
			}

			b, _ := json.Marshal(log)

			userRepo.Publish(c.Request().Context(), string(b), "logger")

			return nil
		},
	}))
	r.Use(middleware.Recover())

	// Cors
	r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "http://localhost"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
}

func SetPrivateMiddleware(r *echo.Echo) {
	config := middleware.JWTConfig{
		Claims:     &utils.JwtCustomClaims{},
		SigningKey: []byte("secret"),
		Skipper: func(c echo.Context) bool {
			// Skip middleware if path is equal 'login'
			if c.Request().URL.Path == "/login" {
				return true
			}
			if c.Request().URL.Path == "/register" {
				return true
			}
			return false
		},
	}
	r.Use(middleware.JWTWithConfig(config))
}
