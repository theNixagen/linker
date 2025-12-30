package user_repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDuplicatedEmail    = errors.New("email already exists")
	ErrDuplicatedUsername = errors.New("username already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type User struct {
	ID             int32
	Email          string
	Password       string
	CreatedAt      time.Time
	ProfilePicture string
	Bio            string
	BannerPicture  string
	Name           string
	Username       string
}

type UserRepository interface {
	Create(ctx context.Context, user User) (int32, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	UpdateBio(ctx context.Context, username, bio string) error
	UpdateProfilePhoto(ctx context.Context, username, objectName string) error
	UpdateBannerPhoto(ctx context.Context, username, objectName string) error
}
