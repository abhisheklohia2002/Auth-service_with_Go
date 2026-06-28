package services_moduleprogress

import (
	"errors"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_moduleDocument "example.com/m/internal/repositories/module_document"
	repositories_moduleprogress "example.com/m/internal/repositories/module_progress"

	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
)

const VideoRequiredPercent = 70

type ModuleProgressService struct {
	moduleProgressRepo     *repositories_moduleprogress.ModuleProgressRepository
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository
	moduleDocumentRepo     *repositories_moduleDocument.ModuleDocumentRepository
}

func NewModuleProgressService(
	moduleProgressRepo *repositories_moduleprogress.ModuleProgressRepository,
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository,
	moduleDocumentRepo *repositories_moduleDocument.ModuleDocumentRepository,
) *ModuleProgressService {
	return &ModuleProgressService{
		moduleProgressRepo:     moduleProgressRepo,
		trainingAssignmentRepo: trainingAssignmentRepo,
		moduleDocumentRepo:     moduleDocumentRepo,
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

func (s *ModuleProgressService) UpdateVideoProgress(
	progressID uint,
	req dto.UpdateVideoProgressRequest,
) (*models.ModuleProgress, error) {
	progress, err := s.moduleProgressRepo.FindByID(progressID)
	if err != nil {
		return nil, err
	}

	if progress == nil {
		return nil, errors.New("module progress not found")
	}

	if progress.Status == "completed" {
		return progress, nil
	}

	hasVideo, err := s.moduleDocumentRepo.ExistsByModuleID(progress.ModuleID)
	if err != nil {
		return nil, err
	}

	if !hasVideo {
		return nil, errors.New("this module does not have video")
	}

	if req.DurationSeconds <= 0 {
		return nil, errors.New("duration_seconds must be greater than zero")
	}

	if req.WatchedSeconds < 0 {
		return nil, errors.New("watched_seconds cannot be negative")
	}

	if req.WatchedSeconds > req.DurationSeconds {
		req.WatchedSeconds = req.DurationSeconds
	}

	watchedPercent := int(float64(req.WatchedSeconds) / float64(req.DurationSeconds) * 100)

	if watchedPercent > 100 {
		watchedPercent = 100
	}

	now := time.Now()

	if progress.Status == "pending" {
		progress.Status = "in_progress"
	}

	if progress.StartedAt == nil {
		progress.StartedAt = &now
	}

	progress.VideoDurationSeconds = req.DurationSeconds

	if req.WatchedSeconds > progress.VideoWatchedSeconds {
		progress.VideoWatchedSeconds = req.WatchedSeconds
	}

	if watchedPercent > progress.VideoWatchedPercent {
		progress.VideoWatchedPercent = watchedPercent
	}

	if progress.VideoWatchedPercent >= VideoRequiredPercent && progress.VideoCompletedAt == nil {
		progress.VideoCompletedAt = &now
	}

	return s.moduleProgressRepo.Update(progress)
}

func (s *ModuleProgressService) UpdateStatus(
	progressID uint,
	req dto.UpdateModuleProgressRequest,
) (*models.ModuleProgress, error) {
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

	if req.Status == "completed" {
		hasVideo, err := s.moduleDocumentRepo.ExistsByModuleID(progress.ModuleID)
		if err != nil {
			return nil, err
		}

		if hasVideo && progress.VideoWatchedPercent < VideoRequiredPercent {
			return nil, errors.New("watch at least 70% of the video before completing this module")
		}

		if progress.StartedAt == nil {
			progress.StartedAt = &now
		}

		progress.CompletedAt = &now
	}

	if req.Status == "in_progress" && progress.StartedAt == nil {
		progress.StartedAt = &now
	}

	progress.Status = req.Status

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
