package models

import "gorm.io/gorm"

type Order struct {
    gorm.Model
    TotalBarang     int     `json:"total_barang" gorm:"not null"`
    TotalHarga      float64 `json:"total_harga" gorm:"not null"`
    MetodePembayaran string `json:"metode_pembayaran" gorm:"not null"`
    Invoice         string  `json:"invoice" gorm:"not null"`
    UserID          uint    `json:"user_id" gorm:"not null"` // Foreign key ke User
    User            User    `gorm:"foreignKey:UserID"`       // Relasi ke User
}