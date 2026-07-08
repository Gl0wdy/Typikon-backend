package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository/models"
)

type PostgresCategoryRepo struct {
	db *gorm.DB
}

func NewPostgresCategoryRepo(db *gorm.DB) *PostgresCategoryRepo {
	return &PostgresCategoryRepo{db: db}
}

func (r *PostgresCategoryRepo) Create(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	model := toCategoryModel(category)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return toCategoryDomain(model), nil
}

func (r *PostgresCategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var model models.CategoryModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}
	return toCategoryDomain(&model), nil
}

func (r *PostgresCategoryRepo) Update(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	model := toCategoryModel(category)

	result := r.db.WithContext(ctx).Model(&models.CategoryModel{ID: category.ID}).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrCategoryNotFound
	}

	return toCategoryDomain(model), nil
}

func (r *PostgresCategoryRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.CategoryModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrCategoryNotFound
	}
	return nil
}

func (r *PostgresCategoryRepo) List(ctx context.Context, limit, offset int) ([]*domain.Category, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var dbModels []models.CategoryModel
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&dbModels).Error; err != nil {
		return nil, err
	}

	categories := make([]*domain.Category, 0, len(dbModels))
	for _, m := range dbModels {
		categories = append(categories, toCategoryDomain(&m))
	}
	return categories, nil
}

func (r *PostgresCategoryRepo) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.CategoryModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
