package services

import (
	"testing"

	"github.com/theNixagen/linker/internal/repositories/user_repository"
)

func TestUserService_GetUser(t *testing.T) {
	ur := user_repository.NewInMemoryUserRepository()
	us := NewUserService(ur)
	ur.Create(t.Context(), user_repository.User{
		ID:       1,
		Email:    "johndoe@example.com",
		Name:     "John Doe",
		Username: "johndoe",
		Password: "jonesdoe",
	})

	user, err := us.GetUser(t.Context(), "johndoe")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.UserName != "johndoe" {
		t.Fatalf("expected username 'johndoe', got %v", user.UserName)
	}
}

func TestUserService_UpdateBio(t *testing.T) {
	ur := user_repository.NewInMemoryUserRepository()
	us := NewUserService(ur)
	ur.Create(t.Context(), user_repository.User{
		ID:       1,
		Email:    "johndoe@example.com",
		Name:     "John Doe",
		Username: "johndoe",
		Password: "jonesdoe",
	})

	us.UpdateBio(t.Context(), "johndoe", "bio")

	user, err := us.GetUser(t.Context(), "johndoe")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Bio != "bio" {
		t.Fatalf("expected bio 'bio', got %v", user.Bio)
	}
}

func TestUserService_UploadProfilePhoto(t *testing.T) {
	ur := user_repository.NewInMemoryUserRepository()
	us := NewUserService(ur)
	ur.Create(t.Context(), user_repository.User{
		ID:       1,
		Email:    "johndoe@example.com",
		Name:     "John Doe",
		Username: "johndoe",
		Password: "jonesdoe",
	})

	us.UploadProfilePhoto(t.Context(), "johndoe", "fototeste.jpg")

	user, err := us.GetUser(t.Context(), "johndoe")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ProfilePicture != "fototeste.jpg" {
		t.Fatalf("expected profile picture 'fototeste.jpg', got %v", user.ProfilePicture)
	}
}

func TestUserService_UploadBanner(t *testing.T) {
	ur := user_repository.NewInMemoryUserRepository()
	us := NewUserService(ur)
	ur.Create(t.Context(), user_repository.User{
		ID:       1,
		Email:    "johndoe@example.com",
		Name:     "John Doe",
		Username: "johndoe",
		Password: "jonesdoe",
	})

	us.UploadBanner(t.Context(), "johndoe", "fototeste.jpg")

	user, err := us.GetUser(t.Context(), "johndoe")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.BannerPicture != "fototeste.jpg" {
		t.Fatalf("expected banner picture 'fototeste.jpg', got %v", user.BannerPicture)
	}
}
