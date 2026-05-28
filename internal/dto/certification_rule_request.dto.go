package dto

type CreateCertificationRequest struct {
	CourseID          uint   `json:"course_id" binding:"required"`
	CertificationName string `json:"certification_name" binding:"required"`
	ValidityDays      int    `json:"validity_days" binding:"required,min=1"`
	RuleID            *uint  `json:"rule_id"`
	IsActive          *bool  `json:"is_active"`
}
type CreateCertificationRuleRequest struct {
	IssueOnCourseCompletion bool `json:"issue_on_course_completion"`
	MinimumScoreRequired    int  `json:"minimum_score_required"`
	ValidityDays            int  `json:"validity_days" binding:"required,min=1"`
	RenewalRequired         bool `json:"renewal_required"`
}

type UpdateCertificationRequest struct {
	CourseID          *uint   `json:"course_id"`
	CertificationName *string `json:"certification_name"`
	ValidityDays      *int    `json:"validity_days"`
	RuleID            *uint   `json:"rule_id"`
	IsActive          *bool   `json:"is_active"`
}


type UpdateCertificationRuleRequest struct {
	IssueOnCourseCompletion *bool `json:"issue_on_course_completion"`
	MinimumScoreRequired    *int  `json:"minimum_score_required"`
	ValidityDays            *int  `json:"validity_days"`
	RenewalRequired          *bool `json:"renewal_required"`
}