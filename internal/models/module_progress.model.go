package models

import "time"

type ModuleProgress struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AssignmentID uint      `gorm:"not null;uniqueIndex:idx_assignment_module" json:"assignment_id"`
	ModuleID     uint      `gorm:"not null;uniqueIndex:idx_assignment_module" json:"module_id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	Status       string    `gorm:"size:30;default:'pending'" json:"status"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`

	Assignment TrainingAssignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	Module     Module             `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	User       User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
}