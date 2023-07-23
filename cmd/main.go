package main

import (
	"log"
	"soundproof/config"
	"soundproof/internal/domain/service"
	"soundproof/internal/storage"
	"soundproof/internal/transport"

	"github.com/gin-gonic/gin"
)

func main() {
	configuration, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Unable to load configuration")
	}

	db := storage.ConnectPostgresDB(configuration)
	storage := storage.NewPostgreSQL(db)
	service := service.NewUserService(storage)
	handler := transport.NewHandler(service)

	router := gin.Default()

	router.POST("/auth/register", handler.RegisterUser)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
