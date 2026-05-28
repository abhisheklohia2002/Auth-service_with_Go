package models

import "time"

type CertificateIssue struct {
	ID                   uint                `gorm:"primaryKey" json:"certificate_issue_id"`
	CertificationID      uint                `gorm:"not null;index" json:"certification_id"`
	UserID               uint                `gorm:"not null;index" json:"user_id"`
	TrainingAssignmentID *uint               `gorm:"index" json:"training_assignment_id"`
	IssuedOn             time.Time           `json:"issued_on"`
	ExpiryDate           time.Time           `json:"expiry_date"`
	CertificateNumber    string              `gorm:"size:100;not null;uniqueIndex" json:"certificate_number"`
	IssueStatus          string              `gorm:"size:50;default:'issued'" json:"issue_status"`
	PDFPath              string              `gorm:"size:500" json:"pdf_path"`
	Certification        Certification       `gorm:"foreignKey:CertificationID" json:"certification,omitempty"`
	User                 User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TrainingAssignment   *TrainingAssignment `gorm:"foreignKey:TrainingAssignmentID" json:"training_assignment,omitempty"`

	Notifications []Notification `gorm:"foreignKey:CertificateIssueID" json:"notifications,omitempty"`
}
