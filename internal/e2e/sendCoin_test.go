package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/avito-shop/internal/api"
	em "github.com/wDRxxx/avito-shop/internal/e2e/models"
	"github.com/wDRxxx/avito-shop/internal/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func TSendCoin(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        *models.SendCoinRequest
		expectedStatus int
		expectedError  string
		authHeader     string
	}{
		{
			name: "successCase",
			reqBody: &models.SendCoinRequest{
				ToUser: users[1].Username,
				Amount: 10,
			},
			expectedStatus: http.StatusOK,
			authHeader:     "Bearer " + users[0].Token,
		},
		{
			name: "no auth header error case",
			reqBody: &models.SendCoinRequest{
				ToUser: users[1].Username,
				Amount: 10,
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "recipient doesn't exist error case",
			reqBody: &models.SendCoinRequest{
				ToUser: "no",
				Amount: 10,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrUserNotFound.Error(),
			authHeader:     "Bearer " + users[0].Token,
		},
		{
			name: "insufficient balance error case",
			reqBody: &models.SendCoinRequest{
				ToUser: users[1].Username,
				Amount: 10000000,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrInsufficientBalance.Error(),
			authHeader:     "Bearer " + users[0].Token,
		},
		{
			name: "send to yourself error case",
			reqBody: &models.SendCoinRequest{
				ToUser: users[0].Username,
				Amount: 10,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrSendToYourself.Error(),
			authHeader:     "Bearer " + users[0].Token,
		},
	}

	url := fmt.Sprintf("%s/api/sendCoin", apiURL)
	client := http.Client{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Fatal("error marshaling request body", err)
			}

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(r))
			if err != nil {
				t.Fatal("error creating request", err)
			}
			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatal("error executing request", err)
			}
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectedError != "" {
				var sendResponse em.ErrorResponse
				err = utils.ReadJSON(res.Body, &sendResponse)
				if err != nil {
					t.Fatal("error reading response body", err)
				}

				require.Equal(t, tt.expectedError, sendResponse.Errors)
			}
		})
	}
}
