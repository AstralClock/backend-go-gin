package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID          uint       `gorm:"primaryKey" json:"id"`
	Email       string     `json:"email" gorm:"unique;not null"`
	Password    string     `json:"-" gorm:"not null"`
	IsVerified  bool       `json:"is_verified" gorm:"default:false"`
	VerifyToken string     `json:"verify_token"`
	UserDetail  UserDetail `gorm:"foreignKey:UserID"` // Relasi one-to-one
}
