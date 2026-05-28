package repositories_moduleprogress

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type ModuleProgressRepository struct {
	db *gorm.DB
}

func NewModuleProgressRepository(db *gorm.DB) *ModuleProgressRepository {
	return &ModuleProgressRepository{db: db}
}

func (r *ModuleProgressRepository) Create(progress *models.ModuleProgress) (*models.ModuleProgress, error) {
	if err := r.db.Create(progress).Error; err != nil {
		return nil, err
	}

	var created models.ModuleProgress
	err := r.db.
		Preload("Module").
		Preload("Assignment").
		First(&created, progress.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ModuleProgressRepository) BulkCreate(progresses []models.ModuleProgress) error {
	if len(progresses) == 0 {
		return nil
	}

	return r.db.Create(&progresses).Error
}

func (r *ModuleProgressRepository) FindByAssignmentID(assignmentID uint) ([]models.ModuleProgress, error) {
	var progresses []models.ModuleProgress
	
	err := r.db.
		Preload("Module").
		Where("assignment_id = ?", assignmentID).
		Order("id ASC").
		Find(&progresses).
		Error

	return progresses, err
}

func (r *ModuleProgressRepository) FindByID(id uint) (*models.ModuleProgress, error) {
	var progress models.ModuleProgress

	err := r.db.
		Preload("Module").
		Preload("Assignment").
		First(&progress, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *ModuleProgressRepository) Update(progress *models.ModuleProgress) (*models.ModuleProgress, error) {
	if err := r.db.Save(progress).Error; err != nil {
		return nil, err
	}

	var updated models.ModuleProgress
	err := r.db.
		Preload("Module").
		Preload("Assignment").
		First(&updated, progress.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ModuleProgressRepository) CountByAssignmentAndStatus(assignmentID uint, status string) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.ModuleProgress{}).
		Where("assignment_id = ? AND status = ?", assignmentID, status).
		Count(&count).
		Error

	return count, err
}

func (r *ModuleProgressRepository) CountByAssignment(assignmentID uint) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.ModuleProgress{}).
		Where("assignment_id = ?", assignmentID).
		Count(&count).
		Error

	return count, err
}

func (r *ModuleProgressRepository) ExistsByAssignmentAndModule(assignmentID uint, moduleID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.ModuleProgress{}).
		Where("assignment_id = ? AND module_id = ?", assignmentID, moduleID).
		Count(&count).
		Error

	return count > 0, err
}
