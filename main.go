package main

import (
	"backend-go-gin/config"
	"backend-go-gin/handlers"
	"backend-go-gin/migrations"
	"github.com/gin-gonic/gin"
)

func init() {
	config.ConnectDB()
}

func main() {
	migrations.Migrate()

	router := gin.Default()

	// Route untuk testing
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	// Public routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	router.Run(":8000")
}