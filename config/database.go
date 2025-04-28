package config

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

// func ConnectDB() {
//     dsn := "host=localhost user=postgres password=123 dbname=revabajuanak port=5433 sslmode=disable TimeZone=Asia/Jakarta"
//     db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//     if err != nil {
//         panic("Gagal konek ke database")
//     }

//     DB = db
//     fmt.Println("Berhasil konek ke database!")
// }

func ConnectDB() {
    dsn := "host=ep-wandering-pine-a1atgr3h-pooler.ap-southeast-1.aws.neon.tech user=neondb_owner password=npg_hBS4VIZa7FyU dbname=neondb port=5432 sslmode=require TimeZone=Asia/Jakarta"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to the database")
    }

    DB = db
    fmt.Println("Successfully connected to the database!")
}

