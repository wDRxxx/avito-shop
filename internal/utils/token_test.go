package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/avito-shop/internal/service/models"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	secretKey := "secret"
	tests := []struct {
		name     string
		user     *models.UserClaims
		duration time.Duration
		hasErr   bool
	}{
		{
			name: "success case",
			user: &models.UserClaims{
				RegisteredClaims: jwt.RegisteredClaims{Subject: "123"},
				Username:         "testuser",
			},
			duration: time.Hour,
			hasErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.user, secretKey, tt.duration)
			if tt.hasErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				claims, err := VerifyToken(token, secretKey)
				require.NoError(t, err)
				require.Equal(t, tt.user.Subject, claims.Subject)
			}
		})
	}
}

func TestUserIDFromToken(t *testing.T) {
	secretKey := "secret"
	tests := []struct {
		name     string
		user     *models.UserClaims
		duration time.Duration
		id       int
		hasErr   bool
	}{
		{
			name: "success case",
			user: &models.UserClaims{
				RegisteredClaims: jwt.RegisteredClaims{Subject: "123"},
				Username:         "testuser",
			},
			duration: time.Hour,
			id:       123,
			hasErr:   false,
		},
		{
			name:     "invalid token",
			user:     nil,
			duration: 0,
			id:       0,
			hasErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token string
			var err error
			if tt.user != nil {
				token, err = GenerateToken(tt.user, secretKey, tt.duration)
				require.NoError(t, err)
			} else {
				token = "invalid.token.here"
			}

			userID, err := UserIDFromToken(token, secretKey)
			require.Equal(t, tt.hasErr, err != nil)
			require.Equal(t, tt.id, userID)
		})
	}
}
