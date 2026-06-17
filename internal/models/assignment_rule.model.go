package models

import "time"

type AssignmentRule struct {
	ID           uint      `json:"assignment_rule_id" gorm:"primaryKey;column:assignment_rule_id"`
	CourseID     uint      `json:"course_id" gorm:"not null"`
	RuleName     string    `json:"rule_name" gorm:"not null"`
	TriggerEvent string    `json:"trigger_event" gorm:"not null"`
	RoleID       uint      `json:"role_id" gorm:"not null"`
	DueDays      int       `json:"due_days" gorm:"not null;default:14"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`

	Course Course `json:"course,omitempty" gorm:"foreignKey:CourseID;references:ID"`
	Role   Role   `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}