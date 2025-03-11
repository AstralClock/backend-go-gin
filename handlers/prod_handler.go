package handlers

import (
    "backend-go-gin/controllers"
    "github.com/gin-gonic/gin"
)

// AddProductHandler handles the request to add a new product
func AddProductHandler(c *gin.Context) {
    controllers.AddProduct(c)
}

// DeleteProductHandler handles the request to delete a product by its ID
func DeleteProductHandler(c *gin.Context) {
    controllers.DeleteProduct(c)
}

// GetAllProductsHandler handles the request to get all products
func GetAllProductsHandler(c *gin.Context) {
    controllers.GetAllProducts(c)
}