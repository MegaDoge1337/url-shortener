package service

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/megadoge1337/url-shortener/internal/domain"
	"github.com/megadoge1337/url-shortener/internal/repository"
	"github.com/megadoge1337/url-shortener/models"
	"github.com/redis/go-redis/v9"
)

type UrlService struct {
	repo  *repository.UrlRepository
	cache *redis.Client
}

func NewUrlService(repo *repository.UrlRepository, cache *redis.Client) *UrlService {
	return &UrlService{
		repo:  repo,
		cache: cache,
	}
}

func (s *UrlService) Create(u domain.URL) (*domain.URL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	raw, err := s.cache.Get(ctx, normalizeURL(u.URL)).Bytes()
	if err == nil {
		var cached domain.URL
		if jsonErr := json.Unmarshal(raw, &cached); jsonErr == nil {
			return &cached, nil
		}
	}

	var urlModel models.URL
	urlModel.ID = u.ID
	urlModel.URL = normalizeURL(u.URL)
	urlModel.Alias = u.Alias

	newUrl, err := s.repo.Create(ctx, &urlModel)
	if err != nil {
		return nil, err
	}

	u.ID = newUrl.ID
	u.URL = newUrl.URL
	u.Alias = newUrl.Alias

	uJson, err := json.Marshal(u)
	if err == nil {
		s.cache.Set(ctx, u.URL, uJson, 1*time.Minute)
	}

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

	return &url, nil
}

func normalizeURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" {
		return "https://" + rawURL
	}
	return rawURL
}
