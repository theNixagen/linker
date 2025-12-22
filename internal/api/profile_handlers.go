package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/theNixagen/linker/internal/db"
	"github.com/theNixagen/linker/internal/domain/links"
	"github.com/theNixagen/linker/internal/domain/user"
	"github.com/theNixagen/linker/internal/services"
)

func (api *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "username")

	profile, err := api.UserService.GetUser(r.Context(), user)
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
	json.NewEncoder(w).Encode(profile)
}

func (api *API) UpdateBio(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := GetTokenClaims(r.Context())

	var req user.UpdateBioRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "unprocessable entity",
		})
		return
	}

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user unauthorized",
		})
		return
	}

	if err := api.UserService.UpdateBio(r.Context(), userClaims.Username, req.Bio); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "user not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user not found",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetTokenClaims(r.Context())

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	info, err := api.FileService.PutObject(r.Context(), header.Filename, file, header.Size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = api.UserService.UploadProfilePhoto(r.Context(), claims.Username, info.Key)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) UploadBanner(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetTokenClaims(r.Context())

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	info, err := api.FileService.PutObject(r.Context(), header.Filename, file, header.Size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = api.UserService.UploadBanner(r.Context(), claims.Username, info.Key)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *API) CreateNewLink(w http.ResponseWriter, r *http.Request) {
	var link links.CreateLink

	claims, ok := GetTokenClaims(r.Context())

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := api.Validator.Struct(link); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	err := api.UserService.CreateLink(r.Context(), claims.Username, link.URL, link.Title, link.Description)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (api *API) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	links, err := api.UserService.GetAllLinksFromAUser(r.Context(), username)
	if err != nil {
		if errors.Is(err, services.ErrLinksNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]db.Link{
		"links": links,
	})
}
