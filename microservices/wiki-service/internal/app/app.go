package app

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/repository"
	wiki_grpc "github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/transport/grpc"
	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/usecase"
)

type App struct {
	grpcServer *grpc.Server
	port       string
}

func New(port string) (*App, error) {
	db, err := repository.NewPostgresConnection(repository.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "typikon",
		Password: "typikon",
		DBName:   "wiki",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	articleRepo := repository.NewPostgresArticleRepo(db)
	articleUsecase := usecase.NewArticleUseCase(articleRepo)
	articlesServer := wiki_grpc.NewArticlesServer(articleUsecase)

	grpcServer := grpc.NewServer()
	wiki_grpc.RegisterArticlesServiceServer(grpcServer, articlesServer)

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
