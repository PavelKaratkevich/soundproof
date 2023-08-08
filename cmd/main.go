package main

import (
	"log"
	"soundproof/config"
	"soundproof/docs"
	"soundproof/internal/domain/service"
	"soundproof/internal/storage"
	transport "soundproof/internal/transport/http"
	"soundproof/pkg/logging"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Soundproof API
// @title           Soundproof service
// @version         2.0
// @contact.name   	Pavel Karatkevich
// @contact.url    	https://www.linkedin.com/in/pavel-karatkevich-236461178/
// @contact.email  	p.korotkevitch@gmail.com
// @license.name  	Apache 2.0
// @license.url   	http://www.apache.org/licenses/LICENSE-2.0.html
// @host      		localhost:8080
// @BasePath  		/api/v1.
func main() {
	// load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Unable to load configuration")
	}

	// Create a new logger
	logger := logging.Logger(cfg.Logging.Format, cfg.Logging.Level)

	// Sync the logger before exiting.
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Fatalf("syncing logger: %v", err)
		}
	}(logger)

	// Log the config
	logger.Debug("config", zap.Any("config", cfg))

	// Wire all layers together
	db := storage.ConnectPostgresDB(logger, cfg)
	storage := storage.NewPostgreSQL(logger, db)
	service := service.NewUserService(logger, storage)
	handler := transport.NewHandler(logger, service)

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Soundproof API"
	docs.SwaggerInfo.Description = "This is a soundproof API performed as a test assignment."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// launch a Gin router
	router := gin.Default()

	// define routes
	router.POST("/auth/register", handler.RegisterUser)
	router.POST("/auth/login", handler.Login)
	router.GET("/user/profile", handler.GetUser)
	router.PUT("/user/profile", handler.UpdateUser)

	// swagger
	if gin.Mode() == gin.DebugMode {
		router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	
	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
