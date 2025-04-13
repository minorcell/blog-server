package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	// Username is unique and used as the primary key for the user
	UserID   uint   `json:"user_id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`

	// Password is hashed with bcrypt
	Password string `json:"password"`

	// 1 : male 0: female
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`

	// 1: admin 2: user
	Role     string `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}