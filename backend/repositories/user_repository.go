package repositories

import (
	"churma-keygen/backend/domain"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) CountByUsername(username string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("username = ?", username).Count(&count).Error
	return count, err
}

func (r *GormUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
