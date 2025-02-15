package httpServer

import (
	"net/http"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func (s *server) authRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Vary", "Authorization")

		_, _, err := s.getAndVerifyHeaderToken(r)
		if err != nil {
			utils.WriteJSONError(api.ErrUnauthorized, w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
