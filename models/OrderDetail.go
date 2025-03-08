package models

import "gorm.io/gorm"

type OrderDetail struct {
    gorm.Model
    OrderID     uint    `json:"order_id" gorm:"not null"`       // Foreign key ke Order
    ProdukID    uint    `json:"produk_id" gorm:"not null"`      // Foreign key ke Produk
    TotalProduk int     `json:"total_produk" gorm:"not null"`
    HargaSatuan float64 `json:"harga_satuan" gorm:"not null"`
    HargaTotal  float64 `json:"harga_total" gorm:"not null"`
    Order       Order   `gorm:"foreignKey:OrderID"`            // Relasi ke Order
    Produk      Produk  `gorm:"foreignKey:ProdukID"`           // Relasi ke Produk
}