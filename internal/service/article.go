package service

import (
	"context"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}
func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	return a.repo.Create(ctx, article)
}
