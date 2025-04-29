// filepath: /C:/Users/Fauzy/Documents/Tugas/SEM 4/SI/Project/dashboard/backend-go-gin/main.go
package main

import (
	"backend-go-gin/config"
	"backend-go-gin/controllers"
	"backend-go-gin/handlers"
	"backend-go-gin/middleware"
	"backend-go-gin/migrations"

	"github.com/gin-gonic/gin"
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

	r.MaxMultipartMemory = 3 << 20 // 3MB
	r.Use(middleware.CORSMiddleware())
	orderController := controllers.NewOrderController()
    // Route untuk testing
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Hello, World!"})
    })

	// Public routes
	r.POST("/login", handlers.Login)
	r.POST("/payment/notification", orderController.HandlePaymentNotification)

	//admin routes
	r.POST("/addproducts", handlers.AddProductHandler)
	r.DELETE("/delproducts/:id", handlers.DeleteProductHandler)
	r.GET("/products", handlers.GetAllProductsHandler)
	r.PUT("/editproducts/:id", handlers.EditProductHandler)
	r.POST("/upload/products/:id", handlers.UploadProductImages)

	// Private routes
	r.POST("/addcart", middleware.AuthMiddleware(), handlers.AddToCart)
	r.POST("/orders", middleware.AuthMiddleware(), orderController.CheckoutSelectedItems)

	r.Static("/uploads", "./uploads")
	r.Run(":8000")
}
