package main

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/app"
	wiki_grpc "github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/transport/grpc"
)

func main() {
	application, err := app.New("50051")
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("failed to run grpc server: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := wiki_grpc.RegisterArticlesServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register articles gateway: %v", err)
	}
	if err := wiki_grpc.RegisterCategoriesServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register categories gateway: %v", err)
	}

	log.Println("HTTP gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve gateway: %v", err)
	}
}
