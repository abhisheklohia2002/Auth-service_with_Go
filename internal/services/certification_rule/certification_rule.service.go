package services_certificationrule

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_certificationrule "example.com/m/internal/repositories/certification_rule"
	repositories_course "example.com/m/internal/repositories/course"
)

type CertificationRuleService struct {
	certificationRuleRepo *repositories_certificationrule.CertificationRuleRepository
	courseRepo            *repositories_course.CourseRepository
}

func NewCertificationRuleService(repo *repositories_certificationrule.CertificationRuleRepository, courseRepo *repositories_course.CourseRepository) *CertificationRuleService {
	return &CertificationRuleService{certificationRuleRepo: repo, courseRepo: courseRepo}
}

func (s *CertificationRuleService) Create(req dto.CreateCertificationRuleRequest) (*models.CertificationRule, error) {
	if req.ValidityDays <= 0 {
		return nil, errors.New("validity days must be greater than zero")
	}
	if req.MinimumScoreRequired < 0 {
		return nil, errors.New("minimum score required cannot be negative")
	}

	rule := models.CertificationRule{
		IssueOnCourseCompletion: req.IssueOnCourseCompletion,
		MinimumScoreRequired:    req.MinimumScoreRequired,
		ValidityDays:            req.ValidityDays,
		RenewalRequired:         req.RenewalRequired,
	}

	return s.certificationRuleRepo.Create(&rule)
}

func (s *CertificationRuleService) FindAll() ([]models.CertificationRule, error) {
	return s.certificationRuleRepo.FindAll()
}

func (s *CertificationRuleService) FindByID(id uint) (*models.CertificationRule, error) {
	return s.certificationRuleRepo.FindByID(id)
}

func (s *CertificationRuleService) Update(id uint, req dto.UpdateCertificationRuleRequest) (*models.CertificationRule, error) {
	rule, err := s.certificationRuleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errors.New("certification rule not found")
	}

	if req.IssueOnCourseCompletion != nil {
		rule.IssueOnCourseCompletion = *req.IssueOnCourseCompletion
	}
	if req.MinimumScoreRequired != nil {
		if *req.MinimumScoreRequired < 0 {
			return nil, errors.New("minimum score required cannot be negative")
		}
		rule.MinimumScoreRequired = *req.MinimumScoreRequired
	}
	if req.ValidityDays != nil {
		if *req.ValidityDays <= 0 {
			return nil, errors.New("validity days must be greater than zero")
		}
		rule.ValidityDays = *req.ValidityDays
	}
	if req.RenewalRequired != nil {
		rule.RenewalRequired = *req.RenewalRequired
	}

	return s.certificationRuleRepo.Update(rule)
}

func (s *CertificationRuleService) Delete(id uint) error {
	rule, err := s.certificationRuleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if rule == nil {
		return errors.New("certification rule not found")
	}

	return s.certificationRuleRepo.Delete(id)
}

func (s *CertificationRuleService) CreateByCourseID(
	courseID uint,
	req dto.CreateCertificationRuleRequest,
) (*models.CertificationRule, error) {
	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	if req.MinimumScoreRequired < 0 {
		return nil, errors.New("minimum score required cannot be negative")
	}

	if req.ValidityDays <= 0 {
		return nil, errors.New("validity days must be greater than zero")
	}

	existingRule, err := s.certificationRuleRepo.FindByCourseID(courseID)
	if err != nil {
		return nil, err
	}

	if existingRule != nil {
		return nil, errors.New("certification rule already exists for this course")
	}

	rule := models.CertificationRule{
		CourseID:                courseID,
		IssueOnCourseCompletion: req.IssueOnCourseCompletion,
		MinimumScoreRequired:    req.MinimumScoreRequired,
		ValidityDays:            req.ValidityDays,
		RenewalRequired:         req.RenewalRequired,
	}

	return s.certificationRuleRepo.Create(&rule)
}
