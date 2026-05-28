package dto

type CreateAssessmentQuestionOptionRequest struct {
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

type CreateAssessmentQuestionRequest struct {
	AssessmentID uint                                    `json:"assessment_id" binding:"required"`
	QuestionText string                                  `json:"question_text" binding:"required"`
	QuestionType string                                  `json:"question_type" binding:"required"`
	Marks        int                                     `json:"marks" binding:"required,min=1"`
	SequenceNo   int                                     `json:"sequence_no" binding:"required,min=1"`
	IsActive     *bool                                   `json:"is_active"`
	Options      []CreateAssessmentQuestionOptionRequest `json:"options"`
}

type UpdateAssessmentQuestionRequest struct {
	QuestionText *string                                  `json:"question_text"`
	QuestionType *string                                  `json:"question_type"`
	Marks        *int                                     `json:"marks"`
	SequenceNo   *int                                     `json:"sequence_no"`
	IsActive     *bool                                    `json:"is_active"`
	Options      *[]CreateAssessmentQuestionOptionRequest `json:"options"`
}
