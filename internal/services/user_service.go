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
		"id":             user.ID,
		"name":           user.Name,
		"email":          user.Email,
		"bio":            user.Bio,
		"profilePicture": user.ProfilePicture,
		"bannerPicture":  user.BannerPicture,
		"exp":            time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(us.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
