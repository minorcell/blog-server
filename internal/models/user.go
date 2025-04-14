package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID uint `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`

	Password string `json:"password"`

	// 0: unknown 1: male 2: female
	Sex     string `json:"sex"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`

	// 1: normal 2: admin
	Role string `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
