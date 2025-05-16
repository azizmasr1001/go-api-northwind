package models

import "time"

type PasswordReset struct {
	ID        int       `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}
