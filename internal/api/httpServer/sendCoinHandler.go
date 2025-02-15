package httpServer

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/models"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func (s *server) SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	claims, err := utils.VerifyToken(token, s.authConfig.TokenSecret())
	if err != nil {
		slog.Error(
			"error verifying token",
			slog.Any("error", err),
			slog.String("token", r.Header.Get("Authorization")),
		)

		_ = utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error(
			"error converting subject to id",
			slog.Any("error", err),
			slog.String("subject", claims.Subject),
		)
		_ = utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	var req *models.SendCoinRequest
	err = utils.ReadReqJSON(w, r, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_ = utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "json:") {
			_ = utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		slog.Error(
			"error reading request on /api/auth",
			slog.Any("error", err),
		)
		_ = utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	if claims.Username == req.ToUser {
		_ = utils.WriteJSONError(api.ErrSendToYourself, w, http.StatusBadRequest)
		return
	}

	err = s.service.SendCoin(r.Context(), req.ToUser, userID, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) {
			_ = utils.WriteJSONError(api.ErrInsufficientBalance, w, http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			_ = utils.WriteJSONError(api.ErrUserNotFound, w, http.StatusBadRequest)
			return
		}

		slog.Error(
			"error sending coin",
			slog.Any("error", err),
			slog.Int("from", userID),
			slog.String("to", req.ToUser),
			slog.Int("amount", req.Amount),
		)

		_ = utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
