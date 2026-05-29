package dto

type CreateAssignmentRuleRequest struct {
	RuleName     string `json:"rule_name" binding:"required"`
	TriggerEvent string `json:"trigger_event" binding:"required"`
	RoleID       uint   `json:"role_id" binding:"required"`
	IsActive     bool   `json:"is_active"`
}

type UpdateAssignmentRuleRequest struct {
	RuleName     *string `json:"rule_name"`
	TriggerEvent *string `json:"trigger_event"`
	RoleID       *uint   `json:"role_id"`
	IsActive     *bool   `json:"is_active"`
}