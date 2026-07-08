package app

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository"
	wiki_grpc "github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/transport/grpc"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/usecase"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/config"
)

type App struct {
	grpcServer *grpc.Server
	port       string
}

func New(port string) (*App, error) {
	cfg := config.LoadConfig()

	db, err := repository.NewPostgresConnection(repository.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	categoryRepo := repository.NewPostgresCategoryRepo(db)
	categoryUsecase := usecase.NewCategoryUseCase(categoryRepo)
	categoriesServer := wiki_grpc.NewCategoriesServer(categoryUsecase)

	articleRepo := repository.NewPostgresArticleRepo(db)
	articleUsecase := usecase.NewArticleUseCase(articleRepo, categoryRepo)
	articlesServer := wiki_grpc.NewArticlesServer(articleUsecase)

	commentRepo := repository.NewPostgresCommentRepo(db)
	commentUsecase := usecase.NewCommentUseCase(commentRepo, articleRepo)
	commentsServer := wiki_grpc.NewCommentsServer(commentUsecase)

	grpcServer := grpc.NewServer()
	wiki_grpc.RegisterArticlesServiceServer(grpcServer, articlesServer)
	wiki_grpc.RegisterCategoriesServiceServer(grpcServer, categoriesServer)
	wiki_grpc.RegisterCommentsServiceServer(grpcServer, commentsServer)

	return &App{
		grpcServer: grpcServer,
		port:       port,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("wiki-service is running on port %s", a.port)
	return a.grpcServer.Serve(lis)
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
}
