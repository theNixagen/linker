package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theNixagen/linker/internal/db"
	"github.com/theNixagen/linker/internal/domain/auth"
	"github.com/theNixagen/linker/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedEmail    = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	pool          *pgxpool.Pool
	queries       *db.Queries
	redisService  *RedisService
	jwtSecret     string
	refreshSecret string
}

func NewAuthService(pool *pgxpool.Pool, redisAddr, jwtSecret, refreshSecret string) *AuthService {
	return &AuthService{
		pool:          pool,
		queries:       db.New(pool),
		redisService:  NewRedisService(redisAddr),
		jwtSecret:     jwtSecret,
		refreshSecret: refreshSecret,
	}
}

func (as *AuthService) CreateUser(ctx context.Context, user user.CreateUser) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 6)

	if err != nil {
		return 0, err
	}

	id, err := as.queries.CreateUser(ctx, db.CreateUserParams{
		Email:    user.Email,
		Name:     pgtype.Text{String: user.Name, Valid: true},
		Username: pgtype.Text{String: user.Username, Valid: true},
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

func (as *AuthService) generateToken(ctx context.Context, user db.User) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	uuid := uuid.NewString()

	as.redisService.Set(ctx, fmt.Sprintf("uuid:%s", user.Username.String), uuid, time.Hour*24)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"typ": "refresh",
		"jti": uuid,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(as.refreshSecret))

	if err != nil {
		return "", "", err
	}

	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, err
}

func (as *AuthService) AuthUser(ctx context.Context, username, password string) (string, string, error) {
	user, err := as.queries.GetUserByUsername(ctx, pgtype.Text{
		String: username,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	tokenString, refreshTokenString, err := as.generateToken(ctx, user)

	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func (as *AuthService) RefreshSession(ctx context.Context, tokenStr string) (string, string, error) {
	claims := &auth.RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(as.refreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("token inválido")
	}

	if claims.Typ != "refresh" {
		return "", "", errors.New("token não é refresh")
	}

	if claims.Sub == "" {
		return "", "", errors.New("invalid token claims")
	}

	uuid, err := as.redisService.Get(ctx, fmt.Sprintf("uuid:%s", claims.Sub))

	if err != nil {
		return "", "", errors.New("invalid token uuid")
	}

	if uuid != claims.Jti {
		return "", "", errors.New("token uuid does not match")
	}

	user, err := as.queries.GetUserByUsername(ctx, pgtype.Text{
		String: claims.Sub,
		Valid:  true,
	})

	tokenString, refreshTokenString, err := as.generateToken(ctx, user)

	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}
