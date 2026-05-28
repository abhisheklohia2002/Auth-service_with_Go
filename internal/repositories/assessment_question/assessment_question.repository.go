package repositories_assessmentquestion

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type AssessmentQuestionRepository struct {
	db *gorm.DB
}

func NewAssessmentQuestionRepository(db *gorm.DB) *AssessmentQuestionRepository {
	return &AssessmentQuestionRepository{db: db}
}

func (r *AssessmentQuestionRepository) Create(question *models.AssessmentQuestion) (*models.AssessmentQuestion, error) {
	if err := r.db.Create(question).Error; err != nil {
		return nil, err
	}

	return r.FindByID(question.ID)
}

func (r *AssessmentQuestionRepository) FindByID(id uint) (*models.AssessmentQuestion, error) {
	var question models.AssessmentQuestion

	err := r.db.
		Preload("Options").
		Preload("Assessment").
		First(&question, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *AssessmentQuestionRepository) FindByAssessmentID(assessmentID uint) ([]models.AssessmentQuestion, error) {
	var questions []models.AssessmentQuestion

	err := r.db.
		Preload("Options").
		Where("assessment_id = ? AND is_active = ?", assessmentID, true).
		Order("sequence_no ASC").
		Find(&questions).
		Error

	return questions, err
}

func (r *AssessmentQuestionRepository) FindAllByAssessmentID(assessmentID uint) ([]models.AssessmentQuestion, error) {
	var questions []models.AssessmentQuestion

	err := r.db.
		Preload("Options").
		Where("assessment_id = ?", assessmentID).
		Order("sequence_no ASC").
		Find(&questions).
		Error

	return questions, err
}

func (r *AssessmentQuestionRepository) Update(question *models.AssessmentQuestion) (*models.AssessmentQuestion, error) {
	if err := r.db.Save(question).Error; err != nil {
		return nil, err
	}

	return r.FindByID(question.ID)
}

func (r *AssessmentQuestionRepository) Delete(id uint) error {
	return r.db.Delete(&models.AssessmentQuestion{}, id).Error
}
