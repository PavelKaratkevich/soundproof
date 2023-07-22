package main

import (
	"log"
	"soundproof/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {

	// configuration, err := config.NewConfig()
	// if err != nil {
	// 	log.Fatalf("Unable to load configuration")
	// }

	db := storage.ConnectPostgresDB()
	storage.NewPostgreSQL(db)


	router := gin.Default()

	router.GET("/auth/register", func(c *gin.Context) {
		c.String(200, "Success")
	})

	if err := router.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
