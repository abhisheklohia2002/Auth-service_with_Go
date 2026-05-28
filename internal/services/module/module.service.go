package services_module

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_module "example.com/m/internal/repositories/module"
)

type ModuleService struct {
	moduleRepo *repositories_module.ModuleRepository
	courseRepo *repositories_course.CourseRepository
}

func NewModuleService(
	moduleRepo *repositories_module.ModuleRepository,
	courseRepo *repositories_course.CourseRepository,
) *ModuleService {
	return &ModuleService{
		moduleRepo: moduleRepo,
		courseRepo: courseRepo,
	}
}

func (s *ModuleService) Create(req dto.CreateModuleRequest) (*models.Module, error) {
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	module := models.Module{
		CourseID:          req.CourseID,
		ModuleTitle:       req.ModuleTitle,
		ModuleDescription: req.ModuleDescription,
		SequenceNo:        req.SequenceNo,
		DueDays:           req.DueDays,
		IsActive:          isActive,
	}

	return s.moduleRepo.Create(&module)
}

func (s *ModuleService) FindAll() ([]models.Module, error) {
	return s.moduleRepo.FindAll()
}

func (s *ModuleService) FindByID(id uint) (*models.Module, error) {
	return s.moduleRepo.FindByID(id)
}

func (s *ModuleService) FindByCourseID(courseID uint) ([]models.Module, error) {
	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	return s.moduleRepo.FindByCourseID(courseID)
}

func (s *ModuleService) Update(id uint, req dto.UpdateModuleRequest) (*models.Module, error) {
	module, err := s.moduleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, errors.New("module not found")
	}

	if req.CourseID != nil {
		course, err := s.courseRepo.FindByID(*req.CourseID)
		if err != nil {
			return nil, err
		}
		if course == nil {
			return nil, errors.New("course not found")
		}

		module.CourseID = *req.CourseID
	}

	if req.ModuleTitle != nil {
		module.ModuleTitle = *req.ModuleTitle
	}

	if req.ModuleDescription != nil {
		module.ModuleDescription = *req.ModuleDescription
	}

	if req.SequenceNo != nil {
		module.SequenceNo = *req.SequenceNo
	}

	if req.DueDays != nil {
		module.DueDays = *req.DueDays
	}

	if req.IsActive != nil {
		module.IsActive = *req.IsActive
	}

	return s.moduleRepo.Update(module)
}

func (s *ModuleService) Delete(id uint) error {
	module, err := s.moduleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if module == nil {
		return errors.New("module not found")
	}

	return s.moduleRepo.Delete(id)
}