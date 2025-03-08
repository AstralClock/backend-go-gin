package models

import "gorm.io/gorm"

type Ulasan struct {
    gorm.Model
    Ulasan    string `json:"ulasan" gorm:"not null"`
    ProdukID  uint   `json:"produk_id" gorm:"not null"` // Foreign key ke Produk
    UserID    uint   `json:"user_id" gorm:"not null"`   // Foreign key ke User
    Produk    Produk `gorm:"foreignKey:ProdukID"`       // Relasi ke Produk
    User      User   `gorm:"foreignKey:UserID"`         // Relasi ke User
}