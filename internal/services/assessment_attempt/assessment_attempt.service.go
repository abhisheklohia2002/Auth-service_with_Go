package services_assessmentattempt

import (
	"encoding/json"
	"errors"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentattempt "example.com/m/internal/repositories/assessment_attempt"
	repositories_assessmentattemptanswer "example.com/m/internal/repositories/assessment_attempt_answer"
	repositories_assessmentquestion "example.com/m/internal/repositories/assessment_question"
	repositories_moduleprogress "example.com/m/internal/repositories/module_progress"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
)

type AssessmentAttemptService struct {
	assessmentAttemptRepo  *repositories_assessmentattempt.AssessmentAttemptRepository
	assessmentRepo         *repositories_assessment.AssessmentRepository
	userRepo               *repositories.UserRepository
	moduleProgressRepo     *repositories_moduleprogress.ModuleProgressRepository
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository
	questionRepo           *repositories_assessmentquestion.AssessmentQuestionRepository
	attemptAnswerRepo      *repositories_assessmentattemptanswer.AssessmentAttemptAnswerRepository
}

func NewAssessmentAttemptService(
	assessmentAttemptRepo *repositories_assessmentattempt.AssessmentAttemptRepository,
	assessmentRepo *repositories_assessment.AssessmentRepository,
	userRepo *repositories.UserRepository,
	moduleProgressRepo *repositories_moduleprogress.ModuleProgressRepository,
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository,
	attemptAnswerRepo *repositories_assessmentattemptanswer.AssessmentAttemptAnswerRepository,
) *AssessmentAttemptService {
	return &AssessmentAttemptService{
		assessmentAttemptRepo:  assessmentAttemptRepo,
		assessmentRepo:         assessmentRepo,
		userRepo:               userRepo,
		moduleProgressRepo:     moduleProgressRepo,
		trainingAssignmentRepo: trainingAssignmentRepo,
		attemptAnswerRepo:      attemptAnswerRepo,
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

	createdAttempt, err := s.assessmentAttemptRepo.Create(&attempt)
	if err != nil {
		return nil, err
	}

	if err := s.completeModuleProgressIfPassed(req.UserID, assessment, resultStatus); err != nil {
		return nil, err
	}

	return createdAttempt, nil
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

func (s *AssessmentAttemptService) completeModuleProgressIfPassed(
	userID uint,
	assessment *models.Assessment,
	resultStatus string,
) error {
	if resultStatus != "passed" {
		return nil
	}

	if assessment.ModuleID == nil {
		return nil
	}

	progress, err := s.moduleProgressRepo.FindByUserAndModule(userID, *assessment.ModuleID)
	if err != nil {
		return err
	}

	if progress == nil {
		return nil
	}

	now := time.Now()

	if progress.StartedAt == nil {
		progress.StartedAt = &now
	}

	progress.Status = "completed"
	progress.CompletedAt = &now

	_, err = s.moduleProgressRepo.Update(progress)
	if err != nil {
		return err
	}

	return s.updateAssignmentIfAllModulesCompleted(progress.AssignmentID)
}

func (s *AssessmentAttemptService) updateAssignmentIfAllModulesCompleted(assignmentID uint) error {
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

func (s *AssessmentAttemptService) Submit(req dto.SubmitAssessmentAttemptRequest) (*models.AssessmentAttempt, error) {
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

	questions, err := s.questionRepo.FindByAssessmentID(req.AssessmentID)
	if err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, errors.New("assessment has no questions")
	}

	scoreObtained, answerRows, err := calculateAssessmentScore(req.Answers, questions)
	if err != nil {
		return nil, err
	}

	resultStatus := "failed"
	if scoreObtained >= passingScore {
		resultStatus = "passed"
	}

	attempt := models.AssessmentAttempt{
		AssessmentID:  req.AssessmentID,
		UserID:        req.UserID,
		AttemptNo:     int(previousAttempts) + 1,
		ScoreObtained: scoreObtained,
		ResultStatus:  resultStatus,
		AttemptedAt:   time.Now(),
	}

	createdAttempt, err := s.assessmentAttemptRepo.Create(&attempt)
	if err != nil {
		return nil, err
	}

	for i := range answerRows {
		answerRows[i].AttemptID = createdAttempt.ID
	}

	if err := s.attemptAnswerRepo.CreateMany(answerRows); err != nil {
		return nil, err
	}

	if err := s.completeModuleProgressIfPassed(req.UserID, assessment, resultStatus); err != nil {
		return nil, err
	}

	return s.assessmentAttemptRepo.FindByID(createdAttempt.ID)
}

func calculateAssessmentScore(
	submittedAnswers []dto.SubmitAssessmentAnswerRequest,
	questions []models.AssessmentQuestion,
) (int, []models.AssessmentAttemptAnswer, error) {
	questionMap := make(map[uint]models.AssessmentQuestion)
	for _, question := range questions {
		questionMap[question.ID] = question
	}

	totalScore := 0
	answerRows := make([]models.AssessmentAttemptAnswer, 0)

	for _, submitted := range submittedAnswers {
		question, exists := questionMap[submitted.QuestionID]
		if !exists {
			return 0, nil, errors.New("submitted answer contains invalid question")
		}

		isCorrect := false
		marksAwarded := 0

		switch question.QuestionType {
		case "single_choice", "true_false":
			isCorrect = isSingleChoiceCorrect(question, submitted.SelectedOptionIDs)

		case "multiple_choice":
			isCorrect = isMultipleChoiceCorrect(question, submitted.SelectedOptionIDs)

		case "text":
			// For now text is manual review. Award 0.
			isCorrect = false
		}

		if isCorrect {
			marksAwarded = question.Marks
			totalScore += question.Marks
		}

		selectedJSON, err := json.Marshal(submitted.SelectedOptionIDs)
		if err != nil {
			return 0, nil, err
		}

		answerRows = append(answerRows, models.AssessmentAttemptAnswer{
			QuestionID:        submitted.QuestionID,
			SelectedOptionIDs: string(selectedJSON),
			TextAnswer:        submitted.TextAnswer,
			IsCorrect:         isCorrect,
			MarksAwarded:      marksAwarded,
		})
	}

	return totalScore, answerRows, nil
}

func isSingleChoiceCorrect(question models.AssessmentQuestion, selectedOptionIDs []uint) bool {
	if len(selectedOptionIDs) != 1 {
		return false
	}

	selected := selectedOptionIDs[0]

	for _, option := range question.Options {
		if option.ID == selected && option.IsCorrect {
			return true
		}
	}

	return false
}

func isMultipleChoiceCorrect(question models.AssessmentQuestion, selectedOptionIDs []uint) bool {
	correctMap := make(map[uint]bool)
	for _, option := range question.Options {
		if option.IsCorrect {
			correctMap[option.ID] = true
		}
	}

	if len(selectedOptionIDs) != len(correctMap) {
		return false
	}

	for _, selectedID := range selectedOptionIDs {
		if !correctMap[selectedID] {
			return false
		}
	}

	return true
}
