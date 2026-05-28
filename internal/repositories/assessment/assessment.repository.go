package repositories_assessment

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type AssessmentRepository struct {
	db *gorm.DB
}

func NewAssessmentRepository(db *gorm.DB) *AssessmentRepository {
	return &AssessmentRepository{db: db}
}

func (r *AssessmentRepository) Create(assessment *models.Assessment) (*models.Assessment, error) {
	if err := r.db.Create(assessment).Error; err != nil {
		return nil, err
	}

	var created models.Assessment
	err := r.db.
		Preload("Course").
		Preload("Module").
		Preload("Rule").
		First(&created, assessment.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *AssessmentRepository) FindAll() ([]models.Assessment, error) {
	var assessments []models.Assessment

	err := r.db.
		Preload("Course").
		Preload("Module").
		Preload("Rule").
		Order("id DESC").
		Find(&assessments).
		Error

	return assessments, err
}

func (r *AssessmentRepository) FindByID(id uint) (*models.Assessment, error) {
	var assessment models.Assessment

	err := r.db.
		Preload("Course").
		Preload("Module").
		Preload("Rule").
		First(&assessment, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &assessment, nil
}

func (r *AssessmentRepository) FindByCourseID(courseID uint) ([]models.Assessment, error) {
	var assessments []models.Assessment

	err := r.db.
		Preload("Course").
		Preload("Module").
		Preload("Rule").
		Where("course_id = ?", courseID).
		Order("id DESC").
		Find(&assessments).
		Error

	return assessments, err
}

func (r *AssessmentRepository) FindByModuleID(moduleID uint) ([]models.Assessment, error) {
	var assessments []models.Assessment

	err := r.db.
		Preload("Course").
		Preload("Module").
		Preload("Rule").
		Where("module_id = ?", moduleID).
		Order("id DESC").
		Find(&assessments).
		Error

	return assessments, err
}

func (r *AssessmentRepository) Update(assessment *models.Assessment) (*models.Assessment, error) {
	if err := r.db.Save(assessment).Error; err != nil {
		return nil, err
	}

	var updated models.Assessment
	err := r.db.
		Preload("Course").
		Preload("Module").
		Preload("Rule").
		First(&updated, assessment.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *AssessmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Assessment{}, id).Error
}
