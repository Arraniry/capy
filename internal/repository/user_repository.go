package repository

import (
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetAll() ([]User, error) {
	var items []User
	result := r.db.Find(&items)
	return items, result.Error
}

func (r *UserRepository) GetByID(id uint) (*User, error) {
	var item User
	result := r.db.First(&item, id)
	return &item, result.Error
}

func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
