package api

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *API) BindRoutes() {
	r := api.Router
	r.Use(middleware.Logger, middleware.AllowContentType("application/json", "multipart/form-data"), api.SetContentTypeMiddleware("application/json"))

	r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	r.Get("/reference", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "http://localhost:8080/docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Linker API",
			},
			DarkMode: true,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlContent))
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.CreateUser)
		r.Post("/login", api.AuthUser)
		r.Post("/refresh-session", api.RefreshSession)
	})

	r.Route("/profile", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(api.AuthMiddleware)
			r.Put("/", api.UpdateBio)
			r.Put("/photo", api.UploadProfilePicture)
			r.Put("/banner", api.UploadBanner)
			r.Post("/link", api.CreateNewLink)
		})
		r.Get("/{username}", api.GetProfile)
		r.Get("/links/{username}", api.GetUserLinks)
	})
}
