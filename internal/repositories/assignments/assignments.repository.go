package repositories_assignments

import (
	"errors"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type AssignmentRuleRepository struct {
	db *gorm.DB
}

func NewAssignmentRuleRepository(db *gorm.DB) *AssignmentRuleRepository {
	return &AssignmentRuleRepository{db: db}
}

func (r *AssignmentRuleRepository) Create(rule *models.AssignmentRule) (*models.AssignmentRule, error) {
	if err := r.db.Create(rule).Error; err != nil {
		return nil, err
	}

	var created models.AssignmentRule
	err := r.db.
		Preload("Course").
		Preload("Role").
		First(&created, rule.ID).
		Error

	return &created, err
}

func (r *AssignmentRuleRepository) FindAll() ([]models.AssignmentRule, error) {
	var rules []models.AssignmentRule

	err := r.db.
		Preload("Course").
		Preload("Role").
		Order("id DESC").
		Find(&rules).
		Error

	return rules, err
}

func (r *AssignmentRuleRepository) FindByID(id uint) (*models.AssignmentRule, error) {
	var rule models.AssignmentRule

	err := r.db.
		Preload("Course").
		Preload("Role").
		First(&rule, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &rule, err
}

func (r *AssignmentRuleRepository) Update(rule *models.AssignmentRule) (*models.AssignmentRule, error) {
	if err := r.db.Save(rule).Error; err != nil {
		return nil, err
	}

	var updated models.AssignmentRule
	err := r.db.
		Preload("Course").
		Preload("Role").
		First(&updated, rule.ID).
		Error

	return &updated, err
}

func (r *AssignmentRuleRepository) Delete(id uint) error {
	return r.db.Delete(&models.AssignmentRule{}, id).Error
}

func (r *AssignmentRuleRepository) FindByCourseID(courseID uint) ([]models.AssignmentRule, error) {
	var rules []models.AssignmentRule

	err := r.db.
		Preload("Course").
		Preload("Role").
		Where("course_id = ?", courseID).
		Order("id DESC").
		Find(&rules).
		Error

	return rules, err
}
