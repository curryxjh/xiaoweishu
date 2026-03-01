package repository

import (
	"context"
	"xiaoweishu/internal/domain"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

type CachedArticleRepository struct {
}

func NewArticleRepository() ArticleRepository {
	return &CachedArticleRepository{}
}

func (c *CachedArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return 1, nil
}
