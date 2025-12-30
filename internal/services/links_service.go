package services

import (
	"context"

	"github.com/theNixagen/linker/internal/repositories/links_repository"
	"github.com/theNixagen/linker/internal/repositories/user_repository"
)

type LinksService struct {
	UserRepository  user_repository.UserRepository
	LinksRepository links_repository.LinksRepository
}

func NewLinksService(userRepository user_repository.UserRepository, linksRepository links_repository.LinksRepository) *LinksService {
	return &LinksService{
		UserRepository:  userRepository,
		LinksRepository: linksRepository,
	}
}

func (ls *LinksService) CreateLink(ctx context.Context, username, url, title, description string) error {
	user, err := ls.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	if err = ls.LinksRepository.CreateLink(ctx, user.ID, url, title, description); err != nil {
		return err
	}

	return nil
}

func (ls *LinksService) GetAllLinksFromAUser(ctx context.Context, username string) ([]links_repository.Link, error) {
	user, err := ls.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	links, err := ls.LinksRepository.FindAllLinksFromAUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return links, nil
}
