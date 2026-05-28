package dto

type CreateAssessmentRuleRequest struct {
	MaxAttempts      int    `json:"max_attempts" binding:"required,min=1"`
	PassingScore     int    `json:"passing_score" binding:"required,min=0"`
	RetakeAllowed    bool   `json:"retake_allowed"`
	EvaluationMethod string `json:"evaluation_method"`
}

type UpdateAssessmentRuleRequest struct {
	MaxAttempts      *int    `json:"max_attempts"`
	PassingScore     *int    `json:"passing_score"`
	RetakeAllowed    *bool   `json:"retake_allowed"`
	EvaluationMethod *string `json:"evaluation_method"`
}
