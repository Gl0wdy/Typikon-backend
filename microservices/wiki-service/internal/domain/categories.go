package domain

import (
	"context"
	"errors"
	"time"
)

type Category struct {
	ID        string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidCategoryTitle = errors.New("invalid category title: cannot be empty")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryTitleTooLong = errors.New("invalid category title: too long")
)

func (c *Category) Validate() error {
	if c.Title == "" {
		return ErrInvalidCategoryTitle
	}
	if len(c.Title) > 100 {
		return ErrCategoryTitleTooLong
	}
	return nil
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	GetByID(ctx context.Context, id string) (*Category, error)
	Exists(ctx context.Context, id string) (bool, error)
	Update(ctx context.Context, category *Category) (*Category, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Category, error)
}
