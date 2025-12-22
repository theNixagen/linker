package user

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
