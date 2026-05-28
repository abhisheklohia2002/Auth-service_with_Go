package services_role

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	
	repositories_role "example.com/m/internal/repositories/Role"
)

type RoleService struct {
	roleRepo *repositories_role.RoleRepository
}

func NewRoleService(roleRepo *repositories_role.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) Create(req dto.CreateRoleRequest) (*models.Role, error) {
	existing, err := s.roleRepo.FindByName(req.RoleName)
	if err == nil && existing.ID != 0 {
		return nil, errors.New("role already exists")
	}

	role := models.Role{
		RoleName:    req.RoleName,
		RoleType:    req.RoleType,
		Description: req.Description,
	}

	return s.roleRepo.Create(&role)
}

func (s *RoleService) FindAll() ([]models.Role, error) {
	return s.roleRepo.FindAll()
}

func (s *RoleService) FindByID(id uint) (*models.Role, error) {
	return s.roleRepo.FindByID(id)
}

func (s *RoleService) Update(id uint, req dto.UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	if req.RoleName != nil {
		role.RoleName = *req.RoleName
	}

	if req.RoleType != nil {
		role.RoleType = *req.RoleType
	}

	if req.Description != nil {
		role.Description = *req.Description
	}

	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	return s.roleRepo.Delete(id)
}
