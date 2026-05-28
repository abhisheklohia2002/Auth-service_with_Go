package repositories_module

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	return &ModuleRepository{db: db}
}

func (r *ModuleRepository) Create(module *models.Module) (*models.Module, error) {
	if err := r.db.Create(module).Error; err != nil {
		return nil, err
	}

	var created models.Module
	err := r.db.
		Preload("Course").
		First(&created, module.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ModuleRepository) FindAll() ([]models.Module, error) {
	var modules []models.Module

	err := r.db.
		Preload("Course").
		Order("course_id ASC, sequence_no ASC").
		Find(&modules).
		Error

	return modules, err
}

func (r *ModuleRepository) FindByID(id uint) (*models.Module, error) {
	var module models.Module

	err := r.db.
		Preload("Course").
		First(&module, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *ModuleRepository) FindByCourseID(courseID uint) ([]models.Module, error) {
	var modules []models.Module

	err := r.db.
		Preload("Course").
		Where("course_id = ? AND is_active = ?", courseID, true).
		Order("sequence_no ASC").
		Find(&modules).
		Error

	return modules, err
}

func (r *ModuleRepository) Update(module *models.Module) (*models.Module, error) {
	if err := r.db.Save(module).Error; err != nil {
		return nil, err
	}

	var updated models.Module
	err := r.db.
		Preload("Course").
		First(&updated, module.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ModuleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Module{}, id).Error
}
