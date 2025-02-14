package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/models"
	"github.com/wDRxxx/avito-shop/internal/models/converter"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/mocks"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func TestUserInfoHandler(t *testing.T) {
	t.Parallel()

	type serviceMockFunc func(mc *minimock.Controller) service.Service

	var (
		ctx     = context.Background()
		userID  = int(gofakeit.Uint8())
		authCfg = config.NewMockAuthConfig()

		inventoryItems = []*sm.InventoryItem{
			{
				ID:       int(gofakeit.Uint8()),
				Title:    gofakeit.BeerName(),
				Quantity: int(gofakeit.Uint8()),
			},
		}

		incomingTransactions = []*sm.IncomingTransaction{
			{
				Amount:         int(gofakeit.Uint8()),
				SenderUsername: gofakeit.Username(),
			},
		}
		outgoingTransactions = []*sm.OutgoingTransaction{
			{
				Amount:            int(gofakeit.Uint8()),
				RecipientUsername: gofakeit.Username(),
			},
		}

		userInfo = &sm.UserInfo{
			Balance:              int(gofakeit.Uint8()),
			InventoryItems:       inventoryItems,
			IncomingTransactions: incomingTransactions,
			OutgoingTransactions: outgoingTransactions,
		}

		expectedBody  = converter.UserInfoFromServiceToApi(userInfo)
		errorResponse = models.ErrorResponse{Errors: api.ErrInternal.Error()}

		serviceErr = errors.New("service error")
	)

	// Генерация токена для успешного случая
	user := &sm.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: fmt.Sprint(userID),
		},
		Username: gofakeit.Username(),
	}

	token, err := utils.GenerateToken(user, authCfg.TokenSecret(), authCfg.TokenTTL())
	require.NoError(t, err)

	type args struct {
		ctx   context.Context
		token string
	}

	tests := []struct {
		name           string
		args           args
		serviceMock    serviceMockFunc
		expectedStatus int
		expectedError  error
		expectedBody   any
	}{
		{
			name: "success case",
			args: args{
				ctx:   ctx,
				token: "Bearer " + token,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.UserInfoMock.Expect(ctx, userID).Return(userInfo, nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedError:  nil,
			expectedBody:   expectedBody,
		},
		{
			name: "invalid token case",
			args: args{
				ctx:   ctx,
				token: "Bearer invalidToken",
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  api.ErrInternal,
			expectedBody:   errorResponse,
		},
		{
			name: "service error case",
			args: args{
				ctx:   ctx,
				token: "Bearer " + token,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.UserInfoMock.Expect(ctx, userID).Return(nil, serviceErr)
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  api.ErrInternal,
			expectedBody:   errorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/info", nil)
			r.Header.Set("Authorization", tt.args.token)

			serviceMock := tt.serviceMock(minimock.NewController(t))
			serv := newMockHTTPServer(authCfg, config.NewMockHttpConfig(), serviceMock)

			serv.UserInfoHandler(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != nil && w.Code != tt.expectedStatus {
				var errorResp *models.ErrorResponse
				err = utils.ReadJSON(w.Body, &errorResp)
				require.NoError(t, err)
				t.Error(errorResp.Errors)
			}

			expectedBody, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			require.Equal(t, expectedBody, w.Body.Bytes())
		})
	}
}
