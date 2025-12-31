package services

import (
	"testing"

	"github.com/theNixagen/linker/internal/domain/user"
	"github.com/theNixagen/linker/internal/repositories/cache_repository"
	"github.com/theNixagen/linker/internal/repositories/user_repository"
)

func TestAuthService_CreateUser(t *testing.T) {
	jwtSecret := "test_jwt"
	refreshSecret := "test_refresh"
	userRepository := user_repository.NewInMemoryUserRepository()
	cacheRepository := cache_repository.NewInMemoryCacheRepository()
	as := NewAuthService(jwtSecret, refreshSecret, userRepository, cacheRepository)
	as.CreateUser(t.Context(), user.CreateUser{
		Name:     "john doe",
		Username: "johndoe",
		Email:    "johndoe@example.com",
		Password: "johndoesupersecretpassword",
	})

	if len(userRepository.Users) != 1 {
		t.Errorf("expected 1 user, got %d", len(userRepository.Users))
	}
}

func TestAuthService_AuthUser(t *testing.T) {
	jwtSecret := "test_jwt"
	refreshSecret := "test_refresh"
	userRepository := user_repository.NewInMemoryUserRepository()
	cacheRepository := cache_repository.NewInMemoryCacheRepository()
	as := NewAuthService(jwtSecret, refreshSecret, userRepository, cacheRepository)
	as.CreateUser(t.Context(), user.CreateUser{
		Name:     "john doe",
		Username: "johndoe",
		Email:    "johndoe@example.com",
		Password: "johndoesupersecretpassword",
	})

	token, _, err := as.AuthUser(t.Context(), "johndoe", "johndoesupersecretpassword")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if token == "" {
		t.Errorf("expected token, got empty string")
	}
}

func TestAuthService_AuthUser_WrongUsername(t *testing.T) {
	jwtSecret := "test_jwt"
	refreshSecret := "test_refresh"
	userRepository := user_repository.NewInMemoryUserRepository()
	cacheRepository := cache_repository.NewInMemoryCacheRepository()
	as := NewAuthService(jwtSecret, refreshSecret, userRepository, cacheRepository)
	as.CreateUser(t.Context(), user.CreateUser{
		Name:     "john doe",
		Username: "johndoe",
		Email:    "johndoe@example.com",
		Password: "johndoesupersecretpassword",
	})

	_, _, err := as.AuthUser(t.Context(), "john", "johndoesupersecretpassword")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestAuthService_AuthUser_WrongPassword(t *testing.T) {
	jwtSecret := "test_jwt"
	refreshSecret := "test_refresh"
	userRepository := user_repository.NewInMemoryUserRepository()
	cacheRepository := cache_repository.NewInMemoryCacheRepository()
	as := NewAuthService(jwtSecret, refreshSecret, userRepository, cacheRepository)
	as.CreateUser(t.Context(), user.CreateUser{
		Name:     "john doe",
		Username: "johndoe",
		Email:    "johndoe@example.com",
		Password: "johndoesupersecretpassword",
	})

	_, _, err := as.AuthUser(t.Context(), "johndoe", "wrongpassword")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
