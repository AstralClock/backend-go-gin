package controllers

import (
	"backend-go-gin/models"
	"errors"
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