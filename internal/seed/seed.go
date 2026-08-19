package seed

import "gorm.io/gorm"

func Run(db *gorm.DB) error {

	if err:= SeedRoles(db); err != nil{
		return err
	}

	return nil
}