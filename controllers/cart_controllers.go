package controllers

import (
	"backend-go-gin/config"
	"backend-go-gin/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartController struct {
	DB *gorm.DB
}

func NewCartController(db *gorm.DB) *CartController {
	return &CartController{DB: db}
}

func (cc *CartController) GetOrCreateActiveCart(userID uint) (models.Cart, error) {
	var cart models.Cart
	err := cc.DB.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error

	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{
			UserID:      userID,
			Status:      "active",
			TotalBarang: 0,
		}
		err = cc.DB.Create(&cart).Error
	}

	return cart, err
}

func (cc *CartController) UpdateCartItem(cartItemID uint, userID uint, quantity int) (models.CartDetail, error) {
	var cartDetail models.CartDetail

	err := cc.DB.Transaction(func(tx *gorm.DB) error {
		// Verify item belongs to user's active cart
		if err := tx.Joins("JOIN carts ON carts.id = cart_details.cart_id").
			Where("cart_details.id = ? AND carts.user_id = ? AND carts.status = ?", cartItemID, userID, "active").
			First(&cartDetail).Error; err != nil {
			return errors.New("item tidak ditemukan di keranjang anda")
		}

		// Update quantity and subtotal
		cartDetail.Quantity = quantity
		cartDetail.Subtotal = cartDetail.Price * float64(quantity)

		return tx.Save(&cartDetail).Error
	})

	return cartDetail, err
}

// New: Delete Cart Item
func (cc *CartController) DeleteCartItem(cartItemID uint, userID uint) error {
	return cc.DB.Transaction(func(tx *gorm.DB) error {
		// Verify item belongs to user's active cart
		var cartDetail models.CartDetail
		if err := tx.Joins("JOIN carts ON carts.id = cart_details.cart_id").
			Where("cart_details.id = ? AND carts.user_id = ? AND carts.status = ?", cartItemID, userID, "active").
			First(&cartDetail).Error; err != nil {
			return errors.New("item tidak ditemukan di keranjang anda")
		}

		// Delete item
		if err := tx.Delete(&cartDetail).Error; err != nil {
			return err
		}

		// Update cart total items
		return tx.Model(&models.Cart{}).
			Where("id = ?", cartDetail.CartID).
			Update("total_barang", gorm.Expr("total_barang - ?", cartDetail.Quantity)).Error
	})
}

type CartsController struct{}

func NewCartsController() *CartController {
	return &CartController{}
}

// GetUserCart - Get cart data for logged in user
func (cc *CartController) GetUserCart(c *gin.Context) {
	// Ambil userID dari JWT
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var cart models.Cart
	if err := config.DB.
		Preload("User").
		Preload("CartDetails", "deleted_at IS NULL"). // Hanya yang belum dihapus
		Where("user_id = ?", userID).
		First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Hitung total item yang belum dihapus
	var activeCartItemsCount int64
	config.DB.Model(&models.CartDetail{}).
		Where("cart_id = ? AND deleted_at IS NULL", cart.ID).
		Count(&activeCartItemsCount)

	// Hitung total harga dari item yang belum dihapus
	var totalPrice float64
	config.DB.Model(&models.CartDetail{}).
		Select("COALESCE(SUM(harga * jumlah), 0)").
		Where("cart_id = ? AND deleted_at IS NULL", cart.ID).
		Scan(&totalPrice)

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart retrieved successfully",
		"data": gin.H{
			"cart":              cart,
			"total_items":       activeCartItemsCount,
			"total_price":       totalPrice,
			"active_cart_items": activeCartItemsCount,
		},
	})
}

// GetUserCartDetails - Get cart details for logged in user
func (cc *CartController) GetUserCartDetails(c *gin.Context) {
    // Ambil userID dari JWT
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var cartDetails []models.CartDetail
    if err := config.DB.
        Preload("Produk").
        Preload("Ukuran").
        Joins("JOIN carts ON carts.id = cart_details.cart_id").
        Where("carts.user_id = ? AND cart_details.deleted_at IS NULL", userID).
        Find(&cartDetails).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart items"})
        return
    }

    // Format response langsung menggunakan model yang sudah di-preload
    type cartItemResponse struct {
        ID         uint     `json:"id"`
        ProdukID   uint     `json:"produk_id"`
        NamaProduk string   `json:"nama_produk"`
        Harga      float64  `json:"harga"`
        Quantity   int      `json:"quantity"`
        Subtotal   float64  `json:"subtotal"`
        Ukuran     string   `json:"ukuran,omitempty"` // Nama ukuran jika ada
        Image      string   `json:"image"`
    }

    var response []cartItemResponse
    for _, item := range cartDetails {
        ukuran := ""
        if item.Ukuran.ID != 0 { // Cek apakah ada ukuran
            ukuran = item.Ukuran.Ukuran
        }

        response = append(response, cartItemResponse{
            ID:         item.ID,
            ProdukID:   item.ProdukID,
            NamaProduk: item.Produk.NamaProduk,
            Harga:      item.Price,
            Quantity:   item.Quantity,
            Subtotal:   item.Subtotal,
            Ukuran:     ukuran,
            Image:      item.Produk.Image,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "message":    "Cart items retrieved successfully",
        "cart_items": response,
    })
}
