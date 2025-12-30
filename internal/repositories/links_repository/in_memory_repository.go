package links_repository

import (
	"context"
)

type InMemoryLinksRepository struct {
	Links []Link
}

func NewInMemoryLinksRepository() *InMemoryLinksRepository {
	return &InMemoryLinksRepository{
		Links: []Link{},
	}
}

func (r *InMemoryLinksRepository) CreateLink(ctx context.Context, UserID int32, Url, Title, Description string) error {
	r.Links = append(r.Links, Link{
		ID:          int32(len(r.Links) + 1),
		UserID:      UserID,
		Url:         Url,
		Title:       Title,
		Description: Description,
	})

	return nil
}

func (r *InMemoryLinksRepository) FindAllLinksFromAUser(ctx context.Context, userID int32) ([]Link, error) {
	var result []Link
	for _, link := range r.Links {
		if link.UserID != userID {
			continue
		}
		result = append(result, Link{
			ID:          link.ID,
			UserID:      link.UserID,
			Url:         link.Url,
			Title:       link.Title,
			Description: link.Description,
			CreatedAt:   link.CreatedAt,
		})
	}

	if len(result) == 0 {
		return nil, ErrLinksNotFound
	}

	return result, nil
}
