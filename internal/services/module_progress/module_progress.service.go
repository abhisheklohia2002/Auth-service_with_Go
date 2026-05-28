package services_moduleprogress

import (
	"errors"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_moduleprogress "example.com/m/internal/repositories/module_progress"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
)

type ModuleProgressService struct {
	moduleProgressRepo     *repositories_moduleprogress.ModuleProgressRepository
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository
}

func NewModuleProgressService(
	moduleProgressRepo *repositories_moduleprogress.ModuleProgressRepository,
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository,
) *ModuleProgressService {
	return &ModuleProgressService{
		moduleProgressRepo:     moduleProgressRepo,
		trainingAssignmentRepo: trainingAssignmentRepo,
	}
}

func (s *ModuleProgressService) FindByAssignmentID(assignmentID uint) ([]models.ModuleProgress, error) {
	assignment, err := s.trainingAssignmentRepo.FindByID(assignmentID)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, errors.New("training assignment not found")
	}

	return s.moduleProgressRepo.FindByAssignmentID(assignmentID)
}

func (s *ModuleProgressService) UpdateStatus(progressID uint, req dto.UpdateModuleProgressRequest) (*models.ModuleProgress, error) {
	progress, err := s.moduleProgressRepo.FindByID(progressID)
	if err != nil {
		return nil, err
	}
	if progress == nil {
		return nil, errors.New("module progress not found")
	}

	if !isValidModuleProgressStatus(req.Status) {
		return nil, errors.New("invalid module progress status")
	}

	now := time.Now()

	progress.Status = req.Status

	if req.Status == "in_progress" && progress.StartedAt == nil {
		progress.StartedAt = &now
	}

	if req.Status == "completed" {
		if progress.StartedAt == nil {
			progress.StartedAt = &now
		}
		progress.CompletedAt = &now
	}

	updated, err := s.moduleProgressRepo.Update(progress)
	if err != nil {
		return nil, err
	}

	if req.Status == "completed" {
		if err := s.updateAssignmentIfAllModulesCompleted(progress.AssignmentID); err != nil {
			return nil, err
		}
	}

	return updated, nil
}

func (s *ModuleProgressService) updateAssignmentIfAllModulesCompleted(assignmentID uint) error {
	total, err := s.moduleProgressRepo.CountByAssignment(assignmentID)
	if err != nil {
		return err
	}

	if total == 0 {
		return nil
	}

	completed, err := s.moduleProgressRepo.CountByAssignmentAndStatus(assignmentID, "completed")
	if err != nil {
		return err
	}

	if completed != total {
		return nil
	}

	assignment, err := s.trainingAssignmentRepo.FindByID(assignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("training assignment not found")
	}

	now := time.Now()
	assignment.Status = "completed"
	assignment.CompletionDate = &now

	_, err = s.trainingAssignmentRepo.Update(assignment)
	return err
}

func isValidModuleProgressStatus(status string) bool {
	switch status {
	case "pending", "in_progress", "completed":
		return true
	default:
		return false
	}
}
