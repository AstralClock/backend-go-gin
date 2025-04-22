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
	migrations.Migrate()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	r.MaxMultipartMemory = 3 << 20 // 3MB
	r.Use(middleware.CORSMiddleware())
	userDetailHandler := handlers.NewUserDetailHandler()
	orderController := controllers.NewOrderController()

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/verify-email", handlers.VerifyEmail)
	r.POST("/payment/notification", orderController.HandlePaymentNotification)

	//admin routes
	r.POST("/addproducts", handlers.AddProductHandler)
	r.DELETE("/delproducts/:id", handlers.DeleteProductHandler)
	r.GET("/products", handlers.GetAllProductsHandler)
	r.PUT("/editproducts/:id", handlers.EditProductHandler)

	// Private routes
	r.POST("/user-detail", middleware.AuthMiddleware(), userDetailHandler.SaveUserDetail)
	r.PUT("/user-detail", middleware.AuthMiddleware(), userDetailHandler.UpdateUserDetail)
	r.POST("/addcart", middleware.AuthMiddleware(), handlers.AddToCart)
	r.POST("/orders", middleware.AuthMiddleware(), orderController.CheckoutSelectedItems)

	r.Run(":8000")
}
