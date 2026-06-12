package services_assessmentquestion

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	"example.com/m/internal/models"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentquestion "example.com/m/internal/repositories/assessment_question"
	repositories_assessmentquestionoption "example.com/m/internal/repositories/assessment_question_option"
	"github.com/xuri/excelize/v2"
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



func (s *AssessmentQuestionService) BulkUploadQuestions(
	ctx context.Context,
	assessmentID uint,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (*dto.BulkQuestionUploadResponse, error) {
	if file == nil || fileHeader == nil {
		return nil, errors.New("file is required")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".xlsx" {
		return nil, errors.New("only .xlsx files are allowed")
	}

	if fileHeader.Size > 10*1024*1024 {
		return nil, errors.New("file size must be less than 10MB")
	}

	assessment, err := s.assessmentRepo.FindByID(assessmentID)
	if err != nil {
		return nil, err
	}

	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("invalid excel file")
	}
	defer xlsx.Close()

	rows, err := xlsx.GetRows("Questions")
	if err != nil {
		return nil, errors.New("sheet 'Questions' not found")
	}

	if len(rows) < 2 {
		return nil, errors.New("excel file has no question rows")
	}

	result := &dto.BulkQuestionUploadResponse{
		TotalRows: len(rows) - 1,
		Errors:    []dto.BulkQuestionError{},
	}

	var questions []dto.CreateAssessmentQuestionRequest

	for i := 1; i < len(rows); i++ {
		rowNumber := i + 1

		question, rowErrors := parseAssessmentQuestionRow(rows[i], rowNumber, assessmentID)
		if len(rowErrors) > 0 {
			result.Errors = append(result.Errors, rowErrors...)
			continue
		}

		questions = append(questions, question)
	}

	result.ValidRows = len(questions)
	result.InvalidRows = result.TotalRows - result.ValidRows

	if result.InvalidRows > 0 {
		result.Success = false
		return result, nil
	}

	if err := s.questionRepo.CreateBulk(ctx, questions); err != nil {
		return nil, err
	}

	result.Success = true
	return result, nil
}



func parseAssessmentQuestionRow(
	row []string,
	rowNumber int,
	assessmentID uint,
) (dto.CreateAssessmentQuestionRequest, []dto.BulkQuestionError) {
	var rowErrors []dto.BulkQuestionError

	get := func(index int) string {
		if index >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[index])
	}

	sequenceNoText := get(0)
	questionText := get(1)
	questionType := strings.ToLower(get(2))
	marksText := get(3)
	isActiveText := strings.ToLower(get(4))
	optionA := get(5)
	optionB := get(6)
	optionC := get(7)
	optionD := get(8)
	correctAnswer := strings.ToUpper(strings.ReplaceAll(get(9), " ", ""))

	if questionText == "" {
		rowErrors = append(rowErrors, dto.BulkQuestionError{
			Row:     rowNumber,
			Field:   "question_text",
			Message: "question_text is required",
		})
	}

	if !isValidQuestionType(questionType) {
		rowErrors = append(rowErrors, dto.BulkQuestionError{
			Row:     rowNumber,
			Field:   "question_type",
			Message: "question_type must be single_choice, multiple_choice, or true_false",
		})
	}

	marks, err := strconv.Atoi(marksText)
	if err != nil || marks <= 0 {
		rowErrors = append(rowErrors, dto.BulkQuestionError{
			Row:     rowNumber,
			Field:   "marks",
			Message: "marks must be a positive number",
		})
	}

	sequenceNo, err := strconv.Atoi(sequenceNoText)
	if err != nil || sequenceNo <= 0 {
		rowErrors = append(rowErrors, dto.BulkQuestionError{
			Row:     rowNumber,
			Field:   "sequence_no",
			Message: "sequence_no must be a positive number",
		})
	}

	isActive := true
	if isActiveText != "" {
		parsedActive, err := strconv.ParseBool(isActiveText)
		if err != nil {
			rowErrors = append(rowErrors, dto.BulkQuestionError{
				Row:     rowNumber,
				Field:   "is_active",
				Message: "is_active must be true or false",
			})
		} else {
			isActive = parsedActive
		}
	}

	options := []dto.CreateAssessmentQuestionOptionRequest{}

	switch questionType {
	case "single_choice", "multiple_choice":
		optionMap := map[string]string{
			"A": optionA,
			"B": optionB,
			"C": optionC,
			"D": optionD,
		}

		for label, text := range optionMap {
			if text != "" {
				options = append(options, dto.CreateAssessmentQuestionOptionRequest{
					OptionText: text,
					IsCorrect:  helper.IsCorrectLabel(label, correctAnswer),
				})
			}
		}

		if len(options) < 2 {
			rowErrors = append(rowErrors, dto.BulkQuestionError{
				Row:     rowNumber,
				Field:   "options",
				Message: "at least 2 options are required",
			})
		}

		if correctAnswer == "" {
			rowErrors = append(rowErrors, dto.BulkQuestionError{
				Row:     rowNumber,
				Field:   "correct_answer",
				Message: "correct_answer is required",
			})
		}

		if questionType == "single_choice" && strings.Contains(correctAnswer, ",") {
			rowErrors = append(rowErrors, dto.BulkQuestionError{
				Row:     rowNumber,
				Field:   "correct_answer",
				Message: "single_choice allows only one correct answer",
			})
		}

	case "true_false":
		if optionA == "" {
			optionA = "true"
		}
		if optionB == "" {
			optionB = "false"
		}

		options = []dto.CreateAssessmentQuestionOptionRequest{
			{
				OptionText: optionA,
				IsCorrect:  helper.IsCorrectLabel("A", correctAnswer),
			},
			{
				OptionText: optionB,
				IsCorrect:  helper.IsCorrectLabel("B", correctAnswer),
			},
		}

		if correctAnswer != "A" && correctAnswer != "B" {
			rowErrors = append(rowErrors, dto.BulkQuestionError{
				Row:     rowNumber,
				Field:   "correct_answer",
				Message: "true_false correct_answer must be A or B",
			})
		}
	}

	if len(rowErrors) > 0 {
		return dto.CreateAssessmentQuestionRequest{}, rowErrors
	}

	return dto.CreateAssessmentQuestionRequest{
		AssessmentID: assessmentID,
		QuestionText: questionText,
		QuestionType: questionType,
		Marks:        marks,
		SequenceNo:   sequenceNo,
		IsActive:     &isActive,
		Options:      options,
	}, nil
}