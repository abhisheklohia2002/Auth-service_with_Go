package services_trainingmapping

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_course "example.com/m/internal/repositories/course"
	repository_trainingmapping "example.com/m/internal/repositories/training_mapping"
)

type TrainingMappingService struct {
	trainingMappingRepo *repository_trainingmapping.TrainingMappingRepository
	roleRepo            *repositories.RoleRepository
	courseRepo          *repositories_course.CourseRepository
}

func NewTrainingMappingService(
	trainingMappingRepo *repository_trainingmapping.TrainingMappingRepository,
	roleRepo *repositories.RoleRepository,
	courseRepo *repositories_course.CourseRepository,
) *TrainingMappingService {
	return &TrainingMappingService{
		trainingMappingRepo: trainingMappingRepo,
		roleRepo:            roleRepo,
		courseRepo:          courseRepo,
	}
}

func (s *TrainingMappingService) Create(req dto.CreateTrainingMappingRequest) (*models.TrainingMapping, error) {
	role, err := s.roleRepo.FindByID(req.RoleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	exists, err := s.trainingMappingRepo.ExistsByRoleAndCourse(req.RoleID, req.CourseID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("training mapping already exists for this role and course")
	}

	activeFlag := true
	if req.ActiveFlag != nil {
		activeFlag = *req.ActiveFlag
	}

	assignmentTrigger := req.AssignmentTrigger
	if assignmentTrigger == "" {
		assignmentTrigger = "manual"
	}

	mapping := models.TrainingMapping{
		RoleID:            req.RoleID,
		CourseID:          req.CourseID,
		IsMandatory:       req.IsMandatory,
		AssignmentTrigger: assignmentTrigger,
		ActiveFlag:        activeFlag,
	}

	return s.trainingMappingRepo.Create(&mapping)
}

func (s *TrainingMappingService) FindAll() ([]models.TrainingMapping, error) {
	return s.trainingMappingRepo.FindAll()
}

func (s *TrainingMappingService) FindByID(id uint) (*models.TrainingMapping, error) {
	return s.trainingMappingRepo.FindByID(id)
}

func (s *TrainingMappingService) FindByRoleID(roleID uint) ([]models.TrainingMapping, error) {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	return s.trainingMappingRepo.FindByRoleID(roleID)
}

func (s *TrainingMappingService) Update(id uint, req dto.UpdateTrainingMappingRequest) (*models.TrainingMapping, error) {
	mapping, err := s.trainingMappingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if mapping == nil {
		return nil, errors.New("training mapping not found")
	}

	newRoleID := mapping.RoleID
	newCourseID := mapping.CourseID

	if req.RoleID != nil {
		role, err := s.roleRepo.FindByID(*req.RoleID)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, errors.New("role not found")
		}
		newRoleID = *req.RoleID
	}

	if req.CourseID != nil {
		course, err := s.courseRepo.FindByID(*req.CourseID)
		if err != nil {
			return nil, err
		}
		if course == nil {
			return nil, errors.New("course not found")
		}
		newCourseID = *req.CourseID
	}

	if newRoleID != mapping.RoleID || newCourseID != mapping.CourseID {
		exists, err := s.trainingMappingRepo.ExistsByRoleAndCourse(newRoleID, newCourseID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("training mapping already exists for this role and course")
		}
	}

	mapping.RoleID = newRoleID
	mapping.CourseID = newCourseID

	if req.IsMandatory != nil {
		mapping.IsMandatory = *req.IsMandatory
	}

	if req.AssignmentTrigger != nil {
		mapping.AssignmentTrigger = *req.AssignmentTrigger
	}

	if req.ActiveFlag != nil {
		mapping.ActiveFlag = *req.ActiveFlag
	}

	return s.trainingMappingRepo.Update(mapping)
}

func (s *TrainingMappingService) Delete(id uint) error {
	mapping, err := s.trainingMappingRepo.FindByID(id)
	if err != nil {
		return err
	}
	if mapping == nil {
		return errors.New("training mapping not found")
	}

	return s.trainingMappingRepo.Delete(id)
}
