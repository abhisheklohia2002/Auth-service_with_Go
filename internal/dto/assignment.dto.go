package dto

type CreateAssignmentRuleRequest struct {
	CourseID     uint   `json:"course_id" binding:"required"`
	RuleName     string `json:"rule_name" binding:"required"`
	TriggerEvent string `json:"trigger_event" binding:"required"`
	RoleID       uint   `json:"role_id" binding:"required"`
	DueDays      int    `json:"due_days" binding:"required"`
	IsActive     bool   `json:"is_active"`
}

type UpdateAssignmentRuleRequest struct {
	CourseID     *uint   `json:"course_id"`
	RuleName     *string `json:"rule_name"`
	TriggerEvent *string `json:"trigger_event"`
	RoleID       *uint   `json:"role_id"`
	DueDays      *int    `json:"due_days"`
	IsActive     *bool   `json:"is_active"`
}
