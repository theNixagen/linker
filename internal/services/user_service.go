package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theNixagen/linker/internal/db"
	"github.com/theNixagen/linker/internal/domain/user"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrLinksNotFound = errors.New("links not found")
)

type UserService struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (us *UserService) GetUser(ctx context.Context, username string) (user.GetUser, error) {
	userFound, err := us.queries.GetUserByUsername(ctx, pgtype.Text{
		String: username,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.GetUser{}, ErrUserNotFound
		}

		return user.GetUser{}, err
	}

	userDto := user.GetUser{
		ID:             userFound.ID,
		Email:          userFound.Email,
		CreatedAt:      userFound.CreatedAt.Time,
		ProfilePicture: userFound.ProfilePicture,
		Bio:            userFound.Bio,
		BannerPicture:  userFound.BannerPicture,
		Name:           userFound.Name.String,
		UserName:       userFound.Username.String,
	}

	return userDto, nil
}

func (us *UserService) UpdateBio(ctx context.Context, username, bio string) error {
	err := us.queries.UpdateBio(ctx, db.UpdateBioParams{
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
		Bio: bio,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (us *UserService) UploadProfilePhoto(ctx context.Context, username, objectName string) error {
	err := us.queries.UpdateProfilePhoto(ctx, db.UpdateProfilePhotoParams{
		ProfilePicture: objectName,
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (us *UserService) UploadBanner(ctx context.Context, username, objectName string) error {
	err := us.queries.UpdateBannerPhoto(ctx, db.UpdateBannerPhotoParams{
		BannerPicture: objectName,
		Username: pgtype.Text{
			String: username,
			Valid:  true,
		},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (us *UserService) CreateLink(ctx context.Context, username, url, title, description string) error {
	user, err := us.queries.GetUserByUsername(ctx, pgtype.Text{
		String: username,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if err = us.queries.CreateLink(ctx, db.CreateLinkParams{
		UserID:      user.ID,
		Url:         url,
		Title:       title,
		Description: description,
	}); err != nil {
		return nil
	}

	return nil
}

func (us *UserService) GetAllLinksFromAUser(ctx context.Context, username string) ([]db.Link, error) {
	user, err := us.queries.GetUserByUsername(ctx, pgtype.Text{
		String: username,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	links, err := us.queries.FindAllLinksFromAUser(ctx, user.ID)
	if err != nil {
		if errors.Is(err, ErrLinksNotFound) {
			return nil, ErrLinksNotFound
		}
		return nil, err
	}

	return links, nil
}
