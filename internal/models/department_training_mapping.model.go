package models

import "time"

type DepartmentTrainingMapping struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	DepartmentID      uint      `gorm:"not null" json:"department_id"`
	CourseID          uint      `gorm:"not null" json:"course_id"`
	IsMandatory       bool      `gorm:"default:true" json:"is_mandatory"`
	AssignmentTrigger string    `json:"assignment_trigger"`
	ActiveFlag        bool      `gorm:"default:true" json:"active_flag"`
	CreatedByUserID   uint      `json:"created_by_user_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Department Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Course     Course     `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	CreatedBy  User       `gorm:"foreignKey:CreatedByUserID" json:"created_by,omitempty"`
}
