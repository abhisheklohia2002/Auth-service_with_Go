package repositories_assessmentrule


import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type AssessmentRuleRepository struct {
	db *gorm.DB
}

func NewAssessmentRuleRepository(db *gorm.DB) *AssessmentRuleRepository {
	return &AssessmentRuleRepository{db: db}
}

func (r *AssessmentRuleRepository) Create(rule *models.AssessmentRule) (*models.AssessmentRule, error) {
	if err := r.db.Create(rule).Error; err != nil {
		return nil, err
	}

	return rule, nil
}

func (r *AssessmentRuleRepository) FindAll() ([]models.AssessmentRule, error) {
	var rules []models.AssessmentRule

	err := r.db.
		Order("id DESC").
		Find(&rules).
		Error

	return rules, err
}

func (r *AssessmentRuleRepository) FindByID(id uint) (*models.AssessmentRule, error) {
	var rule models.AssessmentRule

	err := r.db.First(&rule, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (r *AssessmentRuleRepository) Update(rule *models.AssessmentRule) (*models.AssessmentRule, error) {
	if err := r.db.Save(rule).Error; err != nil {
		return nil, err
	}

	return rule, nil
}

func (r *AssessmentRuleRepository) Delete(id uint) error {
	return r.db.Delete(&models.AssessmentRule{}, id).Error
}