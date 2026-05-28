package dto

import "time"

type CreateTrainingAssignmentRequest struct {
	UserID           uint       `json:"user_id" binding:"required"`
	CourseID         uint       `json:"course_id" binding:"required"`
	AssignedByUserID uint       `json:"assigned_by_user_id" binding:"required"`
	AssignmentSource string     `json:"assignment_source"`
	IsMandatory      bool       `json:"is_mandatory"`
	DueDate          *time.Time `json:"due_date"`
}

type AutoAssignTrainingRequest struct {
	UserID           uint `json:"user_id" binding:"required"`
	AssignedByUserID uint `json:"assigned_by_user_id" binding:"required"`
}

type UpdateTrainingAssignmentStatusRequest struct {
	Status            string     `json:"status" binding:"required"`
	CompletionDate    *time.Time `json:"completion_date"`
	ImprovementStatus string     `json:"improvement_status"`
}
