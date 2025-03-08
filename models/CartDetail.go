package models

import "gorm.io/gorm"

type CartDetail struct {
    gorm.Model
    CartID    uint    `json:"cart_id" gorm:"not null"`    // Foreign key ke Cart
    ProdukID  uint    `json:"produk_id" gorm:"not null"`  // Foreign key ke Produk
    Price     float64 `json:"price" gorm:"not null"`
    Quantity  int     `json:"quantity" gorm:"not null"`
    Subtotal  float64 `json:"subtotal" gorm:"not null"`
    Cart      Cart    `gorm:"foreignKey:CartID"`          // Relasi ke Cart
    Produk    Produk  `gorm:"foreignKey:ProdukID"`       // Relasi ke Produk
}