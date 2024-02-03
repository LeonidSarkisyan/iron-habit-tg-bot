package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

func GenerateJWTToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID

	expiredHours, err := strconv.Atoi(os.Getenv("EXPIRED"))
	if err != nil {
		return "", err
	}

	claims["exp"] = time.Now().Add(time.Hour * time.Duration(expiredHours)).Unix()

	secretKey := []byte(os.Getenv("SECRET_KEY"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractUserIDFromToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("SECRET_KEY"))
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return 0, fmt.Errorf("невозможно получить claims из токена")
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("невозможно получить user_id из токена")
		}

		userID := uint(userIDFloat)
		return userID, nil
	} else {
		return 0, fmt.Errorf("невалидный токен")
	}
}
