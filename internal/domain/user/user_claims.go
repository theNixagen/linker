package user

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
	BannerPicture  string `json:"banner_picture"`
	jwt.RegisteredClaims
}
