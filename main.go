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

	r.MaxMultipartMemory = 3 << 20 // 3MB
	r.Use(middleware.CORSMiddleware())
	userDetailHandler := handlers.NewUserDetailHandler()
	orderController := controllers.NewOrderController()
	userDetail := controllers.NewUserController()
	cartsController := controllers.NewCartsController()

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/loginadmin", handlers.LoginAdmin)
	r.POST("/login", handlers.Login)
	r.GET("/verify-email", handlers.VerifyEmail)
	r.POST("/payment/notification/:id", orderController.HandlePaymentNotification)
	r.GET("getorders", orderController.GetAllOrders)
	r.GET("gerorder/:id", orderController.GetOrderByID)
	r.GET("/users", userDetail.GetAllUserDetails)
	r.GET("/user/:id", userDetail.GetUserDetailByID)

	//admin routes
	r.POST("/addproducts", handlers.AddProductHandler)
	r.DELETE("/delproducts/:id", handlers.DeleteProductHandler)
	r.GET("/products", handlers.GetAllProductsHandler)
	r.PUT("/editproducts/:id", handlers.EditProductHandler)
	r.POST("/upload/products/:id", handlers.UploadProductImages)

	// Private routes
	r.POST("/postuser", middleware.AuthMiddleware(), userDetailHandler.SaveUserDetail)
	r.PUT("/edituser", middleware.AuthMiddleware(), userDetailHandler.UpdateUserDetail)
	r.POST("/addcart", middleware.AuthMiddleware(), handlers.AddToCart)
	r.PUT("/updatecartdetail/:id", middleware.AuthMiddleware(), handlers.UpdateCartItem)
	r.DELETE("/deletecart/:id", middleware.AuthMiddleware(), handlers.DeleteCartItem)
	r.POST("/orders", middleware.AuthMiddleware(), orderController.CheckoutSelectedItems)
	r.GET("/carts", middleware.AuthMiddleware(), cartsController.GetUserCart)
	r.GET("/cartdetail", middleware.AuthMiddleware(), cartsController.GetUserCartDetails)
	r.GET("/getuserlogin", middleware.AuthMiddleware(), userDetail.GetCurrentUserDetail)
	r.GET("/user/orders", middleware.AuthMiddleware(), orderController.GetUserOrderDetails)
	r.GET("/user/orders/:id", middleware.AuthMiddleware(), orderController.GetUserOrderByID)

	r.Run(":8000")
}
