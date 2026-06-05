package models

import "time"

type TrainingAssignment struct {
	ID                uint       `gorm:"primaryKey" json:"assignment_id"`
	UserID            uint       `gorm:"not null;index" json:"user_id"`
	CourseID          uint       `gorm:"not null;index" json:"course_id"`
	AssignedByUserID  uint       `gorm:"not null;index" json:"assigned_by_user_id"`
	AssignmentSource  string     `gorm:"size:100" json:"assignment_source"`
	IsMandatory       bool       `gorm:"default:false" json:"is_mandatory"`
	AssignedDate      time.Time  `json:"assigned_date"`
	DueDate           *time.Time `json:"due_date"`
	CompletionDate    *time.Time `json:"completion_date"`
	Status            string     `gorm:"size:50;default:'assigned'" json:"status"`
	ImprovementStatus string     `gorm:"size:50" json:"improvement_status"`

	User           User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course         Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	AssignedByUser User   `gorm:"foreignKey:AssignedByUserID" json:"assigned_by_user,omitempty"`

	Notifications     []Notification     `gorm:"foreignKey:AssignmentID" json:"notifications,omitempty"`
	CertificateIssues []CertificateIssue `gorm:"foreignKey:TrainingAssignmentID" json:"certificate_issues,omitempty"`
	ModuleProgresses  []ModuleProgress   `gorm:"foreignKey:AssignmentID" json:"module_progresses,omitempty"`

	DepartmentID *uint       `json:"department_id"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}
