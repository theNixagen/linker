package auth

import "github.com/golang-jwt/jwt/v5"

type RefreshClaims struct {
	Sub string `json:"sub"`
	Typ string `json:"typ"`
	Jti string `json:"jti"`
	jwt.RegisteredClaims
}
