package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex"`
	Username     string    `gorm:uniqueIndex;not null"`
	PasswordHash string    `gorm" "not null" json:"-"`
	Stworzony    time.Time `gorm:"autoCreateTime`
}
