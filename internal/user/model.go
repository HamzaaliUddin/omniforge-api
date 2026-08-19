package user

import (
	"time"

	"omniforge-api/internal/role"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	RoleID       uint      `gorm:"not null"`
	Role         role.Role `gorm:"foreignKey:RoleID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}