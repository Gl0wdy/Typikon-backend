package wiki_grpc

import (
	"context"
)

type ArticlesServer struct {
	UnimplementedArticlesServiceServer
}

func (s *ArticlesServer) CreateArticle(ctx context.Context, req *CreateArticleRequest) (*ArticleResponse, error) {
	return &ArticleResponse{
		Id:    "test-id-123",
		Title: req.GetTitle(),
	}, nil
}

func (s *ArticlesServer) GetArticle(ctx context.Context, req *GetArticleRequest) (*ArticleResponse, error) {
	return &ArticleResponse{
		Id:    req.GetId(),
		Title: "Sample Article Title",
	}, nil
}

func (s *ArticlesServer) UpdateArticle(ctx context.Context, req *UpdateArticleRequest) (*ArticleResponse, error) {
	return &ArticleResponse{
		Id:    req.GetId(),
		Title: req.GetTitle(),
	}, nil
}

func (s *ArticlesServer) DeleteArticle(ctx context.Context, req *DeleteArticleRequest) (*DeleteArticleResponse, error) {
	return &DeleteArticleResponse{
		Success: true,
	}, nil
}

func (s *ArticlesServer) GetArticles(ctx context.Context, req *GetArticlesRequest) (*GetArticlesResponse, error) {
	articles := []*ArticleResponse{
		{Id: "1", Title: "Article 1"},
		{Id: "2", Title: "Article 2"},
	}
	return &GetArticlesResponse{
		Articles: articles,
	}, nil
}
