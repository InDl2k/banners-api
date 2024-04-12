package jwt

import (
	"banners/internal/config"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

func CreateToken(username string) (string, error) {
	tokenTTL := config.MustLoad().JWTToken.TTL
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     getRole(username),
		"exp":      time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
		"iat":      time.Now().Unix(),
	})

	token, err := claims.SignedString([]byte(config.MustLoad().JWTToken.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Имитация БД
func getRole(username string) string {
	if username == "admin" {
		return "admin"
	}
	return "user"
}

func ValidateJWT(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("Token is empty")
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.MustLoad().JWTToken.Secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("Invalid token")
	}

	return role, nil
}
