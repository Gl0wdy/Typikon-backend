package repository

import (
	"context"
	"sync"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/google/uuid"
)

type InMemoryArticleRepo struct {
	mu       sync.Mutex
	articles map[string]*domain.Article
}

func NewInMemoryArticleRepo() *InMemoryArticleRepo {
	return &InMemoryArticleRepo{
		articles: make(map[string]*domain.Article),
	}
}

func (r *InMemoryArticleRepo) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	article.ID = uuid.New().String()
	r.articles[article.ID] = article
	return article, nil
}

func (r *InMemoryArticleRepo) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	article, ok := r.articles[id]
	if !ok {
		return nil, domain.ErrArticleNotFound
	}
	return article, nil
}

func (r *InMemoryArticleRepo) Update(ctx context.Context, article *domain.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.articles[article.ID]; !ok {
		return domain.ErrArticleNotFound
	}
	r.articles[article.ID] = article
	return nil
}

func (r *InMemoryArticleRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.articles[id]; !ok {
		return domain.ErrArticleNotFound
	}
	delete(r.articles, id)
	return nil
}

func (r *InMemoryArticleRepo) List(ctx context.Context, limit, offset int) ([]*domain.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]*domain.Article, 0, len(r.articles))
	for _, a := range r.articles {
		result = append(result, a)
	}
	return result, nil
}
