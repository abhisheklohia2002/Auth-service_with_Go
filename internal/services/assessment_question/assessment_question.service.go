package services_assessmentquestion

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentquestion "example.com/m/internal/repositories/assessment_question"
	repositories_assessmentquestionoption "example.com/m/internal/repositories/assessment_question_option"
)

type AssessmentQuestionService struct {
	questionRepo   *repositories_assessmentquestion.AssessmentQuestionRepository
	optionRepo     *repositories_assessmentquestionoption.AssessmentQuestionOptionRepository
	assessmentRepo *repositories_assessment.AssessmentRepository
}

func NewAssessmentQuestionService(
	questionRepo *repositories_assessmentquestion.AssessmentQuestionRepository,
	optionRepo *repositories_assessmentquestionoption.AssessmentQuestionOptionRepository,
	assessmentRepo *repositories_assessment.AssessmentRepository,
) *AssessmentQuestionService {
	return &AssessmentQuestionService{
		questionRepo:   questionRepo,
		optionRepo:     optionRepo,
		assessmentRepo: assessmentRepo,
	}
}

func (s *AssessmentQuestionService) Create(req dto.CreateAssessmentQuestionRequest) (*models.AssessmentQuestion, error) {
	assessment, err := s.assessmentRepo.FindByID(req.AssessmentID)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	if !isValidQuestionType(req.QuestionType) {
		return nil, errors.New("invalid question type")
	}

	if err := validateQuestionOptions(req.QuestionType, req.Options); err != nil {
		return nil, err
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	question := models.AssessmentQuestion{
		AssessmentID: req.AssessmentID,
		QuestionText: req.QuestionText,
		QuestionType: req.QuestionType,
		Marks:        req.Marks,
		SequenceNo:   req.SequenceNo,
		IsActive:     isActive,
	}

	created, err := s.questionRepo.Create(&question)
	if err != nil {
		return nil, err
	}

	options := make([]models.AssessmentQuestionOption, 0)
	for _, opt := range req.Options {
		options = append(options, models.AssessmentQuestionOption{
			QuestionID: created.ID,
			OptionText: opt.OptionText,
			IsCorrect:  opt.IsCorrect,
		})
	}

	if err := s.optionRepo.CreateMany(options); err != nil {
		return nil, err
	}

	return s.questionRepo.FindByID(created.ID)
}

func (s *AssessmentQuestionService) FindByAssessmentID(assessmentID uint) ([]models.AssessmentQuestion, error) {
	assessment, err := s.assessmentRepo.FindByID(assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	return s.questionRepo.FindAllByAssessmentID(assessmentID)
}

func (s *AssessmentQuestionService) FindLearnerQuestions(assessmentID uint) ([]models.AssessmentQuestion, error) {
	assessment, err := s.assessmentRepo.FindByID(assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	questions, err := s.questionRepo.FindByAssessmentID(assessmentID)
	if err != nil {
		return nil, err
	}
	for qi := range questions {
		for oi := range questions[qi].Options {
			questions[qi].Options[oi].IsCorrect = false
		}
	}

	return questions, nil
}

func (s *AssessmentQuestionService) Delete(id uint) error {
	question, err := s.questionRepo.FindByID(id)
	if err != nil {
		return err
	}
	if question == nil {
		return errors.New("question not found")
	}

	if err := s.optionRepo.DeleteByQuestionID(id); err != nil {
		return err
	}

	return s.questionRepo.Delete(id)
}

func isValidQuestionType(questionType string) bool {
	switch questionType {
	case "single_choice", "multiple_choice", "true_false", "text":
		return true
	default:
		return false
	}
}

func validateQuestionOptions(questionType string, options []dto.CreateAssessmentQuestionOptionRequest) error {
	if questionType == "text" {
		return nil
	}

	if len(options) < 2 {
		return errors.New("question must have at least two options")
	}

	correctCount := 0
	for _, option := range options {
		if option.IsCorrect {
			correctCount++
		}
	}

	switch questionType {
	case "single_choice", "true_false":
		if correctCount != 1 {
			return errors.New("single choice and true/false questions must have exactly one correct option")
		}
	case "multiple_choice":
		if correctCount < 1 {
			return errors.New("multiple choice question must have at least one correct option")
		}
	}

	return nil
}


