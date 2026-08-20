package seed

import (
	"omniforge-api/internal/role"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
	roles := []role.Role{
		{Name: role.NameAdmin},
		{Name: role.NameUser},
	}

	for _, item := range roles {
		var existing role.Role

		err := db.Where("name = ?", item.Name).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&item).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}