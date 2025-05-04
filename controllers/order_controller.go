package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"backend-go-gin/config"
	"backend-go-gin/models"
	"backend-go-gin/services"
	"backend-go-gin/utils"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	paymentService *services.PaymentService
}

func NewOrderController() *OrderController {
	return &OrderController{
		paymentService: services.NewPaymentService(),
	}
}

func (oc *OrderController) CheckoutSelectedItems(c *gin.Context) {
	var request struct {
		CartDetailIDs []uint `json:"cart_detail_ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get userID from JWT
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get user data
	var user models.User
	if err := config.DB.First(&user, userIDUint).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Get selected cart details with product info and cart ownership validation
	var cartDetails []models.CartDetail
	if err := config.DB.
		Preload("Produk").
		Joins("JOIN carts ON carts.id = cart_details.cart_id").
		Where("cart_details.id IN ?", request.CartDetailIDs).
		Where("carts.user_id = ?", userIDUint).
		Find(&cartDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart items"})
		return
	}

	if len(cartDetails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid cart items found"})
		return
	}

	// Convert selected cart details to order details
	var orderDetails []models.OrderDetail
	totalBarang := 0
	totalHarga := 0.0

	for _, cartDetail := range cartDetails {
		orderDetail := models.OrderDetail{
			ProdukID:    cartDetail.ProdukID,
			TotalProduk: cartDetail.Quantity,
			HargaSatuan: cartDetail.Price,
			HargaTotal:  cartDetail.Subtotal,
		}
		orderDetails = append(orderDetails, orderDetail)

		totalBarang += cartDetail.Quantity
		totalHarga += cartDetail.Subtotal
	}

	// Create invoice number
	invoiceNumber := "INV-" + utils.RandomString(8)

	// Create order
	order := models.Order{
		TotalBarang:      totalBarang,
		TotalHarga:       totalHarga,
		MetodePembayaran: "midtrans",
		Invoice:          invoiceNumber,
		UserID:           userIDUint,
		Status:           "pending",
	}

	// Save order to database
	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Set OrderID for all order details
	for i := range orderDetails {
		orderDetails[i].OrderID = order.ID
	}

	// Create order details
	if err := config.DB.Create(&orderDetails).Error; err != nil {
		config.DB.Delete(&order)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order details"})
		return
	}

	// Get user details for payment
	var userDetail models.UserDetail
	if err := config.DB.Where("user_id = ?", user.ID).First(&userDetail).Error; err != nil {
		config.DB.Delete(&order)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user details"})
		return
	}

	// Create Midtrans transaction
	snapResp, err := oc.paymentService.CreateSnapTransaction(&order, &user, &userDetail)
	if err != nil {
		config.DB.Delete(&order)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment transaction"})
		return
	}

	// Create payment record
	payment := models.Payment{
		OrderID:            order.ID,
		PaymentToken:       snapResp.Token,
		PaymentRedirectURL: snapResp.RedirectURL,
		PaymentMethod:      "midtrans",
		Status:             "pending",
		Amount:             order.TotalHarga,
		MidtransOrderID:    order.Invoice,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		config.DB.Delete(&order)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save payment data"})
		return
	}

	// Periksa dan hapus cart jika kosong
	if len(cartDetails) > 0 {
		var remainingItems int64
		config.DB.Model(&models.CartDetail{}).Where("cart_id = ?", cartDetails[0].CartID).Count(&remainingItems)
		if remainingItems == 0 {
			config.DB.Delete(&models.Cart{}, cartDetails[0].CartID)
		}
	}

	// Hapus cart details yang sudah dipilih
	if err := config.DB.
		Where("id IN ?", request.CartDetailIDs).
		Delete(&models.CartDetail{}).Error; err != nil {
		// Tidak return error, hanya log karena order sudah berhasil dibuat
		fmt.Printf("Failed to delete cart details: %v\n", err)
	}

	// Periksa dan hapus cart jika kosong
	if len(cartDetails) > 0 {
		var remainingItems int64
		if err := config.DB.
			Model(&models.CartDetail{}).
			Where("cart_id = ?", cartDetails[0].CartID).
			Count(&remainingItems).Error; err != nil {
			fmt.Printf("Failed to check remaining cart items: %v\n", err)
		}

		if remainingItems == 0 {
			if err := config.DB.
				Delete(&models.Cart{}, cartDetails[0].CartID).Error; err != nil {
				fmt.Printf("Failed to delete empty cart: %v\n", err)
			}
		}
	}

	var orderWithDetails models.Order
	if err := config.DB.Preload("User.UserDetail").First(&orderWithDetails, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order details"})
		return
	}

	// Ganti response order dengan orderWithDetails yang sudah termasuk User dan UserDetail
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created and cart items removed",
		"order":   orderWithDetails,
		"payment": gin.H{
			"token":        payment.PaymentToken,
			"redirect_url": payment.PaymentRedirectURL,
		},
	})

}

func (oc *OrderController) HandlePaymentNotification(c *gin.Context) {
	// Get order ID from URL parameter
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	// Find the order to get invoice number
	var order models.Order
	if err := config.DB.First(&order, uint(orderID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Verify payment with Midtrans using invoice number
	transactionStatus, err := oc.paymentService.VerifyPayment(order.Invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find related payment record
	var payment models.Payment
	if err := config.DB.Where("order_id = ?", order.ID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment record not found"})
		return
	}

	// Update status based on Midtrans response
	switch transactionStatus.TransactionStatus {
	case "capture", "settlement":
		payment.Status = "success"
		// Kurangi stok produk jika berhasil
		if err := oc.updateProductStock(order.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}
	case "pending":
		payment.Status = "pending"
	default: // deny, expire, cancel, etc
		payment.Status = "failed"
	}

	// Save payment and order status
	if err := config.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	order.Status = payment.Status
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Payment status updated",
		"order_id": order.ID,
		"status":   payment.Status,
		"invoice":  order.Invoice,
	})
}

// Fungsi baru untuk mengurangi stok produk
func (oc *OrderController) updateProductStock(orderID uint) error {
	// Ambil semua order detail untuk order ini
	var orderDetails []models.OrderDetail
	if err := config.DB.Where("order_id = ?", orderID).Find(&orderDetails).Error; err != nil {
		return err
	}

	// Kurangi stok untuk setiap produk
	for _, detail := range orderDetails {
		var product models.Produk
		if err := config.DB.First(&product, detail.ProdukID).Error; err != nil {
			return err
		}

		// Pastikan stok cukup sebelum dikurangi
		if product.Jumlah < detail.TotalProduk {
			return fmt.Errorf("not enough stock for product %d", product.ID)
		}

		product.Jumlah -= detail.TotalProduk
		if err := config.DB.Save(&product).Error; err != nil {
			return err
		}
	}

	return nil
}

func calculateTotals(details []models.OrderDetail) (int, float64) {
	totalBarang := 0
	totalHarga := 0.0

	for _, detail := range details {
		totalBarang += detail.TotalProduk
		totalHarga += detail.HargaTotal
	}

	return totalBarang, totalHarga
}

// GetOrderByID - Dengan alamat dan telepon
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.
		Preload("User.UserDetail"). // Tambahkan ini untuk load UserDetail
		Preload("OrderDetails.Produk").
		First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Format response dengan alamat & telepon
	response := gin.H{
		"message": "Order retrieved successfully",
		"order": gin.H{
			"id":           order.ID,
			"invoice":      order.Invoice,
			"status":       order.Status,
			"total_barang": order.TotalBarang,
			"total_harga":  order.TotalHarga,
			"created_at":   order.CreatedAt,
			"user": gin.H{
				"id":      order.User.ID,
				"name":    order.User.UserDetail.Nama,
				"email":   order.User.Email,
				"phone":   order.User.UserDetail.Telepon, // Tambahkan telepon
				"address": order.User.UserDetail.Alamat,  // Tambahkan alamat
			},
			"order_details": order.OrderDetails,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (oc *OrderController) GetAllOrders(c *gin.Context) {
	var orders []models.Order

	if err := config.DB.
		Preload("User.UserDetail"). // Load UserDetail
		Preload("OrderDetails.Produk").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Format response untuk multiple orders
	var ordersResponse []gin.H
	for _, order := range orders {
		ordersResponse = append(ordersResponse, gin.H{
			"id":      order.ID,
			"invoice": order.Invoice,
			"status":  order.Status,
			"user": gin.H{
				"id":      order.User.ID,
				"name":    order.User.UserDetail.Nama,
				"phone":   order.User.UserDetail.Telepon,
				"address": order.User.UserDetail.Alamat,
			},
			"created_at": order.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Orders retrieved successfully",
		"orders":  ordersResponse,
	})
}
