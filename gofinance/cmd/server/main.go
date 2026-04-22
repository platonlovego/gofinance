package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yourusername/gofinance/config"
	"github.com/yourusername/gofinance/internal/handler"
	"github.com/yourusername/gofinance/internal/repository"
	"github.com/yourusername/gofinance/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}

	db, err := config.NewPostgres()
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer db.Close()

	repo := repository.New(db)
	svc := service.New(repo)
	h := handler.New(svc)

	router := h.InitRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
