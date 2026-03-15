package service

import (
	"context"
	"time"

	"github.com/megadoge1337/url-shortener/internal/domain"
	"github.com/megadoge1337/url-shortener/internal/repository"
	"github.com/megadoge1337/url-shortener/models"
)

type UrlService struct {
	repo *repository.UrlRepository
}

func NewUrlService(repo *repository.UrlRepository) *UrlService {
	return &UrlService{
		repo: repo,
	}
}

func (s *UrlService) Create(u domain.URL) (*domain.URL, error) {
	var urlModel models.URL
	urlModel.ID = u.ID
	urlModel.URL = u.URL
	urlModel.Alias = u.Alias

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUrl, err := s.repo.Create(ctx, &urlModel)
	if err != nil {
		return nil, err
	}

	u.ID = newUrl.ID
	u.URL = newUrl.URL
	u.Alias = newUrl.Alias
	return &u, nil
}

func (s *UrlService) GetByAlias(alias string) (*domain.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	urlModel, err := s.repo.GetByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}

	url := domain.URL{
		ID:    urlModel.ID,
		URL:   urlModel.URL,
		Alias: urlModel.Alias,
	}

	return &url, err
}
