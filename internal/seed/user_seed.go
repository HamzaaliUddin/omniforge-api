package seed

import (
	"errors"

	"omniforge-api/internal/auth"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	var userRole role.Role

	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	var existingUser user.User

	err := db.Where("email = ?", "user@omniforge.dev").First(&existingUser).Error

	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordHash, err := auth.HashPassword("Password123")
	if err != nil {
		return err
	}

	newUser := user.User{
		Name:         "Demo User",
		Email:        "user@omniforge.dev",
		PasswordHash: passwordHash,
		RoleID:       userRole.ID,
	}

	return db.Create(&newUser).Error
}