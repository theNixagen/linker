package user

import (
	"time"
)

type GetUser struct {
	ID             int32     `json:"id"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	ProfilePicture string    `json:"profile_picture"`
	Bio            string    `json:"bio"`
	BannerPicture  string    `json:"banner_picture"`
	Name           string    `json:"name"`
}
