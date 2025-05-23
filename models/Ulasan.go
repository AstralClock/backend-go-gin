package models

import "gorm.io/gorm"

type Ulasan struct {
    gorm.Model
    Ulasan    string  `json:"ulasan" gorm:"not null"`
    Rating    int     `json:"rating" gorm:"not null"`
    ProdukID  uint    `json:"produk_id" gorm:"not null"`
    UserID    uint    `json:"user_id" gorm:"not null"`
    Produk    Produk  `gorm:"foreignKey:ProdukID"`
    User      User    `gorm:"foreignKey:UserID"`
}