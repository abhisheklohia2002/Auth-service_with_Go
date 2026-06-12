package repositories_department

import (
	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(department *models.Department) error
	GetByID(id uint) (*models.Department, error)
	GetAll() ([]models.Department, error)
	Update(department *models.Department) error
	Delete(id uint) error
	GetActiveUsers(departmentID uint) ([]models.User, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *departmentRepository) GetByID(id uint) (*models.Department, error) {
	var department models.Department
	err := r.db.Preload("Entity").First(&department, id).Error
	if err != nil {
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) GetAll() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Preload("Entity").Order("id DESC").Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) Update(department *models.Department) error {
	return r.db.Save(department).Error
}

func (r *departmentRepository) Delete(id uint) error {
	return r.db.Model(&models.Department{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *departmentRepository) GetActiveUsers(departmentID uint) ([]models.User, error) {
	var users []models.User

	err := r.db.
		Where("department_id = ? AND status = ?", departmentID, "active").
		Find(&users).Error

	return users, err
}
