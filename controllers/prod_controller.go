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
	
	// Preload relasi yang diperlukan
	if err := config.DB.
		Preload("Ukurans").
		Preload("ProdukUkuranStock.Ukuran").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format response
	type ukuranResponse struct {
		ID   uint   `json:"id"`
		Nama string `json:"nama"`
		Stok int    `json:"stok"`
	}

	type produkResponse struct {
		ID         uint             `json:"id"`
		CreatedAt  string           `json:"created_at"`
		UpdatedAt  string           `json:"updated_at"`
		NamaProduk string           `json:"nama_produk"`
		Deskripsi  string           `json:"deskripsi"`
		Kategori   string           `json:"kategori"`
		Tag        string           `json:"tag"`
		Harga      float64          `json:"harga"`
		Jumlah     int              `json:"jumlah"`
		Image      string           `json:"image"`
		Ukurans    []ukuranResponse `json:"ukurans"`
	}

	var response []produkResponse
	
	for _, p := range products {
		// Map ukuran dengan stok
		var ukurans []ukuranResponse
		for _, pu := range p.ProdukUkuranStock {
			ukurans = append(ukurans, ukuranResponse{
				ID:   pu.Ukuran.ID,
				Nama: pu.Ukuran.Ukuran,
				Stok: pu.Stok,
			})
		}

		response = append(response, produkResponse{
			ID:         p.ID,
			NamaProduk: p.NamaProduk,
			Deskripsi:  p.Deskripsi,
			Kategori:   p.Kategori,
			Tag:        p.Tag,
			Harga:      p.Harga,
			Jumlah:     p.Jumlah,
			Image:      p.Image,
			Ukurans:    ukurans,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"products": response,
	})
}

// GetProductByID retrieves a product by its ID
func GetProductByID(c *gin.Context) {
    id := c.Param("id")

    var product models.Produk
    // Find the product by ID
    if err := config.DB.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"product": product})
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