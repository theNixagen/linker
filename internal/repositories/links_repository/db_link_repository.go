package links_repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/theNixagen/linker/internal/db"
)

type DbLinksRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewDbLinksRepository(pool *pgxpool.Pool) *DbLinksRepository {
	return &DbLinksRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *DbLinksRepository) CreateLink(ctx context.Context, UserID int32, Url, Title, Description string) error {
	err := r.queries.CreateLink(ctx, db.CreateLinkParams{
		UserID:      UserID,
		Url:         Url,
		Title:       Title,
		Description: Description,
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *DbLinksRepository) FindAllLinksFromAUser(ctx context.Context, userID int32) ([]Link, error) {
	links, err := r.queries.FindAllLinksFromAUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLinksNotFound
		}
		return nil, err
	}

	if len(links) == 0 {
		return nil, ErrLinksNotFound
	}

	var result []Link
	for _, link := range links {
		result = append(result, Link{
			ID:          link.ID,
			UserID:      link.UserID,
			Url:         link.Url,
			Title:       link.Title,
			Description: link.Description,
			CreatedAt:   link.CreatedAt.Time,
		})
	}

	return result, nil
}
