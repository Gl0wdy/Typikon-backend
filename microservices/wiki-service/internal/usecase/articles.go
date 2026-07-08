package usecase

import (
	"context"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
)

type ArticleUseCase interface {
	Create(ctx context.Context, article *domain.Article) (*domain.Article, error)
	GetByID(ctx context.Context, id string) (*domain.Article, error)
	Update(ctx context.Context, article *domain.Article) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*domain.Article, error)
}

type articleUseCase struct {
	repo         domain.ArticleRepository
	categoryRepo domain.CategoryRepository
}

func NewArticleUseCase(repo domain.ArticleRepository, categoryRepo domain.CategoryRepository) ArticleUseCase {
	return &articleUseCase{repo: repo, categoryRepo: categoryRepo}
}

func (u *articleUseCase) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	if err := article.Validate(); err != nil {
		return nil, err
	}

	if _, err := u.categoryRepo.GetByID(ctx, article.CategoryID); err != nil {
		return nil, err
	}

	createdArticle, err := u.repo.Create(ctx, article)
	if err != nil {
		return nil, err
	}
	return createdArticle, nil
}

func (u *articleUseCase) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	article, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (u *articleUseCase) Update(ctx context.Context, article *domain.Article) error {
	if err := article.Validate(); err != nil {
		return err
	}

	existing, err := u.repo.GetByID(ctx, article.ID)
	if err != nil {
		return err
	}

	article.CreatedAt = existing.CreatedAt
	return u.repo.Update(ctx, article)
}

func (u *articleUseCase) Delete(ctx context.Context, id string) error {
	if _, err := u.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}

func (u *articleUseCase) List(ctx context.Context, limit, offset int) ([]*domain.Article, error) {
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	return u.repo.List(ctx, limit, offset)
}
