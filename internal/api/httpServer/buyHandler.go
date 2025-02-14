package httpServer

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func (s *server) BuyHandler(w http.ResponseWriter, r *http.Request) {
	// Was removed, cause doesn't work with testing request
	// itemTitle := chi.URLParam(r, "item")

	exploded := strings.Split(r.RequestURI, "/")
	itemTitle := exploded[3]

	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	userID, err := utils.UserIDFromToken(token, s.authConfig.TokenSecret())
	if err != nil {
		slog.Error(
			"error verifying token",
			slog.Any("error", err),
			slog.String("token", r.Header.Get("Authorization")),
		)

		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	err = s.service.BuyItem(r.Context(), userID, itemTitle)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			utils.WriteJSONError(api.ErrItemNotFound, w, http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrInsufficientBalance) {
			utils.WriteJSONError(api.ErrInsufficientBalance, w, http.StatusBadRequest)
			return
		}

		slog.Error(
			"error buying item",
			slog.Any("error", err),
			slog.String("item", itemTitle),
			slog.Int("userID", userID),
		)

		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
