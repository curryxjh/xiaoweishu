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
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	} else {
		return a.repo.Create(ctx, article)
	}
}
