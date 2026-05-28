package models

type Module struct {
	ID                uint   `gorm:"primaryKey" json:"module_id"`
	CourseID          uint   `gorm:"not null;index" json:"course_id"`
	ModuleTitle       string `gorm:"size:200;not null" json:"module_title"`
	ModuleDescription string `gorm:"type:text" json:"module_description"`
	SequenceNo        int    `gorm:"not null" json:"sequence_no"`
	DueDays           int    `gorm:"default:0" json:"due_days"`
	IsActive          bool   `gorm:"default:true" json:"is_active"`

	Course      Course       `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Assessments []Assessment `gorm:"foreignKey:ModuleID" json:"assessments,omitempty"`
}