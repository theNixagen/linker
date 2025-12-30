package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theNixagen/linker/internal/db"
	"github.com/theNixagen/linker/internal/domain/user"
	"github.com/theNixagen/linker/internal/repositories/user_repository"
)

type UserService struct {
	UserRepository user_repository.UserRepository
	pool           *pgxpool.Pool
	queries        *db.Queries
}

func NewUserService(userRepository user_repository.UserRepository, pool *pgxpool.Pool) *UserService {
	return &UserService{
		UserRepository: userRepository,
		pool:           pool,
		queries:        db.New(pool),
	}
}

func (us *UserService) GetUser(ctx context.Context, username string) (user.GetUser, error) {
	userFound, err := us.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return user.GetUser{}, err
	}
	userDto := user.GetUser{
		ID:             userFound.ID,
		Email:          userFound.Email,
		CreatedAt:      userFound.CreatedAt,
		ProfilePicture: userFound.ProfilePicture,
		Bio:            userFound.Bio,
		BannerPicture:  userFound.BannerPicture,
		Name:           userFound.Name,
		UserName:       userFound.Username,
	}
	return userDto, nil
}

func (us *UserService) UpdateBio(ctx context.Context, username, bio string) error {
	err := us.UserRepository.UpdateBio(ctx, username, bio)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UploadProfilePhoto(ctx context.Context, username, objectName string) error {
	err := us.UserRepository.UpdateProfilePhoto(ctx, username, objectName)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UploadBanner(ctx context.Context, username, objectName string) error {
	err := us.UserRepository.UpdateBannerPhoto(ctx, username, objectName)

	if err != nil {
		return err
	}
	return nil
}
