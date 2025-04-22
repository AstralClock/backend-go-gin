package models

import "gorm.io/gorm"

type Payment struct {
    gorm.Model
    OrderID            uint    `json:"order_id"`
    PaymentToken       string  `json:"payment_token" gorm:"size:255"`       // Token dari Snap
    PaymentRedirectURL string  `json:"payment_redirect_url" gorm:"size:255"` // URL redirect Midtrans
    PaymentMethod      string  `json:"payment_method" gorm:"size:50"`        // e.g., "credit_card", "gopay"
    Status             string  `json:"status" gorm:"size:20"`                // "pending", "success", "failed"
    Amount             float64 `json:"amount"`                               // Jumlah pembayaran
    MidtransOrderID    string  `json:"midtrans_order_id" gorm:"size:50"`     // Order ID dari Midtrans (unik)
}