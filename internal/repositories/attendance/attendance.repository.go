package repositories_attendance

import (
	"errors"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Create(attendance *models.Attendance) error
	BulkCreate(attendances []models.Attendance) error
	FindByID(id uint) (*models.Attendance, error)
	FindBySessionID(sessionID uint) ([]models.Attendance, error)
	FindByUserID(userID uint) ([]models.Attendance, error)
	FindBySessionAndUser(sessionID uint, userID uint) (*models.Attendance, error)
	Update(attendance *models.Attendance) error
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) BulkCreate(attendances []models.Attendance) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, attendance := range attendances {
			if err := tx.Create(&attendance).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *attendanceRepository) FindByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance

	if err := r.db.First(&attendance, id).Error; err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (r *attendanceRepository) FindBySessionID(sessionID uint) ([]models.Attendance, error) {
	var attendances []models.Attendance

	err := r.db.
		Where("session_id = ?", sessionID).
		Order("user_id ASC").
		Find(&attendances).Error

	return attendances, err
}

func (r *attendanceRepository) FindByUserID(userID uint) ([]models.Attendance, error) {
	var attendances []models.Attendance

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&attendances).Error

	return attendances, err
}

func (r *attendanceRepository) FindBySessionAndUser(sessionID uint, userID uint) (*models.Attendance, error) {
	var attendance models.Attendance

	err := r.db.
		Where("session_id = ? AND user_id = ?", sessionID, userID).
		First(&attendance).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func (r *attendanceRepository) Update(attendance *models.Attendance) error {
	return r.db.Save(attendance).Error
}