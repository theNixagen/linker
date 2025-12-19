package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api API) BindRoutes() {
	r := api.Router
	r.Use(middleware.Logger, middleware.AllowContentType("application/json", "multipart/form-data"), api.SetContentTypeMiddleware("application/json"))

	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.CreateUser)
		r.Post("/login", api.AuthUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Route("/profile", func(r chi.Router) {
			r.Get("/", api.GetProfile)
			r.Put("/", api.UpdateBio)
		})
	})
}
