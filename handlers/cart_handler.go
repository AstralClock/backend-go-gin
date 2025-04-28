package handlers

import (
	"backend-go-gin/config"
	"backend-go-gin/controllers"
	"backend-go-gin/models"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type AddToCartRequest struct {
	ProdukID uint `json:"produk_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}


func AddToCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartCtrl := controllers.NewCartController(config.DB)
	cart, err := cartCtrl.GetOrCreateActiveCart(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat/mengambil keranjang"})
		return
	}

	var produk models.Produk
	if err := config.DB.First(&produk, req.ProdukID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	cartDetail := models.CartDetail{
		CartID:    cart.ID,
		ProdukID:  req.ProdukID,
		Quantity:  req.Quantity,
		Price:     produk.Harga,
		Subtotal:  produk.Harga * float64(req.Quantity),
	}

	if err := config.DB.Create(&cartDetail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan item ke keranjang"})
		return
	}

	config.DB.Model(&cart).Update("total_barang", gorm.Expr("total_barang + ?", req.Quantity))

	c.JSON(http.StatusOK, gin.H{
		"message": "Item berhasil ditambahkan ke keranjang",
		"cart_id": cart.ID,
	})
}


func UpdateCartItem(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartCtrl := controllers.NewCartController(config.DB)
	updatedItem, err := cartCtrl.UpdateCartItem(uint(itemID), userID, req.Quantity)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedItem)
}

func DeleteCartItem(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
		return
	}

	cartCtrl := controllers.NewCartController(config.DB)
	if err := cartCtrl.DeleteCartItem(uint(itemID), userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}