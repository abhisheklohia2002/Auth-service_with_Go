package models

import "time"

type AssessmentAttempt struct {
	ID            uint      `gorm:"primaryKey" json:"attempt_id"`
	AssessmentID  uint      `gorm:"not null;index" json:"assessment_id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	AttemptNo     int       `gorm:"not null" json:"attempt_no"`
	ScoreObtained int       `gorm:"not null" json:"score_obtained"`
	ResultStatus  string    `gorm:"size:50;not null" json:"result_status"`
	AttemptedAt   time.Time `json:"attempted_at"`

	Assessment Assessment                `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
	User       User                      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Answers    []AssessmentAttemptAnswer `gorm:"foreignKey:AttemptID" json:"answers,omitempty"`
}
