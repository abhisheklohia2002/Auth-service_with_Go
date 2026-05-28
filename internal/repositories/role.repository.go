package repositories

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) FindByName(roleName string) (models.Role, error) {
	var role models.Role

	err := r.db.
		Where("role_name = ?", roleName).
		First(&role).
		Error

	return role, err
}

func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role

	err := r.db.First(&role, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &role, nil
}
