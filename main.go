// filepath: /C:/Users/Fauzy/Documents/Tugas/SEM 4/SI/Project/dashboard/backend-go-gin/main.go
package main

import (
    "backend-go-gin/config"
    "backend-go-gin/handlers"
    "backend-go-gin/middleware"
    "github.com/gin-gonic/gin"
    "log"
    "backend-go-gin/migrations"
)

func init() {
    config.ConnectDB()
}

func main() {
    // Connect to the database
    migrations.Migrate()

    // Initialize Gin router
    r := gin.Default()

    // Use CORS middleware
    r.Use(middleware.CORSMiddleware())

    // Define routes
    r.POST("/login", handlers.Login)

    // Route untuk testing
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello, World!"})
    })

    r.POST("/addproducts", handlers.AddProductHandler)
    r.DELETE("/delproducts/:id", handlers.DeleteProductHandler)
    r.GET("/products", handlers.GetAllProductsHandler)

    // Add your protected routes here
    // Protected routes
    protected := r.Group("/")
    protected.Use(middleware.Auth())
    {
        // protected.POST("/products", handlers.AddProductHandler)
    }

    // Start the server
    log.Println("Server started at :8000")
    log.Fatal(r.Run(":8000"))
}