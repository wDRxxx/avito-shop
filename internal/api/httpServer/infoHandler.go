package httpServer

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/models/converter"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func (s *server) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	userID, err := utils.UserIDFromToken(token, s.authConfig.TokenSecret())
	if err != nil {
		slog.Error(
			"error verifying token",
			slog.Any("error", err),
			slog.String("token", r.Header.Get("Authorization")),
		)

		_ = utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	info, err := s.service.UserInfo(r.Context(), userID)
	if err != nil {
		slog.Error(
			"error getting user info",
			slog.Any("error", err),
			slog.Int("userID", userID),
		)
		_ = utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	resp := converter.UserInfoFromServiceToApi(info)
	_ = utils.WriteJSON(&resp, w)
}
