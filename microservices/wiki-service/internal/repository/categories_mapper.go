package repository

import (
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository/models"
)

func toCategoryModel(c *domain.Category) *models.CategoryModel {
	return &models.CategoryModel{
		ID:        c.ID,
		Title:     c.Title,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryDomain(m *models.CategoryModel) *domain.Category {
	return &domain.Category{
		ID:        m.ID,
		Title:     m.Title,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
