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
    if err := config.DB.Save(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Initialize sizes with quantity 0 for the new product
    for sizeID := 1; sizeID <= 5; sizeID++ {
        size := models.ProdukUkuran{
            ProdukID: product.ID,
            UkuranID: uint(sizeID),
            Stok:     0,
        }
        if err := config.DB.Create(&size).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize sizes"})
            return
        }
    }

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

// GetAllProducts returns all products from the database with calculated jumlah
func GetAllProducts(c *gin.Context) {
	var products []models.Produk
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, product := range products {
		var productSizes []models.ProdukUkuran
		// Fetch sizes and calculate total quantity
		if err := config.DB.Where("produk_id = ?", product.ID).Find(&productSizes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalQuantity := 0
		for _, size := range productSizes {
			totalQuantity += size.Stok
		}

		result = append(result, gin.H{
			"ID":          product.ID,
			"nama_produk": product.NamaProduk,
			"deskripsi":   product.Deskripsi,
			"kategori":    product.Kategori,
			"harga":       product.Harga,
			"jumlah":      totalQuantity,
			"image":       product.Image,
		})
	}

	c.JSON(http.StatusOK, gin.H{"products": result})
}

func EditProductSizes(c *gin.Context) {
    id := c.Param("id") // Get the product ID from the request parameters

    var sizeUpdates []struct {
        UkuranID uint `json:"ukuran_id"` // Size ID
        Stok     int  `json:"stok"`      // New stock quantity
    }

    // Bind JSON to sizeUpdates
    if err := c.ShouldBindJSON(&sizeUpdates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Iterate through the updates and apply them
    for _, update := range sizeUpdates {
        var productSize models.ProdukUkuran
        // Find the specific size for the product
        if err := config.DB.Where("produk_id = ? AND ukuran_id = ?", id, update.UkuranID).First(&productSize).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Size with UkuranID %d not found for product %s", update.UkuranID, id)})
            return
        }

        // Update the stock
        productSize.Stok = update.Stok
        if err := config.DB.Save(&productSize).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update size with UkuranID %d", update.UkuranID)})
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"message": "Sizes updated successfully"})
}

// GetProductByID retrieves a product by its ID with calculated jumlah
func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Produk
	// Find the product by ID
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var productSizes []models.ProdukUkuran
	// Fetch sizes and calculate total quantity
	if err := config.DB.Where("produk_id = ?", product.ID).Find(&productSizes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalQuantity := 0
	for _, size := range productSizes {
		totalQuantity += size.Stok
	}

	c.JSON(http.StatusOK, gin.H{
		"product": gin.H{
			"ID":          product.ID,
			"nama_produk": product.NamaProduk,
			"deskripsi":   product.Deskripsi,
			"kategori":    product.Kategori,
			"harga":       product.Harga,
			"jumlah":      totalQuantity,
			"image":       product.Image,
		},
	})
}

// GetProductWithSizes retrieves a product by its ID along with sizes and calculated jumlah
func GetProductWithSizes(c *gin.Context) {
	id := c.Param("id") // Get the product ID from the request parameters

	var product models.Produk
	// Find the product by ID
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var productSizes []models.ProdukUkuran
	// Fetch product sizes and their quantities
	if err := config.DB.Where("produk_id = ?", id).Preload("Ukuran").Find(&productSizes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate the total quantity (jumlah) and simplify the sizes output
	totalQuantity := 0
	sizes := []gin.H{}
	for _, size := range productSizes {
		totalQuantity += size.Stok
		sizes = append(sizes, gin.H{
			"Ukuran": size.Ukuran.Ukuran, // Assuming Ukuran struct has a field `Ukuran` for size name
			"Stok":   size.Stok,
		})
	}

	// Simplify the product output and include the calculated jumlah
	c.JSON(http.StatusOK, gin.H{
		"product": gin.H{
			"ID":          product.ID,
			"nama_produk": product.NamaProduk,
			"deskripsi":   product.Deskripsi,
			"kategori":    product.Kategori,
			"harga":       product.Harga,
			"jumlah":      totalQuantity, // Use the calculated total quantity
			"image":       product.Image,
		},
		"sizes": sizes,
	})
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
