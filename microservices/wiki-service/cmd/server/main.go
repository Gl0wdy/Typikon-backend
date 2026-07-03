package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	wiki_grpc "github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/transport/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Can't listen on :50051: %v", err)
	}

	baseServer := grpc.NewServer()

	articlesImpl := &wiki_grpc.ArticlesServer{}

	wiki_grpc.RegisterArticlesServiceServer(baseServer, articlesImpl)

	log.Println("gRPC server started on :50051...")

	if err := baseServer.Serve(listener); err != nil {
		log.Fatalf("Error while serving: %v", err)
	}
}
