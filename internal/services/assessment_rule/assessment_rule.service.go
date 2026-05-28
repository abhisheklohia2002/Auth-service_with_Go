package services_assessmentrule

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_assessmentrule "example.com/m/internal/repositories/assessment_rule"
)

type AssessmentRuleService struct {
	assessmentRuleRepo *repositories_assessmentrule.AssessmentRuleRepository
}

func NewAssessmentRuleService(assessmentRuleRepo *repositories_assessmentrule.AssessmentRuleRepository) *AssessmentRuleService {
	return &AssessmentRuleService{
		assessmentRuleRepo: assessmentRuleRepo,
	}
}

func (s *AssessmentRuleService) Create(req dto.CreateAssessmentRuleRequest) (*models.AssessmentRule, error) {
	if req.PassingScore < 0 {
		return nil, errors.New("passing score cannot be negative")
	}

	if req.MaxAttempts <= 0 {
		return nil, errors.New("max attempts must be greater than zero")
	}

	evaluationMethod := req.EvaluationMethod
	if evaluationMethod == "" {
		evaluationMethod = "score"
	}

	rule := models.AssessmentRule{
		MaxAttempts:      req.MaxAttempts,
		PassingScore:     req.PassingScore,
		RetakeAllowed:    req.RetakeAllowed,
		EvaluationMethod: evaluationMethod,
	}

	return s.assessmentRuleRepo.Create(&rule)
}

func (s *AssessmentRuleService) FindAll() ([]models.AssessmentRule, error) {
	return s.assessmentRuleRepo.FindAll()
}

func (s *AssessmentRuleService) FindByID(id uint) (*models.AssessmentRule, error) {
	return s.assessmentRuleRepo.FindByID(id)
}

func (s *AssessmentRuleService) Update(id uint, req dto.UpdateAssessmentRuleRequest) (*models.AssessmentRule, error) {
	rule, err := s.assessmentRuleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errors.New("assessment rule not found")
	}

	if req.MaxAttempts != nil {
		if *req.MaxAttempts <= 0 {
			return nil, errors.New("max attempts must be greater than zero")
		}
		rule.MaxAttempts = *req.MaxAttempts
	}

	if req.PassingScore != nil {
		if *req.PassingScore < 0 {
			return nil, errors.New("passing score cannot be negative")
		}
		rule.PassingScore = *req.PassingScore
	}

	if req.RetakeAllowed != nil {
		rule.RetakeAllowed = *req.RetakeAllowed
	}

	if req.EvaluationMethod != nil {
		rule.EvaluationMethod = *req.EvaluationMethod
	}

	return s.assessmentRuleRepo.Update(rule)
}

func (s *AssessmentRuleService) Delete(id uint) error {
	rule, err := s.assessmentRuleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if rule == nil {
		return errors.New("assessment rule not found")
	}

	return s.assessmentRuleRepo.Delete(id)
}
