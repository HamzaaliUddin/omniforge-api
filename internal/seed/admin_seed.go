package seed

import (
	"errors"

	"omniforge-api/internal/auth"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	var adminRole role.Role

	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var existingUser user.User

	err := db.Where("email = ?", "admin@omniforge.dev").First(&existingUser).Error

	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordHash, err := auth.HashPassword("Admin123")
	if err != nil {
		return err
	}

	admin := user.User{
		Name:         "OmniForge Admin",
		Email:        "admin@omniforge.dev",
		PasswordHash: passwordHash,
		RoleID:       adminRole.ID,
	}

	return db.Create(&admin).Error
}