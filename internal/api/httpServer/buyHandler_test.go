package httpServer

import (
	"context"
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

func TestBuyHandler(t *testing.T) {
	t.Parallel()

	type serviceMockFunc func(mc *minimock.Controller) service.Service

	var (
		ctx       = context.Background()
		itemTitle = "title"
		authCfg   = config.NewMockAuthConfig()

		serviceErr = errors.New("service error")
		userID     = int(gofakeit.Uint8())

		user = &sm.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: fmt.Sprint(userID),
			},
			Username: gofakeit.Username(),
		}
	)

	token, err := utils.GenerateToken(user, authCfg.TokenSecret(), authCfg.TokenTTL())
	require.NoError(t, err)

	type args struct {
		ctx   context.Context
		item  string
		token string
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
				ctx:   ctx,
				item:  itemTitle,
				token: "Bearer " + token,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.BuyItemMock.Expect(ctx, userID, itemTitle).Return(nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name: "item not found case",
			args: args{
				ctx:   ctx,
				item:  itemTitle,
				token: "Bearer " + token,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.BuyItemMock.Expect(ctx, userID, itemTitle).Return(service.ErrItemNotFound)
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrItemNotFound,
		},
		{
			name: "insufficient balance case",
			args: args{
				ctx:   ctx,
				item:  itemTitle,
				token: "Bearer " + token,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.BuyItemMock.Expect(ctx, userID, itemTitle).Return(service.ErrInsufficientBalance)
				return mockService
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  api.ErrInsufficientBalance,
		},
		{
			name: "service error case",
			args: args{
				ctx:   ctx,
				item:  itemTitle,
				token: "Bearer " + token,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.BuyItemMock.Expect(ctx, userID, itemTitle).Return(serviceErr)
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  api.ErrInternal,
		},
		{
			name: "service error case",
			args: args{
				ctx:   ctx,
				item:  itemTitle,
				token: "Bearer " + token,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.BuyItemMock.Expect(ctx, userID, itemTitle).Return(serviceErr)
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  api.ErrInternal,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/buy/%s", tt.args.item), nil)
			r.Header.Set("Authorization", tt.args.token)

			serviceMock := tt.serviceMock(minimock.NewController(t))
			serv := newMockHTTPServer(authCfg, config.NewMockHttpConfig(), serviceMock)

			serv.BuyHandler(w, r)

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
