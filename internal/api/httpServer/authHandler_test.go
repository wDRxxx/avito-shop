package httpServer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/models"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/mocks"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func TestAuthHandler(t *testing.T) {
	t.Parallel()

	type serviceMockFunc func(mc *minimock.Controller) service.Service

	var (
		ctx       = context.Background()
		mockToken = "mockedToken"

		req = &models.AuthRequest{
			Username: gofakeit.Username(),
			Password: gofakeit.Password(true, true, true, true, false, 12),
		}

		serviceErr = errors.New("service error")
	)

	type args struct {
		ctx     context.Context
		request *models.AuthRequest
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
				request: req,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.UserTokenMock.Expect(ctx, req.Username, req.Password).Return(mockToken, nil)
				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name: "wrong credentials error case",
			args: args{
				ctx:     ctx,
				request: req,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.UserTokenMock.Expect(ctx, req.Username, req.Password).Return("", service.ErrWrongCredentials)
				return mockService
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  api.ErrWrongCredentials,
		},
		{
			name: "service error case",
			args: args{
				ctx:     ctx,
				request: req,
			},
			serviceMock: func(mc *minimock.Controller) service.Service {
				mockService := mocks.NewServiceMock(mc)
				mockService.UserTokenMock.Expect(ctx, req.Username, req.Password).Return("", serviceErr)
				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  api.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody, err := json.Marshal(tt.args.request)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(reqBody))

			serviceMock := tt.serviceMock(minimock.NewController(t))
			serv := newMockHTTPServer(config.NewMockAuthConfig(), config.NewMockHttpConfig(), serviceMock)

			serv.AuthHandler(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if w.Code != tt.expectedStatus {
				var errorResp *models.ErrorResponse
				err = utils.ReadJSON(w.Body, &errorResp)
				require.NoError(t, err)
				t.Error(errorResp.Errors)
			}
		})
	}
}
