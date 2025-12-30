package links_repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrLinksNotFound = errors.New("links not found")
)

type Link struct {
	ID          int32
	UserID      int32
	Url         string
	Title       string
	Description string
	CreatedAt   time.Time
}

type LinksRepository interface {
	CreateLink(ctx context.Context, UserID int32, Url, Title, Description string) error
	FindAllLinksFromAUser(ctx context.Context, userID int32) ([]Link, error)
}
