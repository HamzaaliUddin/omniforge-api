package seed

import "gorm.io/gorm"

func Run(db *gorm.DB) error {
	if err := SeedRoles(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedAdmin(db); err != nil {
		return err
	}

	return nil
}