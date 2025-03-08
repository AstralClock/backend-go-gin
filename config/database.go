package config

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := "host=localhost user=postgres password=postgres dbname=revabajuanak port=5432 sslmode=disable TimeZone=Asia/Jakarta"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Gagal konek ke database")
    }

    DB = db
    fmt.Println("Berhasil konek ke database!")
}