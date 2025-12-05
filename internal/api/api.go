package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/theNixagen/linker/internal/services"
)

type API struct {
	Router      *chi.Mux
	Validator   *validator.Validate
	UserService services.UserService
}
