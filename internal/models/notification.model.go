package models

import "time"

type Notification struct {
	ID                 uint       `gorm:"primaryKey" json:"notification_id"`
	UserID             uint       `gorm:"not null;index" json:"user_id"`
	AssignmentID       *uint      `gorm:"index" json:"assignment_id"`
	CertificateIssueID *uint      `gorm:"index" json:"certificate_issue_id"`
	NotificationType   string     `gorm:"size:100;not null" json:"notification_type"`
	Message            string     `gorm:"type:text;not null" json:"message"`
	TargetAudience     string     `gorm:"size:100" json:"target_audience"`
	SentAt             *time.Time `json:"sent_at"`
	ReadStatus         bool       `gorm:"default:false" json:"read_status"`

	User             User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Assignment       *TrainingAssignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	CertificateIssue *CertificateIssue   `gorm:"foreignKey:CertificateIssueID" json:"certificate_issue,omitempty"`
}