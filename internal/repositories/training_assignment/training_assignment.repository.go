package repositories_trainingassignment

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type TrainingAssignmentRepository struct {
	db *gorm.DB
}

func NewTrainingAssignmentRepository(db *gorm.DB) *TrainingAssignmentRepository {
	return &TrainingAssignmentRepository{db: db}
}

func (r *TrainingAssignmentRepository) Create(assignment *models.TrainingAssignment) (*models.TrainingAssignment, error) {
	if err := r.db.Create(assignment).Error; err != nil {
		return nil, err
	}

	var created models.TrainingAssignment
	err := r.db.
		Preload("User").
		Preload("Course").
		Preload("AssignedByUser").
		First(&created, assignment.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *TrainingAssignmentRepository) FindAll() ([]models.TrainingAssignment, error) {
	var assignments []models.TrainingAssignment

	err := r.db.
		Preload("User").
		Preload("Course").
		Preload("AssignedByUser").
		Order("id DESC").
		Find(&assignments).
		Error

	return assignments, err
}

func (r *TrainingAssignmentRepository) FindByID(id uint) (*models.TrainingAssignment, error) {
	var assignment models.TrainingAssignment

	err := r.db.
		Preload("User").
		Preload("Course").
		Preload("AssignedByUser").
		First(&assignment, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (r *TrainingAssignmentRepository) FindByUserID(userID uint) ([]models.TrainingAssignment, error) {
	var assignments []models.TrainingAssignment

	err := r.db.
		Preload("User").
		Preload("Course").
		Preload("AssignedByUser").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&assignments).
		Error

	return assignments, err
}

func (r *TrainingAssignmentRepository) ExistsByUserAndCourse(userID uint, courseID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.TrainingAssignment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TrainingAssignmentRepository) Update(assignment *models.TrainingAssignment) (*models.TrainingAssignment, error) {
	if err := r.db.Save(assignment).Error; err != nil {
		return nil, err
	}

	var updated models.TrainingAssignment
	err := r.db.
		Preload("User").
		Preload("Course").
		Preload("AssignedByUser").
		First(&updated, assignment.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *TrainingAssignmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.TrainingAssignment{}, id).Error
}



func (r *TrainingAssignmentRepository) FindDepartmentAssignments() ([]models.TrainingAssignment, error) {
	var assignments []models.TrainingAssignment

	err := r.db.
		Preload("User").
		Preload("Course").
		Preload("Department").
		Preload("AssignedByUser").
		Where("assignment_source = ?", "department").
		Order("id DESC").
		Find(&assignments).Error

	if err != nil {
		return nil, err
	}

	return assignments, nil
}