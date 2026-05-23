package repositories

import (
	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (d *DepartmentRepository) CreateDepartment(department *models.Department) (models.Department, error) {
	if err := d.db.Create(department).Error; err != nil {
		return models.Department{}, err
	}

	return *department, nil

}

func (d *DepartmentRepository) ListDepartments() ([]models.Department, error) {
	var listDepartment []models.Department
	err := d.db.Find(&listDepartment).Error
	if err != nil {
		return nil, err
	}
	return listDepartment, nil
}

func (d *DepartmentRepository) UserIdByDepartment(userId int) ([]models.Department, error) {
	var departmentLists []models.Department

	err := d.db.
		Preload("User").
		Where("user_id = ?", userId).
		Find(&departmentLists).Error

	return departmentLists, err
}
