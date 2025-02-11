package httpServer

import (
	"net/http"

	"github.com/pkg/errors"
)

var errInvalidAuthHeader = errors.New("invalid oauth header")

func (s *server) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Vary", "Authorization")

		_, _, err := s.getAndVerifyHeaderToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
