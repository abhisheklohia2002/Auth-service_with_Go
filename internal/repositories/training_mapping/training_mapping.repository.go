package repository_trainingmapping

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type TrainingMappingRepository struct {
	db *gorm.DB
}

func NewTrainingMappingRepository(db *gorm.DB) *TrainingMappingRepository {
	return &TrainingMappingRepository{
		db: db,
	}
}

func (r *TrainingMappingRepository) Create(mapping *models.TrainingMapping) (*models.TrainingMapping, error) {
	if err := r.db.Create(mapping).Error; err != nil {
		return nil, err
	}

	var created models.TrainingMapping
	err := r.db.
		Preload("Role").
		Preload("Course").
		First(&created, mapping.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *TrainingMappingRepository) FindAll() ([]models.TrainingMapping, error) {
	var mappings []models.TrainingMapping

	err := r.db.
		Preload("Role").
		Preload("Course").
		Order("id DESC").
		Find(&mappings).
		Error

	return mappings, err
}

func (r *TrainingMappingRepository) FindByID(id uint) (*models.TrainingMapping, error) {
	var mapping models.TrainingMapping

	err := r.db.
		Preload("Role").
		Preload("Course").
		First(&mapping, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &mapping, nil
}

func (r *TrainingMappingRepository) FindByRoleID(roleID uint) ([]models.TrainingMapping, error) {
	var mappings []models.TrainingMapping

	err := r.db.
		Preload("Role").
		Preload("Course").
		Where("role_id = ?", roleID).
		Order("id DESC").
		Find(&mappings).
		Error

	return mappings, err
}

func (r *TrainingMappingRepository) ExistsByRoleAndCourse(roleID uint, courseID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.TrainingMapping{}).
		Where("role_id = ? AND course_id = ?", roleID, courseID).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TrainingMappingRepository) Update(mapping *models.TrainingMapping) (*models.TrainingMapping, error) {
	if err := r.db.Save(mapping).Error; err != nil {
		return nil, err
	}

	var updated models.TrainingMapping
	err := r.db.
		Preload("Role").
		Preload("Course").
		First(&updated, mapping.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *TrainingMappingRepository) Delete(id uint) error {
	return r.db.Delete(&models.TrainingMapping{}, id).Error
}

func (r *TrainingMappingRepository) FindActiveByRoleID(roleID uint) ([]models.TrainingMapping, error) {
	var mappings []models.TrainingMapping

	err := r.db.
		Preload("Course").
		Where("role_id = ? AND active_flag = ?", roleID, true).
		Find(&mappings).
		Error

	return mappings, err
}
