package usecase

import (
	"context"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
)

type CategoryUseCase interface {
	Create(ctx context.Context, category *domain.Category) (*domain.Category, error)
	GetByID(ctx context.Context, id string) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) (*domain.Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*domain.Category, error)
}

type categoryUseCase struct {
	repo domain.CategoryRepository
}

func NewCategoryUseCase(repo domain.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{repo: repo}
}

func (u *categoryUseCase) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if err := category.Validate(); err != nil {
		return nil, err
	}
	return u.repo.Create(ctx, category)
}

func (u *categoryUseCase) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	category, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *categoryUseCase) Update(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if err := category.Validate(); err != nil {
		return nil, err
	}
	return u.repo.Update(ctx, category)
}

func (u *categoryUseCase) Delete(ctx context.Context, id string) error {
	_, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, id)
}

func (u *categoryUseCase) List(ctx context.Context, limit, offset int) ([]*domain.Category, error) {
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	return u.repo.List(ctx, limit, offset)
}
