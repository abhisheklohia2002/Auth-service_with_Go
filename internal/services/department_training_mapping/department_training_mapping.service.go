package services_department_training_mapping

import (
	"errors"
	"strings"

	"example.com/m/internal/models"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_department "example.com/m/internal/repositories/department"
	repositories_department_training_mapping "example.com/m/internal/repositories/department_training_mapping"
)

type CreateDepartmentTrainingMappingRequest struct {
	DepartmentID      uint   `json:"department_id" binding:"required"`
	CourseID          uint   `json:"course_id" binding:"required"`
	IsMandatory       bool   `json:"is_mandatory"`
	AssignmentTrigger string `json:"assignment_trigger"`
	ActiveFlag        bool   `json:"active_flag"`
	CreatedByUserID   uint   `json:"created_by_user_id"`
}

type DepartmentTrainingMappingService interface {
	Create(req CreateDepartmentTrainingMappingRequest) (*models.DepartmentTrainingMapping, error)
	GetAll() ([]models.DepartmentTrainingMapping, error)
	GetByDepartmentID(departmentID uint) ([]models.DepartmentTrainingMapping, error)
}

type departmentTrainingMappingService struct {
	mappingRepo    repositories_department_training_mapping.DepartmentTrainingMappingRepository
	departmentRepo repositories_department.DepartmentRepository
	courseRepo     repositories_course.CourseRepository
}

func NewDepartmentTrainingMappingService(
	mappingRepo repositories_department_training_mapping.DepartmentTrainingMappingRepository,
	departmentRepo repositories_department.DepartmentRepository,
	courseRepo repositories_course.CourseRepository,
) DepartmentTrainingMappingService {
	return &departmentTrainingMappingService{
		mappingRepo:    mappingRepo,
		departmentRepo: departmentRepo,
		courseRepo:     courseRepo,
	}
}

func (s *departmentTrainingMappingService) Create(
	req CreateDepartmentTrainingMappingRequest,
) (*models.DepartmentTrainingMapping, error) {
	if req.DepartmentID == 0 {
		return nil, errors.New("department_id is required")
	}

	if req.CourseID == 0 {
		return nil, errors.New("course_id is required")
	}

	department, err := s.departmentRepo.GetByID(req.DepartmentID)
	if err != nil {
		return nil, errors.New("department not found")
	}

	if !department.IsActive {
		return nil, errors.New("department is inactive")
	}

	_, err = s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, errors.New("course not found")
	}

	exists, err := s.mappingRepo.Exists(req.DepartmentID, req.CourseID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("department already mapped to this course")
	}

	activeFlag := req.ActiveFlag
	if !activeFlag {
		activeFlag = true
	}

	assignmentTrigger := strings.TrimSpace(req.AssignmentTrigger)
	if assignmentTrigger == "" {
		assignmentTrigger = "manual"
	}

	mapping := &models.DepartmentTrainingMapping{
		DepartmentID:      req.DepartmentID,
		CourseID:          req.CourseID,
		IsMandatory:       req.IsMandatory,
		AssignmentTrigger: assignmentTrigger,
		ActiveFlag:        activeFlag,
		CreatedByUserID:   req.CreatedByUserID,
	}

	if err := s.mappingRepo.Create(mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}

func (s *departmentTrainingMappingService) GetAll() ([]models.DepartmentTrainingMapping, error) {
	return s.mappingRepo.GetAll()
}

func (s *departmentTrainingMappingService) GetByDepartmentID(
	departmentID uint,
) ([]models.DepartmentTrainingMapping, error) {
	if departmentID == 0 {
		return nil, errors.New("department_id is required")
	}

	return s.mappingRepo.GetByDepartmentID(departmentID)
}
