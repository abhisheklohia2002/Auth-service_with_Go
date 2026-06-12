package services_entity

import (
	"errors"
	"strings"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_entity "example.com/m/internal/repositories/entity"
)

type EntityService interface {
	Create(req dto.CreateEntityRequest) (*models.Entity, error)
	GetByID(id uint) (*models.Entity, error)
	GetAll() ([]models.Entity, error)
	Update(id uint, req dto.UpdateEntityRequest) (*models.Entity, error)
	Delete(id uint) error
}

type entityService struct {
	entityRepo repositories_entity.EntityRepository
}

func NewEntityService(entityRepo repositories_entity.EntityRepository) EntityService {
	return &entityService{entityRepo: entityRepo}
}

func (s *entityService) Create(req dto.CreateEntityRequest) (*models.Entity, error) {
	name := strings.TrimSpace(req.EntityName)

	if name == "" {
		return nil, errors.New("entity name is required")
	}

	exists, err := s.entityRepo.ExistsByName(name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("entity already exists")
	}

	entity := &models.Entity{
		EntityName:  name,
		EntityType:  strings.TrimSpace(req.EntityType),
		Description: strings.TrimSpace(req.Description),
		IsActive:    true,
	}

	if err := s.entityRepo.Create(entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *entityService) GetByID(id uint) (*models.Entity, error) {
	entity, err := s.entityRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, errors.New("entity not found")
	}

	return entity, nil
}

func (s *entityService) GetAll() ([]models.Entity, error) {
	return s.entityRepo.GetAll()
}

func (s *entityService) Update(id uint, req dto.UpdateEntityRequest) (*models.Entity, error) {
	entity, err := s.entityRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, errors.New("entity not found")
	}

	if strings.TrimSpace(req.EntityName) != "" {
		entity.EntityName = strings.TrimSpace(req.EntityName)
	}

	if strings.TrimSpace(req.EntityType) != "" {
		entity.EntityType = strings.TrimSpace(req.EntityType)
	}

	if strings.TrimSpace(req.Description) != "" {
		entity.Description = strings.TrimSpace(req.Description)
	}

	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}

	if err := s.entityRepo.Update(entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *entityService) Delete(id uint) error {
	entity, err := s.entityRepo.GetByID(id)
	if err != nil {
		return err
	}

	if entity == nil {
		return errors.New("entity not found")
	}

	return s.entityRepo.Delete(id)
}