package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository/models"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewPostgresConnection(cfg DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := db.AutoMigrate(&models.ArticleModel{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}
