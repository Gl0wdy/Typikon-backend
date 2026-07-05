package wiki_grpc

import (
	"context"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/usecase"
)

type ArticlesServer struct {
	UnimplementedArticlesServiceServer
	usecase usecase.ArticleUseCase
}

func NewArticlesServer(uc usecase.ArticleUseCase) *ArticlesServer {
	return &ArticlesServer{usecase: uc}
}

func (s *ArticlesServer) CreateArticle(ctx context.Context, req *CreateArticleRequest) (*ArticleResponse, error) {
	article := &domain.Article{
		Title:      req.GetTitle(),
		Content:    req.GetContent(),
		CategoryID: req.GetCategoryId(),
		UserID:     req.GetUserId(),
	}

	created, err := s.usecase.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	return &ArticleResponse{
		Id:         created.ID,
		Title:      created.Title,
		Content:    created.Content,
		CategoryId: created.CategoryID,
		UserId:     created.UserID,
		CreatedAt:  created.CreatedAt.String(),
	}, nil
}

func (s *ArticlesServer) GetArticle(ctx context.Context, req *GetArticleRequest) (*ArticleResponse, error) {
	article, err := s.usecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &ArticleResponse{
		Id:         article.ID,
		Title:      article.Title,
		Content:    article.Content,
		CategoryId: article.CategoryID,
		UserId:     article.UserID,
		CreatedAt:  article.CreatedAt.String(),
	}, nil
}

func (s *ArticlesServer) UpdateArticle(ctx context.Context, req *UpdateArticleRequest) (*ArticleResponse, error) {
	article := &domain.Article{
		ID:         req.GetId(),
		Title:      req.GetTitle(),
		Content:    req.GetContent(),
		CategoryID: req.GetCategoryId(),
	}

	if err := s.usecase.Update(ctx, article); err != nil {
		return nil, err
	}

	return &ArticleResponse{
		Id:         article.ID,
		Title:      article.Title,
		Content:    article.Content,
		CategoryId: article.CategoryID,
		UserId:     article.UserID,
	}, nil
}

func (s *ArticlesServer) DeleteArticle(ctx context.Context, req *DeleteArticleRequest) (*DeleteArticleResponse, error) {
	if err := s.usecase.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &DeleteArticleResponse{
		Success: true,
	}, nil
}

func (s *ArticlesServer) ListArticles(ctx context.Context, req *ListArticlesRequest) (*ListArticlesResponse, error) {
	articles, err := s.usecase.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, err
	}

	response := &ListArticlesResponse{}
	for _, article := range articles {
		response.Articles = append(response.Articles, &ArticleResponse{
			Id:         article.ID,
			Title:      article.Title,
			Content:    article.Content,
			CategoryId: article.CategoryID,
			UserId:     article.UserID,
			CreatedAt:  article.CreatedAt.String(),
		})
	}

	return response, nil
}
