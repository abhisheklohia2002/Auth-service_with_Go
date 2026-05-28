package repositories_course

import (
	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (r *CourseRepository) Create(course *models.Course) (*models.Course, error) {
	if err := r.db.Create(course).Error; err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) FindAll() ([]models.Course, error) {
	var courses []models.Course

	err := r.db.
		Preload("CreatedByUser").
		Order("created_at DESC").
		Find(&courses).
		Error

	return courses, err
}

func (r *CourseRepository) FindByID(id uint) (*models.Course, error) {
	var course models.Course

	err := r.db.
		Preload("CreatedByUser").
		Preload("Modules").
		First(&course, id).
		Error

	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepository) Update(course *models.Course) (*models.Course, error) {
	if err := r.db.Save(course).Error; err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}
