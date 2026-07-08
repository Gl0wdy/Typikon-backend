package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/app"
	wiki_grpc "github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/transport/grpc"
)

func main() {
	grpcPort := "50051"
	application, err := app.New(grpcPort)
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

	grpcEndpoint := "127.0.0.1:" + grpcPort

	if err := wiki_grpc.RegisterArticlesServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		log.Fatalf("failed to register articles gateway: %v", err)
	}

	if err := wiki_grpc.RegisterCategoriesServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		log.Fatalf("failed to register categories gateway: %v", err)
	}

	if err := wiki_grpc.RegisterCommentsServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		log.Fatalf("failed to register comments gateway: %v", err)
	}

	log.Println("HTTP gateway listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve gateway: %v", err)
	}
}
