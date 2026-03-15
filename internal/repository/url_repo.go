package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/megadoge1337/url-shortener/models"
)

type UrlRepository struct {
	db *sql.DB
}

func NewUrlRepository(db *sql.DB) *UrlRepository {
	return &UrlRepository{
		db: db,
	}
}

func (r *UrlRepository) Create(ctx context.Context, u *models.URL) (*models.URL, error) {
	err := u.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UrlRepository) GetById(ctx context.Context, id int) (*models.URL, error) {
	url, err := models.FindURL(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (r *UrlRepository) GetByAlias(ctx context.Context, alias string) (*models.URL, error) {
	url, err := models.Urls(
		qm.Where("alias = ?", alias),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return url, nil
}
