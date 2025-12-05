package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api API) BindRoutes() {
	r := api.Router
	r.Use(middleware.Logger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.CreateUser)
	})
}
