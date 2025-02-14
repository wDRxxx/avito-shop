package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	em "github.com/wDRxxx/avito-shop/internal/e2e/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func TBuy(t *testing.T) {
	tests := []struct {
		name         string
		authHeader   string
		expectStatus int
	}{
		{
			name:         "success case",
			authHeader:   "Bearer " + users[0].Token,
			expectStatus: http.StatusOK,
		},
		{
			name:         "no auth header error case",
			authHeader:   "",
			expectStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/api/buy/%s", "http://localhost:8080", "cup")
			client := http.Client{}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err, "error creating request")

			req.Header.Set("Content-Type", "application/json")
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			res, err := client.Do(req)
			require.NoError(t, err, "error executing request")
			defer res.Body.Close()

			assert.Equal(t, tt.expectStatus, res.StatusCode)

			if res.StatusCode != http.StatusOK {
				var buyResponse em.ErrorResponse
				err = utils.ReadJSON(res.Body, &buyResponse)
				require.NoError(t, err, "error reading response body")
				assert.NotEmpty(t, buyResponse.Errors, "expected error message from API")
			}
		})
	}
}
