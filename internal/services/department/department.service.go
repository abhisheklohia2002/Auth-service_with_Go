package services_department


import (
	"errors"
	"strings"

	"example.com/m/internal/models"
	repositories_department "example.com/m/internal/repositories/department"
)

type CreateDepartmentRequest struct {
	DepartmentName string `json:"department_name" binding:"required"`
	Description    string `json:"description"`
	IsActive       *bool  `json:"is_active"`
}

type UpdateDepartmentRequest struct {
	DepartmentName string `json:"department_name"`
	Description    string `json:"description"`
	IsActive       *bool  `json:"is_active"`
}

type DepartmentService interface {
	Create(req CreateDepartmentRequest) (*models.Department, error)
	GetByID(id uint) (*models.Department, error)
	GetAll() ([]models.Department, error)
	Update(id uint, req UpdateDepartmentRequest) (*models.Department, error)
	Delete(id uint) error
}

type departmentService struct {
	departmentRepo repositories_department.DepartmentRepository
}

func NewDepartmentService(
	departmentRepo repositories_department.DepartmentRepository,
) DepartmentService {
	return &departmentService{
		departmentRepo: departmentRepo,
	}
}

func (s *departmentService) Create(req CreateDepartmentRequest) (*models.Department, error) {
	name := strings.TrimSpace(req.DepartmentName)
	if name == "" {
		return nil, errors.New("department_name is required")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	department := &models.Department{
		DepartmentName: name,
		Description:    strings.TrimSpace(req.Description),
		IsActive:       isActive,
	}

	if err := s.departmentRepo.Create(department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *departmentService) GetByID(id uint) (*models.Department, error) {
	return s.departmentRepo.GetByID(id)
}

func (s *departmentService) GetAll() ([]models.Department, error) {
	return s.departmentRepo.GetAll()
}

func (s *departmentService) Update(id uint, req UpdateDepartmentRequest) (*models.Department, error) {
	department, err := s.departmentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.DepartmentName) != "" {
		department.DepartmentName = strings.TrimSpace(req.DepartmentName)
	}

	department.Description = strings.TrimSpace(req.Description)

	if req.IsActive != nil {
		department.IsActive = *req.IsActive
	}

	if err := s.departmentRepo.Update(department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *departmentService) Delete(id uint) error {
	return s.departmentRepo.Delete(id)
}