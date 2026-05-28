package repositories_assessmentquestionoption

import (
	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type AssessmentQuestionOptionRepository struct {
	db *gorm.DB
}

func NewAssessmentQuestionOptionRepository(db *gorm.DB) *AssessmentQuestionOptionRepository {
	return &AssessmentQuestionOptionRepository{db: db}
}

func (r *AssessmentQuestionOptionRepository) CreateMany(options []models.AssessmentQuestionOption) error {
	if len(options) == 0 {
		return nil
	}

	return r.db.Create(&options).Error
}

func (r *AssessmentQuestionOptionRepository) DeleteByQuestionID(questionID uint) error {
	return r.db.
		Where("question_id = ?", questionID).
		Delete(&models.AssessmentQuestionOption{}).
		Error
}
