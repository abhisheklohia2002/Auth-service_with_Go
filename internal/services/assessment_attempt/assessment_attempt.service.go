package services_assessmentattempt

import (
	"errors"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentattempt "example.com/m/internal/repositories/assessment_attempt"
)

type AssessmentAttemptService struct {
	assessmentAttemptRepo *repositories_assessmentattempt.AssessmentAttemptRepository
	assessmentRepo        *repositories_assessment.AssessmentRepository
	userRepo              *repositories.UserRepository
}

func NewAssessmentAttemptService(
	assessmentAttemptRepo *repositories_assessmentattempt.AssessmentAttemptRepository,
	assessmentRepo *repositories_assessment.AssessmentRepository,
	userRepo *repositories.UserRepository,
) *AssessmentAttemptService {
	return &AssessmentAttemptService{
		assessmentAttemptRepo: assessmentAttemptRepo,
		assessmentRepo:        assessmentRepo,
		userRepo:              userRepo,
	}
}

func (s *AssessmentAttemptService) Create(req dto.CreateAssessmentAttemptRequest) (*models.AssessmentAttempt, error) {
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	assessment, err := s.assessmentRepo.FindByID(req.AssessmentID)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	if !assessment.IsActive {
		return nil, errors.New("assessment is inactive")
	}

	if req.ScoreObtained > assessment.MaxScore {
		return nil, errors.New("score obtained cannot be greater than max score")
	}

	hasPassed, err := s.assessmentAttemptRepo.HasPassed(req.UserID, req.AssessmentID)
	if err != nil {
		return nil, err
	}
	if hasPassed {
		return nil, errors.New("user has already passed this assessment")
	}

	previousAttempts, err := s.assessmentAttemptRepo.CountByUserAndAssessment(req.UserID, req.AssessmentID)
	if err != nil {
		return nil, err
	}

	maxAttempts := 1
	retakeAllowed := false
	passingScore := assessment.PassingScore

	if assessment.Rule != nil {
		maxAttempts = assessment.Rule.MaxAttempts
		retakeAllowed = assessment.Rule.RetakeAllowed

		if assessment.Rule.PassingScore > 0 {
			passingScore = assessment.Rule.PassingScore
		}
	}

	if previousAttempts > 0 && !retakeAllowed {
		return nil, errors.New("retake is not allowed for this assessment")
	}

	if previousAttempts >= int64(maxAttempts) {
		return nil, errors.New("maximum attempts reached")
	}

	resultStatus := "failed"
	if req.ScoreObtained >= passingScore {
		resultStatus = "passed"
	}

	attempt := models.AssessmentAttempt{
		AssessmentID:  req.AssessmentID,
		UserID:        req.UserID,
		AttemptNo:     int(previousAttempts) + 1,
		ScoreObtained: req.ScoreObtained,
		ResultStatus:  resultStatus,
		AttemptedAt:   time.Now(),
	}

	return s.assessmentAttemptRepo.Create(&attempt)
}

func (s *AssessmentAttemptService) FindByUserID(userID uint) ([]models.AssessmentAttempt, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return s.assessmentAttemptRepo.FindByUserID(userID)
}

func (s *AssessmentAttemptService) FindByAssessmentID(assessmentID uint) ([]models.AssessmentAttempt, error) {
	assessment, err := s.assessmentRepo.FindByID(assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	return s.assessmentAttemptRepo.FindByAssessmentID(assessmentID)
}

func (s *AssessmentAttemptService) FindByUserAndAssessment(userID uint, assessmentID uint) ([]models.AssessmentAttempt, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	assessment, err := s.assessmentRepo.FindByID(assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	return s.assessmentAttemptRepo.FindByUserAndAssessment(userID, assessmentID)
}