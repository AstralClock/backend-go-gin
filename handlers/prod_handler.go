package handlers

import (
	"backend-go-gin/config"
	"backend-go-gin/controllers"
	"backend-go-gin/models"
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

func GetProductByIDHandler(c *gin.Context) {
    id := c.Param("id")

    var produk models.Produk
    if err := config.DB.Preload("Ukurans").
        Preload("ProdukUkuranStock.Ukuran").
        First(&produk, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    // Membuat struktur untuk ukuran dengan stok
    var ukurans []gin.H
    for _, pu := range produk.ProdukUkuranStock {
        ukurans = append(ukurans, gin.H{
            "id":   pu.Ukuran.ID,
            "nama": pu.Ukuran.Ukuran, // Pastikan ini sesuai field di model Ukuran
            "stok": pu.Stok,
        })
    }

    // Membuat response dengan wrapper data
    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "id":          produk.ID,
            "nama_produk": produk.NamaProduk,
            "deskripsi":   produk.Deskripsi,
            "kategori":    produk.Kategori,
            "tag":         produk.Tag,
            "harga":       produk.Harga,
            "jumlah":      produk.Jumlah,
            "image":       produk.Image,
            "created_at":  produk.CreatedAt,
            "updated_at":  produk.UpdatedAt,
            "ukurans":     ukurans,
        },
        "meta": gin.H{
            "status": "success",
        },
    })
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
