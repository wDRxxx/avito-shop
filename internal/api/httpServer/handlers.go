package httpServer

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/models"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func (s *server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req *models.AuthRequest
	err := utils.ReadReqJSON(w, r, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "json: unknown field") {
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		slog.Error(
			"error reading request on /api/auth",
			slog.Any("error", err),
		)
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	token, err := s.service.UserToken(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			utils.WriteJSONError(err, w, http.StatusUnauthorized)
			return
		}

		slog.Error(
			"error getting user token",
			slog.Any("error", err),
		)

		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	utils.WriteJSON(&models.AuthResponse{
		Token: token,
	}, w)
}

func (s *server) BuyHandler(w http.ResponseWriter, r *http.Request) {
	itemTitle := chi.URLParam(r, "item")
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
			utils.WriteJSONError(err, w, http.StatusBadRequest)

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

func (s *server) SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := utils.VerifyToken(strings.Split(r.Header.Get("Authorization"), " ")[1], s.authConfig.TokenSecret())
	if err != nil {
		slog.Error(
			"error verifying token",
			slog.Any("error", err),
			slog.String("token", r.Header.Get("Authorization")),
		)

		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		slog.Error(
			"error converting subject to id",
			slog.Any("error", err),
			slog.String("subject", claims.Subject),
		)
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	var req *models.SendCoinRequest
	err = utils.ReadReqJSON(w, r, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "json: unknown field") {
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		slog.Error(
			"error reading request on /api/auth",
			slog.Any("error", err),
		)
		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	if claims.Username == req.ToUser {
		utils.WriteJSONError(api.ErrSendToYourself, w, http.StatusBadRequest)
		return
	}

	err = s.service.SendCoin(r.Context(), req.ToUser, userID, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) {
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		slog.Error(
			"error sending coin",
			slog.Any("error", err),
			slog.Int("from", userID),
			slog.String("to", req.ToUser),
			slog.Int("amount", req.Amount),
		)

		utils.WriteJSONError(api.ErrInternal, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
