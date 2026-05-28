package dto

type CreateModuleRequest struct {
	CourseID          uint   `json:"course_id" binding:"required"`
	ModuleTitle       string `json:"module_title" binding:"required"`
	ModuleDescription string `json:"module_description"`
	SequenceNo        int    `json:"sequence_no" binding:"required"`
	DueDays           int    `json:"due_days"`
	IsActive          *bool  `json:"is_active"`
}

type UpdateModuleRequest struct {
	CourseID          *uint   `json:"course_id"`
	ModuleTitle       *string `json:"module_title"`
	ModuleDescription *string `json:"module_description"`
	SequenceNo        *int    `json:"sequence_no"`
	DueDays           *int    `json:"due_days"`
	IsActive          *bool   `json:"is_active"`
}