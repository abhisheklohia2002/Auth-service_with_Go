package models

import "time"

type TrainingSession struct {
	ID              uint  `gorm:"primaryKey" json:"session_id"`
	CourseID        uint  `gorm:"not null;index" json:"course_id"`
	ModuleID        *uint `gorm:"index" json:"module_id,omitempty"`
	CreatedByUserID uint  `gorm:"not null;index" json:"created_by_user_id"`

	SessionTitle string    `gorm:"size:255;not null" json:"session_title"`
	SessionType  string    `gorm:"size:50;not null" json:"session_type"` // online/offline/hybrid/self_paced
	StartTime    time.Time `gorm:"not null" json:"start_time"`
	EndTime      time.Time `gorm:"not null" json:"end_time"`

	Location    string `gorm:"size:255" json:"location"`
	MeetingLink string `gorm:"size:500" json:"meeting_link"`

	IsMandatory bool   `gorm:"default:true" json:"is_mandatory"`
	Status      string `gorm:"size:50;default:'scheduled'" json:"status"` // scheduled/completed/cancelled

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Course        Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Module        *Module `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	CreatedByUser User    `gorm:"foreignKey:CreatedByUserID" json:"created_by_user,omitempty"`

	Attendances []Attendance `gorm:"foreignKey:SessionID" json:"attendances,omitempty"`
}

type TrainingSessionFilter struct {
	CourseID    string
	ModuleID    string
	Status      string
	SessionType string
}
