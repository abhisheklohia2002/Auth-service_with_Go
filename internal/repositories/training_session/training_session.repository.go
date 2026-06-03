package repositories_trainingsession

import (
	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type TrainingSessionRepository interface {
	Create(session *models.TrainingSession) error
	FindAll(filter models.TrainingSessionFilter) ([]models.TrainingSession, error)
	FindByID(id uint) (*models.TrainingSession, error)
	Update(session *models.TrainingSession) error
	Delete(session *models.TrainingSession) error
}

type trainingSessionRepository struct {
	db *gorm.DB
}

func NewTrainingSessionRepository(db *gorm.DB) TrainingSessionRepository {
	return &trainingSessionRepository{db: db}
}

func (r *trainingSessionRepository) Create(session *models.TrainingSession) error {
	return r.db.Create(session).Error
}

func (r *trainingSessionRepository) FindAll(filter models.TrainingSessionFilter) ([]models.TrainingSession, error) {
	var sessions []models.TrainingSession

	query := r.db.Model(&models.TrainingSession{})

	if filter.CourseID != "" {
		query = query.Where("course_id = ?", filter.CourseID)
	}

	if filter.ModuleID != "" {
		query = query.Where("module_id = ?", filter.ModuleID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.SessionType != "" {
		query = query.Where("session_type = ?", filter.SessionType)
	}

	err := query.Order("start_time ASC").Find(&sessions).Error
	return sessions, err
}

func (r *trainingSessionRepository) FindByID(id uint) (*models.TrainingSession, error) {
	var session models.TrainingSession

	if err := r.db.First(&session, id).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *trainingSessionRepository) Update(session *models.TrainingSession) error {
	return r.db.Save(session).Error
}

func (r *trainingSessionRepository) Delete(session *models.TrainingSession) error {
	return r.db.Delete(session).Error
}