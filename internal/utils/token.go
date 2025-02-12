package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/wDRxxx/avito-shop/internal/service/models"
)

func GenerateToken(user *models.UserClaims, secretKey string, duration time.Duration) (string, error) {
	claims := models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "avito-shop",
			Subject:   user.Subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string, secretKey string) (*models.UserClaims, error) {
	t, err := jwt.ParseWithClaims(token, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*models.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func UserIDFromToken(token string, secretKey string) (int, error) {
	claims, err := VerifyToken(token, secretKey)
	if err != nil {
		return 0, err
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
