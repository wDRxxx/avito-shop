package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	em "github.com/wDRxxx/avito-shop/internal/e2e/models"
	"github.com/wDRxxx/avito-shop/internal/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func TAuth(t *testing.T) {
	tests := []struct {
		name         string
		username     string
		password     string
		expectStatus int
	}{
		{
			name:         "success case",
			username:     gofakeit.Username(),
			password:     gofakeit.Password(true, true, true, true, false, 12),
			expectStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/api/auth", apiURL)
			client := http.Client{}

			reqBody := &models.AuthRequest{
				Username: tt.username,
				Password: tt.password,
			}
			r, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(r))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			res, err := client.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tt.expectStatus, res.StatusCode)

			var authResponse em.AuthResponse
			err = utils.ReadJSON(res.Body, &authResponse)
			require.NoError(t, err)
			assert.Empty(t, authResponse.Errors)

			claims, err := utils.VerifyToken(authResponse.Token, os.Getenv("TOKEN_SECRET"))
			require.NoError(t, err)
			require.Equal(t, tt.username, claims.Username)

			users = append(users, em.User{
				Username: tt.username,
				Token:    authResponse.Token,
			})
		})
	}
}
