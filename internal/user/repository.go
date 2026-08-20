package user

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User

	err := r.db.
		Preload("Role").
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByID(id uint) (*User, error) {
	var user User

	err := r.db.
		Preload("Role").
		First(&user, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindAll() ([]User, error) {
	var users []User

	if err := r.db.
		Preload("Role").
		Find(&users).
		Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *Repository) Update(user *User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}

	if err := r.db.
		Preload("Role").
		First(user, user.ID).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(
		&User{},
		id,
	).Error
}