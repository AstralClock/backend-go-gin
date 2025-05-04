package handlers

import (
	"backend-go-gin/controllers"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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

func EditProductHandler(c *gin.Context) {
    controllers.EditProduct(c)
}

func UploadProductImages(c *gin.Context) {
    productID := c.Param("id") // Get the product ID from the URL parameter

    if productID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
        return
    }

    // Create the directory for the product's images
    uploadDir := filepath.Join("uploads", "products", productID)
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
        return
    }

    // Get the uploaded files
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
        return
    }

    files := form.File["images"] // "images" is the key for the uploaded files
    for i, file := range files {
        // Save each file with a sequential name (e.g., 1.jpg, 2.jpg, etc.)
        filePath := filepath.Join(uploadDir, fmt.Sprintf("%d.jpg", i+1))
        if err := c.SaveUploadedFile(file, filePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully"})
}