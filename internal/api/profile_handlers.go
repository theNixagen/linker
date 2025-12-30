package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/theNixagen/linker/internal/domain/links"
	"github.com/theNixagen/linker/internal/domain/user"
	"github.com/theNixagen/linker/internal/repositories/links_repository"
	"github.com/theNixagen/linker/internal/repositories/user_repository"
)

// GetProfile godoc
// @Summary      Busca um perfil pelo nome de usuario
// @Tags         profile
// @Produce      json
// @Param        username  path      string  true  "username"
// @Success      200  {object}  user.GetUser
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /profile/{username} [get]
func (api *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "username")

	profile, err := api.UserService.GetUser(r.Context(), user)
	if err != nil {
		if errors.Is(err, user_repository.ErrUserNotFound) {
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

	profileURL, err := api.FileService.GetSignedURL(r.Context(), profile.ProfilePicture, api.FileService.BucketName)
	if err != nil {
		profile.ProfilePicture = ""
	}

	bannerURL, err := api.FileService.GetSignedURL(r.Context(), profile.BannerPicture, api.FileService.BucketName)
	if err != nil {
		profile.BannerPicture = ""
	}

	if profileURL != nil {
		profile.ProfilePicture = profileURL.String()
	}

	if bannerURL != nil {
		profile.BannerPicture = bannerURL.String()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// UpdateBio godoc
// @Summary      Atualiza a bio de um perfil
// @Tags         profile
// @Produce      json
// @Accept       json
// @Param        request  body  user.UpdateBioRequest  true  "payload"
// @Security BearerAuth
// @Success      204  {object}  nil
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      422  {object}  map[string]string
// @Router       /profile/bio [put]
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
		if errors.Is(err, user_repository.ErrUserNotFound) {
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

// UpdatePhoto godoc
// @Summary      Atualiza a foto de perfil
// @Tags         profile
// @Produce      json
// @Accept       multipart/form-data
// @Param        photo  formData  file  true  "Imagem (jpg/png/webp)"
// @Security BearerAuth
// @Success      204  {object}  nil
// @Failure      404  {object}  nil
// @Failure      500  {object}  nil
// @Failure      401  {object}  nil
// @Failure      422  {object}  nil
// @Router       /profile/photo [put]
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
		if errors.Is(err, user_repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateBanner godoc
// @Summary      Atualiza a foto de banner
// @Tags         profile
// @Produce      json
// @Accept       multipart/form-data
// @Param        photo  formData  file  true  "Imagem (jpg/png/webp)"
// @Security BearerAuth
// @Success      204  {object}  nil
// @Failure      404  {object}  nil
// @Failure      500  {object}  nil
// @Failure      401  {object}  nil
// @Failure      422  {object}  nil
// @Router       /profile/banner [put]
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
		if errors.Is(err, user_repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateLink godoc
// @Summary      Cria um novo link para um usuário autenticado
// @Tags         profile
// @Produce      json
// @Accept       json
// @Param        request  body  links.CreateLink  true  "payload"
// @Security BearerAuth
// @Success      204  {object}  nil
// @Failure      404  {object}  nil
// @Failure      500  {object}  nil
// @Failure      401  {object}  nil
// @Failure      422  {object}  nil
// @Router       /profile/link [post]
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

	err := api.LinksService.CreateLink(r.Context(), claims.Username, link.URL, link.Title, link.Description)
	if err != nil {
		if errors.Is(err, user_repository.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetUserLinks godoc
// @Summary      Busca os links de um usuario
// @Tags         profile
// @Produce      json
// @Param        username  path      string  true  "username"
// @Success      200  {object}  user.GetUser
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /profile/links/{username} [get]
func (api *API) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	links, err := api.LinksService.GetAllLinksFromAUser(r.Context(), username)
	if err != nil {
		if errors.Is(err, links_repository.ErrLinksNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"links": links,
	})
}
