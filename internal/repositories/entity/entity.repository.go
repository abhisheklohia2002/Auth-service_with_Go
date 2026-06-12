package repositories_entity

import (
	"errors"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type EntityRepository interface {
	Create(entity *models.Entity) error
	GetByID(id uint) (*models.Entity, error)
	GetAll() ([]models.Entity, error)
	Update(entity *models.Entity) error
	Delete(id uint) error
	ExistsByName(name string) (bool, error)
}

type entityRepository struct {
	db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) EntityRepository {
	return &entityRepository{db: db}
}

func (r *entityRepository) Create(entity *models.Entity) error {
	return r.db.Create(entity).Error
}

func (r *entityRepository) GetByID(id uint) (*models.Entity, error) {
	var entity models.Entity

	err := r.db.
		Preload("Departments").
		First(&entity, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *entityRepository) GetAll() ([]models.Entity, error) {
	var entities []models.Entity

	err := r.db.
		Preload("Departments").
		Order("id DESC").
		Find(&entities).Error

	return entities, err
}

func (r *entityRepository) Update(entity *models.Entity) error {
	return r.db.Save(entity).Error
}

func (r *entityRepository) Delete(id uint) error {
	return r.db.Model(&models.Entity{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

func (r *entityRepository) ExistsByName(name string) (bool, error) {
	var count int64

	err := r.db.Model(&models.Entity{}).
		Where("name = ?", name).
		Count(&count).Error

	return count > 0, err
}