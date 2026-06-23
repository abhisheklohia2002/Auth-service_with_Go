package models

import "time"

type Course struct {
	ID                   uint      `gorm:"primaryKey" json:"course_id"`
	CourseTitle          string    `gorm:"size:200;not null" json:"course_title"`
	CourseDescription    string    `gorm:"type:text" json:"course_description"`
	CourseType           string    `gorm:"size:50;not null" json:"course_type"`
	IsActive             bool      `gorm:"default:true" json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	TotalDurationMinutes uint      `gorm:"not null;default:0" json:"total_duration_minutes"`
	CreatedByUserID      uint      `gorm:"not null" json:"created_by_user_id"`
	CreatedByUser        User      `gorm:"foreignKey:CreatedByUserID" json:"created_by_user,omitempty"`

	Modules             []Module             `gorm:"foreignKey:CourseID" json:"modules,omitempty"`
	Assessments         []Assessment         `gorm:"foreignKey:CourseID" json:"assessments,omitempty"`
	TrainingMappings    []TrainingMapping    `gorm:"foreignKey:CourseID" json:"training_mappings,omitempty"`
	TrainingAssignments []TrainingAssignment `gorm:"foreignKey:CourseID" json:"training_assignments,omitempty"`
	Certifications      []Certification      `gorm:"foreignKey:CourseID" json:"certifications,omitempty"`
}
