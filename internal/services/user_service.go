package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theNixagen/linker/internal/db"
	"github.com/theNixagen/linker/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedEmail    = errors.New("this email is already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type UserService struct {
	pool      *pgxpool.Pool
	queries   *db.Queries
	jwtSecret string
}

func NewUserService(pool *pgxpool.Pool, jwtSecret string) UserService {
	return UserService{
		pool:      pool,
		queries:   db.New(pool),
		jwtSecret: jwtSecret,
	}
}

func (us *UserService) CreateUser(ctx context.Context, user user.CreateUser) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 6)
	if err != nil {
		return 0, err
	}

	id, err := us.queries.CreateUser(ctx, db.CreateUserParams{
		Email:    user.Email,
		Name:     pgtype.Text{String: user.Name, Valid: true},
		Password: string(hashedPassword),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ColumnName == "email" && pgErr.Code == "23505" {
				return 0, ErrDuplicatedEmail
			}
			return 0, errors.New("could not insert user")
		}
		return 0, err
	}

	return int(id), nil
}

func (us *UserService) AuthUser(ctx context.Context, email, password string) (string, error) {
	user, err := us.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(us.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *UserService) GetUser(ctx context.Context, email string) (user.GetUser, error) {
	userFound, err := us.queries.GetUserByEmail(ctx, email)
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
	}

	return userDto, nil
}

func (us *UserService) UpdateBio(ctx context.Context, email, bio string) error {
	err := us.queries.UpdateBio(ctx, db.UpdateBioParams{
		Email: email,
		Bio:   bio,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}

func (us *UserService) UploadProfilePhoto(ctx context.Context, email, objectName string) error {
	err := us.queries.UpdateProfilePhoto(ctx, db.UpdateProfilePhotoParams{
		ProfilePicture: objectName,
		Email:          email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (us *UserService) UploadBanner(ctx context.Context, email, objectName string) error {
	err := us.queries.UpdateBannerPhoto(ctx, db.UpdateBannerPhotoParams{
		BannerPicture: objectName,
		Email:         email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
