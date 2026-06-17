package models

import "time"

type Notification struct {
	ID                 uint       `gorm:"primaryKey" json:"notification_id"`

	CreatedByID        uint       `gorm:"not null;index" json:"created_by_id"`

	AssignmentID       *uint      `gorm:"index" json:"assignment_id"`
	CertificateIssueID *uint      `gorm:"index" json:"certificate_issue_id"`

	NotificationType   string     `gorm:"size:100;not null;index" json:"notification_type"`
	Title              string     `gorm:"size:255" json:"title"`
	Message            string     `gorm:"type:text;not null" json:"message"`

	TargetAudience     string     `gorm:"size:100;index" json:"target_audience"`
	SentAt             *time.Time `json:"sent_at"`

	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	CreatedBy          User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`

	Assignment         *TrainingAssignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	CertificateIssue   *CertificateIssue   `gorm:"foreignKey:CertificateIssueID" json:"certificate_issue,omitempty"`

	Recipients         []NotificationRecipient `gorm:"foreignKey:NotificationID" json:"recipients,omitempty"`
}