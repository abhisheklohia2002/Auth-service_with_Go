package repositories_moduleDocument

import (
	"errors"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

type ModuleDocumentRepository struct {
	db *gorm.DB
}

func NewModuleDocumentRepository(db *gorm.DB) *ModuleDocumentRepository {
	return &ModuleDocumentRepository{db: db}
}

func (r *ModuleDocumentRepository) Create(document *models.ModuleDocument) (*models.ModuleDocument, error) {
	if err := r.db.Create(document).Error; err != nil {
		return nil, err
	}

	var created models.ModuleDocument

	err := r.db.
		Preload("Module").
		First(&created, document.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ModuleDocumentRepository) FindByID(id uint) (*models.ModuleDocument, error) {
	var document models.ModuleDocument

	err := r.db.
		Preload("Module").
		First(&document, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (r *ModuleDocumentRepository) FindByModuleID(moduleID uint) ([]models.ModuleDocument, error) {
	var documents []models.ModuleDocument

	err := r.db.
		Where("module_id = ? AND is_active = ?", moduleID, true).
		Order("created_at ASC").
		Find(&documents).
		Error

	return documents, err
}

func (r *ModuleDocumentRepository) Update(document *models.ModuleDocument) (*models.ModuleDocument, error) {
	if err := r.db.Save(document).Error; err != nil {
		return nil, err
	}

	var updated models.ModuleDocument

	err := r.db.
		Preload("Module").
		First(&updated, document.ID).
		Error

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ModuleDocumentRepository) Delete(id uint) error {
	return r.db.Delete(&models.ModuleDocument{}, id).Error
}

func (r *ModuleDocumentRepository) DeleteDocumentByPublicId(publicID string) error {
	return r.db.
		Where("public_id = ?", publicID).
		Delete(&models.ModuleDocument{}).
		Error
}

func (r *ModuleDocumentRepository) ExistsByModuleID(moduleID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.ModuleVideo{}).
		Where("module_id = ? AND video_url <> ''", moduleID).
		Count(&count).
		Error

	return count > 0, err
}
