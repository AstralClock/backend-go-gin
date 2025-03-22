package main

import (
	"backend-go-gin/config"
	"backend-go-gin/handlers"
	"backend-go-gin/migrations"
	"backend-go-gin/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	config.ConnectDB()
}

func main() {
	migrations.Migrate()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	router.MaxMultipartMemory = 3 << 20 // 3MB

	userDetailHandler := handlers.NewUserDetailHandler()

	// Public routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)
	router.GET("/verify-email", handlers.VerifyEmail)

	// Private routes
	router.POST("/user-detail", middleware.AuthMiddleware(), userDetailHandler.SaveUserDetail)

	router.Run(":8000")
}