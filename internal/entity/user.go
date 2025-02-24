package entity

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password,omitempty" gorm:"not null"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// TODO: Tambahkan field sesuai kebutuhan
}
