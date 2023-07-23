package main

import (
	"log"
	"soundproof/config"
	"soundproof/internal/domain/service"
	"soundproof/internal/storage"
	"soundproof/internal/transport"
	"soundproof/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Unable to load configuration")
	}

	// Create a new logger.
	logger := logging.Logger(cfg.Logging.Format, cfg.Logging.Level)

	// Sync the logger before exiting.
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Fatalf("syncing logger: %v", err)
		}
	}(logger)

	// Log the config.
	logger.Debug("config", zap.Any("config", cfg))

	db := storage.ConnectPostgresDB(logger)
	storage := storage.NewPostgreSQL(logger, db)
	service := service.NewUserService(logger, storage)
	handler := transport.NewHandler(logger, service)

	router := gin.Default()

	router.POST("/auth/register", handler.RegisterUser)

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
