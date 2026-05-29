package services_assignments

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_assignments "example.com/m/internal/repositories/assignments"
)

type AssignmentRuleService struct {
	assignmentRuleRepo *repositories_assignments.AssignmentRuleRepository
	roleRepo           *repositories.RoleRepository
}

func NewAssignmentRuleService(
	assignmentRuleRepo *repositories_assignments.AssignmentRuleRepository,
	roleRepo *repositories.RoleRepository,
) *AssignmentRuleService {
	return &AssignmentRuleService{
		assignmentRuleRepo: assignmentRuleRepo,
		roleRepo:           roleRepo,
	}
}

func (s *AssignmentRuleService) Create(req dto.CreateAssignmentRuleRequest) (*models.AssignmentRule, error) {
	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	rule := models.AssignmentRule{
		RuleName:     req.RuleName,
		TriggerEvent: req.TriggerEvent,
		RoleID:       req.RoleID,
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
