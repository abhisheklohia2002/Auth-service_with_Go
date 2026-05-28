package repositories

import (
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