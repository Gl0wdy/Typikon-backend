// internal/repository/mapper.go
package repository

import (
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository/models"
)

func toModel(a *domain.Article) *models.ArticleModel {
	return &models.ArticleModel{
		ID:         a.ID,
		Title:      a.Title,
		Content:    a.Content,
		CategoryID: a.CategoryID,
		UserID:     a.UserID,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}

func toDomain(m *models.ArticleModel) *domain.Article {
	return &domain.Article{
		ID:         m.ID,
		Title:      m.Title,
		Content:    m.Content,
		CategoryID: m.CategoryID,
		UserID:     m.UserID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
