package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository/models"
)

type PostgresArticleRepo struct {
	db *gorm.DB
}

func NewPostgresArticleRepo(db *gorm.DB) *PostgresArticleRepo {
	return &PostgresArticleRepo{db: db}
}

func (r *PostgresArticleRepo) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	model := toArticleModel(article)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return toArticleDomain(model), nil
}

func (r *PostgresArticleRepo) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	var model models.ArticleModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrArticleNotFound
		}
		return nil, err
	}

	return toArticleDomain(&model), nil
}

func (r *PostgresArticleRepo) Update(ctx context.Context, article *domain.Article) error {
	model := toArticleModel(article)

	result := r.db.WithContext(ctx).Model(&models.ArticleModel{}).
		Where("id = ?", article.ID).
		Updates(model)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrArticleNotFound
	}
	return nil
}

func (r *PostgresArticleRepo) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.ArticleModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrArticleNotFound
	}
	return nil
}

func (r *PostgresArticleRepo) List(ctx context.Context, limit, offset int) ([]*domain.Article, error) {
	var dbModels []models.ArticleModel

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&dbModels).Error; err != nil {
		return nil, err
	}

	articles := make([]*domain.Article, 0, len(dbModels))
	for _, m := range dbModels {
		articles = append(articles, toArticleDomain(&m))
	}
	return articles, nil
}
