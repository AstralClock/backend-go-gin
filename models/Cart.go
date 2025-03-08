package models

import "gorm.io/gorm"

type Cart struct {
    gorm.Model
    TotalBarang int    `json:"total_barang" gorm:"not null"`
    Status      string `json:"status" gorm:"not null"`
    UserID      uint   `json:"user_id" gorm:"not null"` // Foreign key ke User
    User        User   `gorm:"foreignKey:UserID"`       // Relasi ke User
}