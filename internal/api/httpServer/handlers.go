package httpServer

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/wDRxxx/avito-shop/internal/models"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func (s *server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req *models.AuthRequest
	err := utils.ReadReqJSON(w, r, &req)
	if err != nil {
		if strings.Contains(err.Error(), "json: unknown field") {
			utils.WriteJSONError(err, w, http.StatusBadRequest)
			return
		}

		slog.Error(
			"error reading request on /api/auth",
			slog.Any("error", err),
		)
		utils.WriteJSONError(err, w)
		return
	}

	token, err := s.usersService.UserToken(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrWrongCredentials) {
			utils.WriteJSONError(err, w, http.StatusUnauthorized)
			return
		}

		slog.Error(
			"error getting user token",
			slog.Any("error", err),
		)

		utils.WriteJSONError(err, w)
		return
	}

	utils.WriteJSON(&models.AuthResponse{
		Token: token,
	}, w)
}
