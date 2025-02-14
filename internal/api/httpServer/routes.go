package httpServer

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *server) setRoutes() {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Route("/api", func(mux chi.Router) {
		mux.Post("/auth", s.AuthHandler)

		mux.Group(func(mux chi.Router) {
			mux.Use(s.authRequiredMiddleware)

			mux.Get("/buy/{item}", s.BuyHandler)
			mux.Post("/sendCoin", s.SendCoinHandler)
			mux.Get("/info", s.UserInfoHandler)
		})
	})

	s.mux = mux
}
