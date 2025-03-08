package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    ID       uint   `gorm:"primaryKey" json:"id"`
    Email    string `json:"email" gorm:"unique;not null"`
    Password string `json:"-" gorm:"not null"`
    UserDetail UserDetail `gorm:"foreignKey:UserID"` // Relasi one-to-one
}