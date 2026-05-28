package dto

type CreateTrainingMappingRequest struct {
	RoleID            uint   `json:"role_id" binding:"required"`
	CourseID          uint   `json:"course_id" binding:"required"`
	IsMandatory       bool   `json:"is_mandatory"`
	AssignmentTrigger string `json:"assignment_trigger"`
	ActiveFlag        *bool  `json:"active_flag"`
}

type UpdateTrainingMappingRequest struct {
	RoleID            *uint   `json:"role_id"`
	CourseID          *uint   `json:"course_id"`
	IsMandatory       *bool   `json:"is_mandatory"`
	AssignmentTrigger *string `json:"assignment_trigger"`
	ActiveFlag        *bool   `json:"active_flag"`
}
