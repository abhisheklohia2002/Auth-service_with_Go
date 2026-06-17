package repositories_certificationrule

import (
	"errors"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type CertificationRuleRepository struct {
	db *gorm.DB
}

func NewCertificationRuleRepository(db *gorm.DB) *CertificationRuleRepository {
	return &CertificationRuleRepository{db: db}
}

func (r *CertificationRuleRepository) Create(
	rule *models.CertificationRule,
) (*models.CertificationRule, error) {
	if err := r.db.Create(rule).Error; err != nil {
		return nil, err
	}

	var created models.CertificationRule

	err := r.db.
		Preload("Course").
		First(&created, rule.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *CertificationRuleRepository) FindAll() ([]models.CertificationRule, error) {
	var rules []models.CertificationRule

	err := r.db.
		Preload("Course").
		Order("id DESC").
		Find(&rules).
		Error

	return rules, err
}

func (r *CertificationRuleRepository) FindByID(id uint) (*models.CertificationRule, error) {
	var rule models.CertificationRule

	err := r.db.First(&rule, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &rule, nil
}

func (r *CertificationRuleRepository) Update(rule *models.CertificationRule) (*models.CertificationRule, error) {
	if err := r.db.Save(rule).Error; err != nil {
		return nil, err
	}
	return rule, nil
}

func (r *CertificationRuleRepository) Delete(id uint) error {
	return r.db.Delete(&models.CertificationRule{}, id).Error
}

func (r *CertificationRuleRepository) FindByCourseID(
	courseID uint,
) (*models.CertificationRule, error) {
	var rule models.CertificationRule

	err := r.db.
		Preload("Course").
		Where("course_id = ?", courseID).
		First(&rule).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &rule, nil
}
