package repositories_CertificateIssue

import (
	"errors"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type CertificateIssueRepository struct {
	db *gorm.DB
}

func NewCertificateIssueRepository(db *gorm.DB) *CertificateIssueRepository {
	return &CertificateIssueRepository{db: db}
}

func (r *CertificateIssueRepository) Create(issue *models.CertificateIssue) (*models.CertificateIssue, error) {
	if err := r.db.Create(issue).Error; err != nil {
		return nil, err
	}

	var created models.CertificateIssue
	err := r.db.
		Preload("Certification.Course").
		Preload("Certification.Rule").
		Preload("User.Role").
		Preload("TrainingAssignment").
		First(&created, issue.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *CertificateIssueRepository) FindByID(id uint) (*models.CertificateIssue, error) {
	var issue models.CertificateIssue

	err := r.db.
		Preload("Certification.Course").
		Preload("Certification.Rule").
		Preload("User.Role").
		Preload("TrainingAssignment").
		First(&issue, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &issue, nil
}

func (r *CertificateIssueRepository) FindByUserID(userID uint) ([]models.CertificateIssue, error) {
	var issues []models.CertificateIssue

	err := r.db.
		Preload("Certification.Course").
		Preload("Certification.Rule").
		Preload("User.Role").
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&issues).
		Error

	return issues, err
}

func (r *CertificateIssueRepository) ExistsByUserAndCertification(userID uint, certificationID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.CertificateIssue{}).
		Where("user_id = ? AND certification_id = ? AND issue_status = ?", userID, certificationID, "issued").
		Count(&count).
		Error

	return count > 0, err
}

func (r *CertificateIssueRepository) Update(issue *models.CertificateIssue) (*models.CertificateIssue, error) {
	if err := r.db.Save(issue).Error; err != nil {
		return nil, err
	}

	return r.FindByID(issue.ID)
}

func (r *CertificateIssueRepository) FindByCertificateNumber(certificateNumber string) (*models.CertificateIssue, error) {
	var issue models.CertificateIssue

	err := r.db.
		Preload("Certification.Course").
		Preload("Certification.Rule").
		Preload("User.Role").
		Preload("TrainingAssignment").
		Where("certificate_number = ?", certificateNumber).
		First(&issue).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &issue, nil
}
