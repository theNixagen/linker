package services

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theNixagen/linker/internal/db"
	"github.com/theNixagen/linker/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedEmail = errors.New("this email is already used")
)

type UserService struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool:    pool,
		queries: db.New(pool),
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
