package models

import "time"

type Attendance struct {
	ID             uint  `gorm:"primaryKey" json:"attendance_id"`
	SessionID      uint  `gorm:"not null;index;uniqueIndex:idx_session_user" json:"session_id"`
	UserID         uint  `gorm:"not null;index;uniqueIndex:idx_session_user" json:"user_id"`
	MarkedByUserID *uint `gorm:"index" json:"marked_by_user_id,omitempty"`

	Status string `gorm:"size:50;not null" json:"status"`

	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`

	DurationMinutes  int    `gorm:"default:0" json:"duration_minutes"`
	AttendanceSource string `gorm:"size:50;default:'manual'" json:"attendance_source"`

	Remarks string `gorm:"size:500" json:"remarks"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
