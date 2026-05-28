package services_trainingassignment

import (
	"errors"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
	repository_trainingmapping "example.com/m/internal/repositories/training_mapping"
)

type TrainingAssignmentService struct {
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository
	trainingMappingRepo    *repository_trainingmapping.TrainingMappingRepository
	userRepo               *repositories.UserRepository
	courseRepo             *repositories_course.CourseRepository
}

func NewTrainingAssignmentService(
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository,
	trainingMappingRepo *repository_trainingmapping.TrainingMappingRepository,
	userRepo *repositories.UserRepository,
	courseRepo *repositories_course.CourseRepository,
) *TrainingAssignmentService {
	return &TrainingAssignmentService{
		trainingAssignmentRepo: trainingAssignmentRepo,
		trainingMappingRepo:    trainingMappingRepo,
		userRepo:               userRepo,
		courseRepo:             courseRepo,
	}
}

func (s *TrainingAssignmentService) CreateManual(req dto.CreateTrainingAssignmentRequest) (*models.TrainingAssignment, error) {
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	assignedByUser, err := s.userRepo.FindByID(req.AssignedByUserID)
	if err != nil {
		return nil, err
	}
	if assignedByUser == nil {
		return nil, errors.New("assigned by user not found")
	}

	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	exists, err := s.trainingAssignmentRepo.ExistsByUserAndCourse(req.UserID, req.CourseID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("course already assigned to this user")
	}

	source := req.AssignmentSource
	if source == "" {
		source = "manual"
	}

	assignment := models.TrainingAssignment{
		UserID:           req.UserID,
		CourseID:         req.CourseID,
		AssignedByUserID: req.AssignedByUserID,
		AssignmentSource: source,
		IsMandatory:      req.IsMandatory,
		AssignedDate:     time.Now(),
		DueDate:          req.DueDate,
		Status:           "assigned",
	}

	return s.trainingAssignmentRepo.Create(&assignment)
}

func (s *TrainingAssignmentService) AutoAssignByUserRole(req dto.AutoAssignTrainingRequest) ([]models.TrainingAssignment, error) {
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	assignedByUser, err := s.userRepo.FindByID(req.AssignedByUserID)
	if err != nil {
		return nil, err
	}
	if assignedByUser == nil {
		return nil, errors.New("assigned by user not found")
	}

	mappings, err := s.trainingMappingRepo.FindActiveByRoleID(user.RoleID)
	if err != nil {
		return nil, err
	}

	if len(mappings) == 0 {
		return []models.TrainingAssignment{}, nil
	}

	createdAssignments := make([]models.TrainingAssignment, 0)

	for _, mapping := range mappings {
		exists, err := s.trainingAssignmentRepo.ExistsByUserAndCourse(user.ID, mapping.CourseID)
		if err != nil {
			return nil, err
		}

		if exists {
			continue
		}

		assignment := models.TrainingAssignment{
			UserID:           user.ID,
			CourseID:         mapping.CourseID,
			AssignedByUserID: req.AssignedByUserID,
			AssignmentSource: "auto_role_mapping",
			IsMandatory:      mapping.IsMandatory,
			AssignedDate:     time.Now(),
			Status:           "assigned",
		}

		created, err := s.trainingAssignmentRepo.Create(&assignment)
		if err != nil {
			return nil, err
		}

		createdAssignments = append(createdAssignments, *created)
	}

	return createdAssignments, nil
}

func (s *TrainingAssignmentService) FindAll() ([]models.TrainingAssignment, error) {
	return s.trainingAssignmentRepo.FindAll()
}

func (s *TrainingAssignmentService) FindByID(id uint) (*models.TrainingAssignment, error) {
	return s.trainingAssignmentRepo.FindByID(id)
}

func (s *TrainingAssignmentService) FindByUserID(userID uint) ([]models.TrainingAssignment, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return s.trainingAssignmentRepo.FindByUserID(userID)
}

func (s *TrainingAssignmentService) UpdateStatus(id uint, req dto.UpdateTrainingAssignmentStatusRequest) (*models.TrainingAssignment, error) {
	assignment, err := s.trainingAssignmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, errors.New("training assignment not found")
	}

	if !isValidAssignmentStatus(req.Status) {
		return nil, errors.New("invalid assignment status")
	}

	assignment.Status = req.Status

	if req.Status == "completed" {
		now := time.Now()
		if req.CompletionDate != nil {
			assignment.CompletionDate = req.CompletionDate
		} else {
			assignment.CompletionDate = &now
		}
	}

	if req.ImprovementStatus != "" {
		assignment.ImprovementStatus = req.ImprovementStatus
	}

	return s.trainingAssignmentRepo.Update(assignment)
}

func (s *TrainingAssignmentService) Delete(id uint) error {
	assignment, err := s.trainingAssignmentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("training assignment not found")
	}

	return s.trainingAssignmentRepo.Delete(id)
}

func isValidAssignmentStatus(status string) bool {
	switch status {
	case "assigned", "in_progress", "completed", "overdue", "cancelled":
		return true
	default:
		return false
	}
}
