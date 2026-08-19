package role

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:50;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}