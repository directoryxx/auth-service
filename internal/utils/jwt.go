package utils

import (
	"auth/internal/domain"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	Uuid string `json:"uuid"`
	jwt.StandardClaims
}

func GenerateToken(user *domain.User, uuid string) (token interface{}, err error) {
	jwtHourExpire := os.Getenv("JWT_EXPIRE_HOUR")
	convJwtHour, _ := strconv.Atoi(jwtHourExpire)
	// Set custom claims
	claims := &JwtCustomClaims{
		uuid,
		jwt.StandardClaims{
			Issuer:    "Auth Service",
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(convJwtHour)).Unix(),
		},
	}

	// Create token with claims
	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := tokenJwt.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return t, err
}
