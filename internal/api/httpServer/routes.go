package httpServer

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *server) setRoutes() {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

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
