package repositories_assessmentattemptanswer

import (
	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type AssessmentAttemptAnswerRepository struct {
	db *gorm.DB
}

func NewAssessmentAttemptAnswerRepository(db *gorm.DB) *AssessmentAttemptAnswerRepository {
	return &AssessmentAttemptAnswerRepository{db: db}
}

func (r *AssessmentAttemptAnswerRepository) CreateMany(answers []models.AssessmentAttemptAnswer) error {
	if len(answers) == 0 {
		return nil
	}

	return r.db.Create(&answers).Error
}

func (r *AssessmentAttemptAnswerRepository) FindByAttemptID(attemptID uint) ([]models.AssessmentAttemptAnswer, error) {
	var answers []models.AssessmentAttemptAnswer

	err := r.db.
		Preload("Question.Options").
		Where("attempt_id = ?", attemptID).
		Find(&answers).
		Error

	return answers, err
}
