package httpServer

import (
	"bytes"
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
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/mocks"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func TestSendCoinHandler(t *testing.T) {
	t.Parallel()

	type serviceMockFunc func(mc *minimock.Controller) service.Service

	var (
		ctx     = context.Background()
		authCfg = config.NewMockAuthConfig()

		serviceErr = errors.New("service error")
		userID     = int(gofakeit.Uint8())
		toUser     = gofakeit.Username()
		amount     = int(gofakeit.Uint8())
		username   = gofakeit.Username()

		user = &sm.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: fmt.Sprint(userID),
			},
			Username: username,
		}

		reqBody = &models.SendCoinRequest{
			ToUser: toUser,
			Amount: amount,
		}
	)

	token, err := utils.GenerateToken(user, authCfg.TokenSecret(), authCfg.TokenTTL())
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		token   string
		request *models.SendCoinRequest
	}

	tests := []struct {
		name           string
		args           args
		serviceMock    serviceMockFunc
		expectedStatus int
		expectedError  error
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				token:   "Bearer " + token,
				request: reqBody,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.SendCoinMock.Expect(ctx, toUser, userID, amount).Return(nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedError:  nil,
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
		},
		{
			name: "send to yourself case",
			args: args{
				ctx:   ctx,
				token: "Bearer " + token,
				request: &models.SendCoinRequest{
					ToUser: username,
					Amount: amount,
				},
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrSendToYourself,
		},
		{
			name: "service err case",
			args: args{
				ctx:     ctx,
				token:   "Bearer " + token,
				request: reqBody,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.SendCoinMock.Expect(ctx, toUser, userID, amount).Return(serviceErr)
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  api.ErrInternal,
		},
		{
			name: "user not found err case",
			args: args{
				ctx:     ctx,
				token:   "Bearer " + token,
				request: reqBody,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.SendCoinMock.Expect(ctx, toUser, userID, amount).Return(service.ErrUserNotFound)
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrUserNotFound,
		},
		{
			name: "insufficient balance err case",
			args: args{
				ctx:     ctx,
				token:   "Bearer " + token,
				request: reqBody,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.SendCoinMock.Expect(ctx, toUser, userID, amount).Return(service.ErrInsufficientBalance)
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrInsufficientBalance,
		},
		{
			name: "wrong set amount err case",
			args: args{
				ctx:   ctx,
				token: "Bearer " + token,
				request: &models.SendCoinRequest{
					ToUser: toUser,
					Amount: 0,
				},
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrWrongSendCoinAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody, err := json.Marshal(tt.args.request)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(reqBody))
			r.Header.Set("Authorization", tt.args.token)

			serviceMock := tt.serviceMock(minimock.NewController(t))
			serv := newMockHTTPServer(authCfg, config.NewMockHttpConfig(), serviceMock)

			serv.SendCoinHandler(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != nil && w.Code != tt.expectedStatus {
				var errorResp *models.ErrorResponse
				err = utils.ReadJSON(w.Body, &errorResp)
				require.NoError(t, err)
				t.Error(errorResp.Errors)
			}
		})
	}
}
