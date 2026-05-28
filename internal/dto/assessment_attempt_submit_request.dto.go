package dto

type SubmitAssessmentAnswerRequest struct {
	QuestionID        uint   `json:"question_id" binding:"required"`
	SelectedOptionIDs []uint `json:"selected_option_ids"`
	TextAnswer        string `json:"text_answer"`
}

type SubmitAssessmentAttemptRequest struct {
	AssessmentID uint                            `json:"assessment_id" binding:"required"`
	UserID       uint                            `json:"user_id" binding:"required"`
	Answers      []SubmitAssessmentAnswerRequest `json:"answers" binding:"required"`
}
