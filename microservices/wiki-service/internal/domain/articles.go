package domain

import (
	"context"
	"errors"
	"time"
)

type Article struct {
	ID         string
	Title      string
	Content    string
	CategoryID string
	UserID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var (
	ErrInvalidTitle    = errors.New("invalid title: cannot be empty")
	ErrTitleTooLong    = errors.New("invalid title: too long")
	ErrArticleNotFound = errors.New("article not found")
)

func (a *Article) Validate() error {
	if a.Title == "" {
		return ErrInvalidTitle
	}
	if len(a.Title) > 255 {
		return ErrTitleTooLong
	}
	return nil
}

type ArticleRepository interface {
	Create(ctx context.Context, article *Article) (*Article, error)
	GetByID(ctx context.Context, id string) (*Article, error)
	Update(ctx context.Context, article *Article) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Article, error)
}
