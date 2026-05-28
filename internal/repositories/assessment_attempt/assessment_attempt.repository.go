package repositories_assessmentattempt



import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type AssessmentAttemptRepository struct {
	db *gorm.DB
}

func NewAssessmentAttemptRepository(db *gorm.DB) *AssessmentAttemptRepository {
	return &AssessmentAttemptRepository{db: db}
}

func (r *AssessmentAttemptRepository) Create(attempt *models.AssessmentAttempt) (*models.AssessmentAttempt, error) {
	if err := r.db.Create(attempt).Error; err != nil {
		return nil, err
	}

	var created models.AssessmentAttempt
	err := r.db.
		Preload("Assessment").
		Preload("User.Role").
		First(&created, attempt.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *AssessmentAttemptRepository) FindByID(id uint) (*models.AssessmentAttempt, error) {
	var attempt models.AssessmentAttempt

	err := r.db.
		Preload("Assessment.Rule").
		Preload("User.Role").
		First(&attempt, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &attempt, nil
}

func (r *AssessmentAttemptRepository) FindByUserID(userID uint) ([]models.AssessmentAttempt, error) {
	var attempts []models.AssessmentAttempt

	err := r.db.
		Preload("Assessment.Rule").
		Preload("User.Role").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&attempts).
		Error

	return attempts, err
}

func (r *AssessmentAttemptRepository) FindByAssessmentID(assessmentID uint) ([]models.AssessmentAttempt, error) {
	var attempts []models.AssessmentAttempt

	err := r.db.
		Preload("Assessment.Rule").
		Preload("User.Role").
		Where("assessment_id = ?", assessmentID).
		Order("id DESC").
		Find(&attempts).
		Error

	return attempts, err
}

func (r *AssessmentAttemptRepository) FindByUserAndAssessment(userID uint, assessmentID uint) ([]models.AssessmentAttempt, error) {
	var attempts []models.AssessmentAttempt

	err := r.db.
		Preload("Assessment.Rule").
		Preload("User.Role").
		Where("user_id = ? AND assessment_id = ?", userID, assessmentID).
		Order("attempt_no ASC").
		Find(&attempts).
		Error

	return attempts, err
}

func (r *AssessmentAttemptRepository) CountByUserAndAssessment(userID uint, assessmentID uint) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.AssessmentAttempt{}).
		Where("user_id = ? AND assessment_id = ?", userID, assessmentID).
		Count(&count).
		Error

	return count, err
}

func (r *AssessmentAttemptRepository) HasPassed(userID uint, assessmentID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.AssessmentAttempt{}).
		Where("user_id = ? AND assessment_id = ? AND result_status = ?", userID, assessmentID, "passed").
		Count(&count).
		Error

	return count > 0, err
}