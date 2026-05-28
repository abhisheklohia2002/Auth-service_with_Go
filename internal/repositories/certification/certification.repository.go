package repositories_certification

import (
	"errors"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type CertificationRepository struct {
	db *gorm.DB
}

func NewCertificationRepository(db *gorm.DB) *CertificationRepository {
	return &CertificationRepository{db: db}
}

func (r *CertificationRepository) Create(cert *models.Certification) (*models.Certification, error) {
	if err := r.db.Create(cert).Error; err != nil {
		return nil, err
	}

	var created models.Certification
	err := r.db.
		Preload("Course").
		Preload("Rule").
		First(&created, cert.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *CertificationRepository) FindAll() ([]models.Certification, error) {
	var certs []models.Certification

	err := r.db.
		Preload("Course").
		Preload("Rule").
		Order("id DESC").
		Find(&certs).
		Error

	return certs, err
}

func (r *CertificationRepository) FindByID(id uint) (*models.Certification, error) {
	var cert models.Certification

	err := r.db.
		Preload("Course").
		Preload("Rule").
		First(&cert, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (r *CertificationRepository) FindByCourseID(courseID uint) ([]models.Certification, error) {
	var certs []models.Certification

	err := r.db.
		Preload("Course").
		Preload("Rule").
		Where("course_id = ? AND is_active = ?", courseID, true).
		Order("id DESC").
		Find(&certs).
		Error

	return certs, err
}

func (r *CertificationRepository) Update(cert *models.Certification) (*models.Certification, error) {
	if err := r.db.Save(cert).Error; err != nil {
		return nil, err
	}

	var updated models.Certification
	err := r.db.
		Preload("Course").
		Preload("Rule").
		First(&updated, cert.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *CertificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Certification{}, id).Error
}