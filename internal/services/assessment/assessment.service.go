package services_assessment

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentrule "example.com/m/internal/repositories/assessment_rule"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_module "example.com/m/internal/repositories/module"
)

type AssessmentService struct {
	assessmentRepo     *repositories_assessment.AssessmentRepository
	assessmentRuleRepo *repositories_assessmentrule.AssessmentRuleRepository
	courseRepo         *repositories_course.CourseRepository
	moduleRepo         *repositories_module.ModuleRepository
}

func NewAssessmentService(
	assessmentRepo *repositories_assessment.AssessmentRepository,
	assessmentRuleRepo *repositories_assessmentrule.AssessmentRuleRepository,
	courseRepo *repositories_course.CourseRepository,
	moduleRepo *repositories_module.ModuleRepository,
) *AssessmentService {
	return &AssessmentService{
		assessmentRepo:     assessmentRepo,
		assessmentRuleRepo: assessmentRuleRepo,
		courseRepo:         courseRepo,
		moduleRepo:         moduleRepo,
	}
}

func (s *AssessmentService) Create(req dto.CreateAssessmentRequest) (*models.Assessment, error) {
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	if req.ModuleID != nil {
		module, err := s.moduleRepo.FindByID(*req.ModuleID)
		if err != nil {
			return nil, err
		}
		if module == nil {
			return nil, errors.New("module not found")
		}
		if module.CourseID != req.CourseID {
			return nil, errors.New("module does not belong to this course")
		}
	}

	if req.RuleID != nil {
		rule, err := s.assessmentRuleRepo.FindByID(*req.RuleID)
		if err != nil {
			return nil, err
		}
		if rule == nil {
			return nil, errors.New("assessment rule not found")
		}
	}

	if req.PassingScore > req.MaxScore {
		return nil, errors.New("passing score cannot be greater than max score")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	assessment := models.Assessment{
		CourseID:        req.CourseID,
		ModuleID:        req.ModuleID,
		AssessmentTitle: req.AssessmentTitle,
		AssessmentType:  req.AssessmentType,
		MaxScore:        req.MaxScore,
		PassingScore:    req.PassingScore,
		RuleID:          req.RuleID,
		IsActive:        isActive,
	}

	return s.assessmentRepo.Create(&assessment)
}

func (s *AssessmentService) FindAll() ([]models.Assessment, error) {
	return s.assessmentRepo.FindAll()
}

func (s *AssessmentService) FindByID(id uint) (*models.Assessment, error) {
	return s.assessmentRepo.FindByID(id)
}

func (s *AssessmentService) FindByCourseID(courseID uint) ([]models.Assessment, error) {
	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	return s.assessmentRepo.FindByCourseID(courseID)
}

func (s *AssessmentService) FindByModuleID(moduleID uint) ([]models.Assessment, error) {
	module, err := s.moduleRepo.FindByID(moduleID)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, errors.New("module not found")
	}

	return s.assessmentRepo.FindByModuleID(moduleID)
}

func (s *AssessmentService) Update(id uint, req dto.UpdateAssessmentRequest) (*models.Assessment, error) {
	assessment, err := s.assessmentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	newCourseID := assessment.CourseID
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

	if req.ModuleID != nil {
		module, err := s.moduleRepo.FindByID(*req.ModuleID)
		if err != nil {
			return nil, err
		}
		if module == nil {
			return nil, errors.New("module not found")
		}
		if module.CourseID != newCourseID {
			return nil, errors.New("module does not belong to this course")
		}
		assessment.ModuleID = req.ModuleID
	}

	if req.RuleID != nil {
		rule, err := s.assessmentRuleRepo.FindByID(*req.RuleID)
		if err != nil {
			return nil, err
		}
		if rule == nil {
			return nil, errors.New("assessment rule not found")
		}
		assessment.RuleID = req.RuleID
	}

	assessment.CourseID = newCourseID

	if req.AssessmentTitle != nil {
		assessment.AssessmentTitle = *req.AssessmentTitle
	}

	if req.AssessmentType != nil {
		assessment.AssessmentType = *req.AssessmentType
	}

	if req.MaxScore != nil {
		if *req.MaxScore <= 0 {
			return nil, errors.New("max score must be greater than zero")
		}
		assessment.MaxScore = *req.MaxScore
	}

	if req.PassingScore != nil {
		if *req.PassingScore < 0 {
			return nil, errors.New("passing score cannot be negative")
		}
		assessment.PassingScore = *req.PassingScore
	}

	if assessment.PassingScore > assessment.MaxScore {
		return nil, errors.New("passing score cannot be greater than max score")
	}

	if req.IsActive != nil {
		assessment.IsActive = *req.IsActive
	}

	return s.assessmentRepo.Update(assessment)
}

func (s *AssessmentService) Delete(id uint) error {
	assessment, err := s.assessmentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if assessment == nil {
		return errors.New("assessment not found")
	}

	return s.assessmentRepo.Delete(id)
}
