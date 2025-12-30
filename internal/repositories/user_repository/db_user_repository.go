package user_repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theNixagen/linker/internal/db"
)

type DbUserRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewDbUserRepository(pool *pgxpool.Pool) *DbUserRepository {
	return &DbUserRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *DbUserRepository) Create(ctx context.Context, user User) (int32, error) {
	if _, err := r.GetUserByUsername(ctx, user.Username); err == nil {
		return 0, ErrDuplicatedUsername
	}

	result, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:    user.Email,
		Password: user.Password,
		Name:     pgtype.Text{String: user.Name, Valid: true},
		Username: pgtype.Text{String: user.Username, Valid: true},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, ErrDuplicatedEmail
			}
			return 0, errors.New("could not insert user")
		}
		return 0, err
	}

	return result, nil
}

func (r *DbUserRepository) GetUserByUsername(ctx context.Context, username string) (User, error) {
	user, err := r.queries.GetUserByUsername(ctx, pgtype.Text{
		String: username,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return User{
		ID:             user.ID,
		Email:          user.Email,
		Password:       user.Password,
		CreatedAt:      user.CreatedAt.Time,
		Name:           user.Name.String,
		Username:       user.Username.String,
		Bio:            user.Bio,
		BannerPicture:  user.BannerPicture,
		ProfilePicture: user.ProfilePicture,
	}, nil
}

func (r *DbUserRepository) UpdateBio(ctx context.Context, username, bio string) error {
	err := r.queries.UpdateBio(ctx, db.UpdateBioParams{Username: pgtype.Text{String: username, Valid: true}, Bio: bio})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (r *DbUserRepository) UpdateProfilePhoto(ctx context.Context, username, objectName string) error {
	err := r.queries.UpdateProfilePhoto(ctx, db.UpdateProfilePhotoParams{Username: pgtype.Text{String: username, Valid: true}, ProfilePicture: objectName})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *DbUserRepository) UpdateBannerPhoto(ctx context.Context, username, objectName string) error {
	return nil
}
