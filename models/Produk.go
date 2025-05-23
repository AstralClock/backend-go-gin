package models

import "gorm.io/gorm"

type Produk struct {
    gorm.Model
    NamaProduk string    `json:"nama_produk" gorm:"not null"`
    Deskripsi  string    `json:"deskripsi" gorm:"not null"`
    Kategori   string    `json:"kategori" gorm:"not null"`
    Tag        string    `json:"tag" gorm:"not null"`
    Harga      float64   `json:"harga" gorm:"not null"`
    Jumlah     int       `json:"jumlah" gorm:"not null"`
    Image      string    `json:"image" gorm:"not null"`
    Ukurans    []Ukuran  `gorm:"many2many:produk_ukurans;"` // Perubahan utama di sini
    ProdukUkuranStock []ProdukUkuran   `gorm:"foreignKey:ProdukID"`
    Ulasan []Ulasan `gorm:"foreignKey:ProdukID"`
}

type Model struct {
    ID        uint           `json:"id" gorm:"primarykey"`
}