package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/theNixagen/linker/internal/domain/user"
	"github.com/theNixagen/linker/internal/services"
)

func (api *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user user.CreateUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "could not decode json",
		})
		return
	}

	if err := api.Validator.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmail) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "email already exists",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "internal server error",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{
		"id": id,
	})
}

func (api *API) AuthUser(w http.ResponseWriter, r *http.Request) {
	var user user.AuthUser

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "could not decode json",
		})
		return
	}

	if err := api.Validator.Struct(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	token, err := api.UserService.AuthUser(r.Context(), user.Username, user.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid credentials",
			})
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Unauthorized",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
