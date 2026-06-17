package services_assignments

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_assignments "example.com/m/internal/repositories/assignments"
	repositories_course "example.com/m/internal/repositories/course"
)

type AssignmentRuleService struct {
	assignmentRuleRepo *repositories_assignments.AssignmentRuleRepository
	roleRepo           *repositories.RoleRepository
	courseRepo         *repositories_course.CourseRepository
}

func NewAssignmentRuleService(
	assignmentRuleRepo *repositories_assignments.AssignmentRuleRepository,
	roleRepo *repositories.RoleRepository,
	courseRepo *repositories_course.CourseRepository,
) *AssignmentRuleService {
	return &AssignmentRuleService{
		assignmentRuleRepo: assignmentRuleRepo,
		roleRepo:           roleRepo,
		courseRepo:         courseRepo,
	}
}

func (s *AssignmentRuleService) Create(req dto.CreateAssignmentRuleRequest) (*models.AssignmentRule, error) {
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, errors.New("role not found")
	}

	if req.DueDays <= 0 {
		return nil, errors.New("due days must be greater than 0")
	}

	rule := models.AssignmentRule{
		CourseID:     req.CourseID,
		RuleName:     req.RuleName,
		TriggerEvent: req.TriggerEvent,
		RoleID:       req.RoleID,
		DueDays:      req.DueDays,
		IsActive:     req.IsActive,
	}

	return s.assignmentRuleRepo.Create(&rule)
}

func (s *AssignmentRuleService) FindAll() ([]models.AssignmentRule, error) {
	return s.assignmentRuleRepo.FindAll()
}

func (s *AssignmentRuleService) FindByID(id uint) (*models.AssignmentRule, error) {
	return s.assignmentRuleRepo.FindByID(id)
}

func (s *AssignmentRuleService) Update(id uint, req dto.UpdateAssignmentRuleRequest) (*models.AssignmentRule, error) {
	rule, err := s.assignmentRuleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if rule == nil {
		return nil, errors.New("assignment rule not found")
	}

	if req.CourseID != nil {
		course, err := s.courseRepo.FindByID(*req.CourseID)
		if err != nil {
			return nil, err
		}

		if course == nil {
			return nil, errors.New("course not found")
		}

		rule.CourseID = *req.CourseID
	}

	if req.RuleName != nil {
		rule.RuleName = *req.RuleName
	}

	if req.TriggerEvent != nil {
		rule.TriggerEvent = *req.TriggerEvent
	}

	if req.RoleID != nil {
		role, err := s.roleRepo.FindByID(*req.RoleID)
		if err != nil {
			return nil, err
		}

		if role == nil {
			return nil, errors.New("role not found")
		}

		rule.RoleID = *req.RoleID
	}

	if req.DueDays != nil {
		if *req.DueDays <= 0 {
			return nil, errors.New("due days must be greater than 0")
		}

		rule.DueDays = *req.DueDays
	}

	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	return s.assignmentRuleRepo.Update(rule)
}

func (s *AssignmentRuleService) Delete(id uint) error {
	rule, err := s.assignmentRuleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if rule == nil {
		return errors.New("assignment rule not found")
	}

	return s.assignmentRuleRepo.Delete(id)
}

func (s *AssignmentRuleService) FindByCourseID(courseID uint) ([]models.AssignmentRule, error) {
	return s.assignmentRuleRepo.FindByCourseID(courseID)
}
