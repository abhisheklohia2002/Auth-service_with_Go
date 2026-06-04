package repositories_department_training_mapping

import (
	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type DepartmentTrainingMappingRepository interface {
	Create(mapping *models.DepartmentTrainingMapping) error
	GetByID(id uint) (*models.DepartmentTrainingMapping, error)
	GetAll() ([]models.DepartmentTrainingMapping, error)
	GetByDepartmentID(departmentID uint) ([]models.DepartmentTrainingMapping, error)
	Exists(departmentID uint, courseID uint) (bool, error)
}

type departmentTrainingMappingRepository struct {
	db *gorm.DB
}

func NewDepartmentTrainingMappingRepository(db *gorm.DB) DepartmentTrainingMappingRepository {
	return &departmentTrainingMappingRepository{db: db}
}

func (r *departmentTrainingMappingRepository) Create(mapping *models.DepartmentTrainingMapping) error {
	return r.db.Create(mapping).Error
}

func (r *departmentTrainingMappingRepository) GetByID(id uint) (*models.DepartmentTrainingMapping, error) {
	var mapping models.DepartmentTrainingMapping
	err := r.db.
		Preload("Department").
		Preload("Course").
		First(&mapping, id).Error

	if err != nil {
		return nil, err
	}

	return &mapping, nil
}

func (r *departmentTrainingMappingRepository) GetAll() ([]models.DepartmentTrainingMapping, error) {
	var mappings []models.DepartmentTrainingMapping

	err := r.db.
		Preload("Department").
		Preload("Course").
		Order("id DESC").
		Find(&mappings).Error

	return mappings, err
}

func (r *departmentTrainingMappingRepository) GetByDepartmentID(departmentID uint) ([]models.DepartmentTrainingMapping, error) {
	var mappings []models.DepartmentTrainingMapping

	err := r.db.
		Preload("Department").
		Preload("Course").
		Where("department_id = ? AND active_flag = ?", departmentID, true).
		Find(&mappings).Error

	if err != nil {
		return nil, err
	}

	return mappings, nil
}

func (r *departmentTrainingMappingRepository) Exists(departmentID uint, courseID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.DepartmentTrainingMapping{}).
		Where("department_id = ? AND course_id = ?", departmentID, courseID).
		Count(&count).Error

	return count > 0, err
}
