package controllers

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context) {
    var product models.Produk

    // Bind JSON to product model
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Save product to database
    if err := config.DB.Create(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Update the image field with the product ID
    product.Image = fmt.Sprintf("%d", product.ID)
    config.DB.Save(&product)

    c.JSON(http.StatusOK, gin.H{
        "message": "Product added successfully",
        "product": product,
    })
}

// DeleteProduct deletes a product by its ID
func DeleteProduct(c *gin.Context) {
    id := c.Param("id")

    var product models.Produk
    if err := config.DB.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    if err := config.DB.Delete(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetAllProducts returns all products from the database
func GetAllProducts(c *gin.Context) {
    var products []models.Produk
    if err := config.DB.Find(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"products": products})
}

// EditProduct updates a product by its ID
func EditProduct(c *gin.Context) {
    id := c.Param("id")
    var product models.Produk

    // Find the product by ID
    if err := config.DB.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    // Bind JSON to product model
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Save updated product to database
    if err := config.DB.Save(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
}