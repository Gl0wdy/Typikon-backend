package main

import (
	"log"

	"github.com/Gl0wdy/Typikon-backend/microservices/wiki-service/internal/app"
)

func main() {
	application, err := app.New("50051")
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
