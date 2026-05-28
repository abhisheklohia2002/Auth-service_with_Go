package dto

type CreateAssessmentRequest struct {
	CourseID        uint  `json:"course_id" binding:"required"`
	ModuleID        *uint `json:"module_id"`
	AssessmentTitle string `json:"assessment_title" binding:"required"`
	AssessmentType  string `json:"assessment_type" binding:"required"`
	MaxScore         int    `json:"max_score" binding:"required,min=1"`
	PassingScore    int    `json:"passing_score" binding:"required,min=0"`
	RuleID          *uint  `json:"rule_id"`
	IsActive        *bool  `json:"is_active"`
}

type UpdateAssessmentRequest struct {
	CourseID        *uint   `json:"course_id"`
	ModuleID        *uint   `json:"module_id"`
	AssessmentTitle *string `json:"assessment_title"`
	AssessmentType  *string `json:"assessment_type"`
	MaxScore         *int    `json:"max_score"`
	PassingScore    *int    `json:"passing_score"`
	RuleID          *uint   `json:"rule_id"`
	IsActive        *bool   `json:"is_active"`
}