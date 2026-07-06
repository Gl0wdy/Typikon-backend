package wiki_grpc

import (
	"context"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/domain"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/usecase"
)

type CategoriesServer struct {
	UnimplementedCategoriesServiceServer
	usecase usecase.CategoryUseCase
}

func NewCategoriesServer(uc usecase.CategoryUseCase) *CategoriesServer {
	return &CategoriesServer{usecase: uc}
}

func (s *CategoriesServer) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error) {
	category := &domain.Category{
		Title: req.GetTitle(),
	}

	created, err := s.usecase.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:    created.ID,
		Title: created.Title,
	}, nil
}

func (s *CategoriesServer) GetCategory(ctx context.Context, req *GetCategoryRequest) (*CategoryResponse, error) {
	category, err := s.usecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:    category.ID,
		Title: category.Title,
	}, nil
}

func (s *CategoriesServer) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	category := &domain.Category{
		ID:    req.GetId(),
		Title: req.GetNewTitle(),
	}

	updated, err := s.usecase.Update(ctx, category)
	if err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Id:    updated.ID,
		Title: updated.Title,
	}, nil
}

func (s *CategoriesServer) DeleteCategory(ctx context.Context, req *DeleteCategoryRequest) (*DeleteCategoryResponse, error) {
	err := s.usecase.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &DeleteCategoryResponse{
		Success: true,
	}, nil
}

func (s *CategoriesServer) ListCategories(ctx context.Context, req *ListCategoriesRequest) (*ListCategoriesResponse, error) {
	categories, err := s.usecase.List(ctx, int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return nil, err
	}

	resp := &ListCategoriesResponse{}
	for _, category := range categories {
		resp.Categories = append(resp.Categories, &CategoryResponse{
			Id:    category.ID,
			Title: category.Title,
		})
	}

	return resp, nil
}
