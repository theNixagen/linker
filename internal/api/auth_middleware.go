package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/theNixagen/linker/internal/domain/user"
)

type contextKey string

const TokenClaimsKey contextKey = "tokenClaims"

func GetTokenClaims(ctx context.Context) (user.UserClaims, bool) {
	claims, ok := ctx.Value(TokenClaimsKey).(user.UserClaims)
	return claims, ok
}

func verifyToken(tokenString string, secretKey string) (user.UserClaims, error) {
	claims := &user.UserClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return user.UserClaims{}, err
	}

	if token.Valid {
		return *claims, nil
	}

	return user.UserClaims{}, fmt.Errorf("invalid token")
}

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			encoder.Encode(map[string]string{
				"message": "Unauthorized",
			})
			return
		}

		_, token, _ := strings.Cut(bearerToken, " ")

		claims, err := verifyToken(token, api.JwtSecret)

		if err != nil {
			fmt.Println(err)
			encoder.Encode(map[string]string{
				"message": "Invalid token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), TokenClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
