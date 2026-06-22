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

func (r *ModuleRepository) UpsertByModuleID(
	video *models.ModuleVideo,
) (*models.ModuleVideo, error) {
	var existingVideo models.ModuleVideo

	err := r.db.
		Where("module_id = ?", video.ModuleID).
		First(&existingVideo).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.Create(video).Error; err != nil {
				return nil, err
			}

			return video, nil
		}

		return nil, err
	}

	existingVideo.CourseID = video.CourseID
	existingVideo.Title = video.Title
	existingVideo.VideoName = video.VideoName
	existingVideo.VideoURL = video.VideoURL
	existingVideo.VideoPublicID = video.VideoPublicID
	existingVideo.VideoSize = video.VideoSize
	existingVideo.VideoType = video.VideoType
	existingVideo.IsActive = video.IsActive

	if err := r.db.Save(&existingVideo).Error; err != nil {
		return nil, err
	}

	return &existingVideo, nil
}

func (r *ModuleRepository) FindByModuleID(
	moduleID uint,
) (*models.ModuleVideo, error) {
	var video models.ModuleVideo

	err := r.db.
		Where("module_id = ? AND is_active = ?", moduleID, true).
		First(&video).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &video, nil
}
