package models

type AssessmentRule struct {
	ID               uint   `gorm:"primaryKey" json:"assessment_rule_id"`
	MaxAttempts      int    `gorm:"not null;default:1" json:"max_attempts"`
	PassingScore     int    `gorm:"not null" json:"passing_score"`
	RetakeAllowed    bool   `gorm:"default:false" json:"retake_allowed"`
	EvaluationMethod string `gorm:"size:100" json:"evaluation_method"`

	Assessments []Assessment `gorm:"foreignKey:RuleID" json:"assessments,omitempty"`
}