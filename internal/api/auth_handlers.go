package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/theNixagen/linker/internal/domain/user"
	"github.com/theNixagen/linker/internal/repositories/user_repository"
	"github.com/theNixagen/linker/internal/services"
)

// CreateUser godoc
// @Summary      Cria um novo usuario
// @Tags         auth
// @Produce      json
// @Accept       json
// @Param        request  body  user.CreateUser  true  "payload"
// @Success      201  {object}  map[string]int
// @Failure      409  {object}  map[string]any
// @Failure      500  {object}  nil
// @Failure      400  {object}  map[string]any
// @Failure      422  {object}  map[string]string
// @Router       /users [post]
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

	id, err := api.AuthService.CreateUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, user_repository.ErrDuplicatedEmail) || errors.Is(err, user_repository.ErrDuplicatedUsername) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error(),
			})
			return
		}
		println(err.Error())
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

// login godoc
// @Summary      autentica um usuário
// @Tags         auth
// @Produce      json
// @Accept       json
// @Param        request  body  user.AuthUser  true  "payload"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  nil
// @Failure      401  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      422  {object}  map[string]string
// @Router       /users/login [post]
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

	token, refreshToken, err := api.AuthService.AuthUser(r.Context(), user.Username, user.Password)

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

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/users/refresh-session",
		MaxAge:   24 * 60 * 60,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (api *API) RefreshSession(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, refreshToken, err := api.AuthService.RefreshSession(r.Context(), refreshTokenCookie.Value)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/users/refresh-session",
		MaxAge:   24 * 60 * 60,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
