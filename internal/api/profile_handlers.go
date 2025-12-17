package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/theNixagen/linker/internal/services"
)

func (api *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := GetTokenClaims(r.Context())

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user unauthorized",
		})
		return
	}

	user, err := api.UserService.GetUser(r.Context(), userClaims.Email)

	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "user not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return
}
