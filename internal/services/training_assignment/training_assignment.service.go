package services_trainingassignment

import (
	"errors"
	"log"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_module "example.com/m/internal/repositories/module"
	repositories_moduleprogress "example.com/m/internal/repositories/module_progress"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
	repository_trainingmapping "example.com/m/internal/repositories/training_mapping"
)

type TrainingAssignmentService struct {
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository
	trainingMappingRepo    *repository_trainingmapping.TrainingMappingRepository
	userRepo               *repositories.UserRepository
	courseRepo             *repositories_course.CourseRepository
	moduleRepo             *repositories_module.ModuleRepository
	moduleProgressRepo     *repositories_moduleprogress.ModuleProgressRepository
}

type CreateDepartmentTrainingAssignmentRequest struct {
	DepartmentID     uint       `json:"department_id" binding:"required"`
	CourseID         uint       `json:"course_id" binding:"required"`
	AssignedByUserID uint       `json:"assigned_by_user_id" binding:"required"`
	IsMandatory      bool       `json:"is_mandatory"`
	DueDate          *time.Time `json:"due_date"`
}
type DepartmentTrainingAssignmentResponse struct {
	DepartmentID         uint `json:"department_id"`
	CourseID             uint `json:"course_id"`
	TotalUsers           int  `json:"total_users"`
	AssignedCount        int  `json:"assigned_count"`
	SkippedExistingCount int  `json:"skipped_existing_count"`
}

func NewTrainingAssignmentService(
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository,
	trainingMappingRepo *repository_trainingmapping.TrainingMappingRepository,
	userRepo *repositories.UserRepository,
	courseRepo *repositories_course.CourseRepository,
	moduleRepo *repositories_module.ModuleRepository,
	moduleProgressRepo *repositories_moduleprogress.ModuleProgressRepository,
) *TrainingAssignmentService {
	return &TrainingAssignmentService{
		trainingAssignmentRepo: trainingAssignmentRepo,
		trainingMappingRepo:    trainingMappingRepo,
		userRepo:               userRepo,
		courseRepo:             courseRepo,
		moduleRepo:             moduleRepo,
		moduleProgressRepo:     moduleProgressRepo,
	}
}

func (s *TrainingAssignmentService) createModuleProgressForAssignment(assignment *models.TrainingAssignment) error {
	modules, err := s.moduleRepo.FindByCourseID(assignment.CourseID)
	if err != nil {
		return err
	}

	progresses := make([]models.ModuleProgress, 0)

	for _, module := range modules {
		log.Println("module found:", module.ID, module.ModuleTitle)
		exists, err := s.moduleProgressRepo.ExistsByAssignmentAndModule(assignment.ID, module.ID)
		if err != nil {
			return err
		}

		if exists {
			continue
		}

		progresses = append(progresses, models.ModuleProgress{
			AssignmentID: assignment.ID,
			ModuleID:     module.ID,
			UserID:       assignment.UserID,
			Status:       "pending",
		})
	}
	log.Println("module progress rows to create:", len(progresses))
	return s.moduleProgressRepo.BulkCreate(progresses)
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

	created, err := s.trainingAssignmentRepo.Create(&assignment)
	if err != nil {
		return nil, err
	}

	if err := s.createModuleProgressForAssignment(created); err != nil {
		return nil, err
	}

	return created, nil
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

		if err := s.createModuleProgressForAssignment(created); err != nil {
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



func (s *TrainingAssignmentService) AssignCourseToDepartment(
	req CreateDepartmentTrainingAssignmentRequest,
) (*DepartmentTrainingAssignmentResponse, error) {
	if req.DepartmentID == 0 {
		return nil, errors.New("department_id is required")
	}

	if req.CourseID == 0 {
		return nil, errors.New("course_id is required")
	}

	if req.AssignedByUserID == 0 {
		return nil, errors.New("assigned_by_user_id is required")
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

	users, err := s.userRepo.FindActiveByDepartmentID(req.DepartmentID)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("no active users found in this department")
	}

	modules, err := s.moduleRepo.FindByCourseID(req.CourseID)
	if err != nil {
		return nil, err
	}

	if len(modules) == 0 {
		return nil, errors.New("course has no modules")
	}

	assignedCount := 0
	skippedExistingCount := 0

	for _, user := range users {
		exists, err := s.trainingAssignmentRepo.ExistsByUserAndCourse(user.ID, req.CourseID)
		if err != nil {
			return nil, err
		}

		if exists {
			skippedExistingCount++
			continue
		}

		assignment := models.TrainingAssignment{
			UserID:           user.ID,
			CourseID:         req.CourseID,
			AssignedByUserID: req.AssignedByUserID,
			AssignmentSource: "department",
			IsMandatory:      req.IsMandatory,
			AssignedDate:     time.Now(),
			DueDate:          req.DueDate,
			Status:           "assigned",
			DepartmentID:     &req.DepartmentID,
		}

		created, err := s.trainingAssignmentRepo.Create(&assignment)
		if err != nil {
			return nil, err
		}

		if err := s.createModuleProgressForAssignment(created); err != nil {
			return nil, err
		}

		assignedCount++
	}

	return &DepartmentTrainingAssignmentResponse{
		DepartmentID:         req.DepartmentID,
		CourseID:             req.CourseID,
		TotalUsers:           len(users),
		AssignedCount:        assignedCount,
		SkippedExistingCount: skippedExistingCount,
	}, nil
}