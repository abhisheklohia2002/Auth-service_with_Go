package models

type CertificationRule struct {
	ID uint `gorm:"primaryKey" json:"certification_rule_id"`

	CourseID uint `gorm:"not null;index" json:"course_id"`

	IssueOnCourseCompletion bool `gorm:"default:true" json:"issue_on_course_completion"`
	MinimumScoreRequired   int  `gorm:"not null" json:"minimum_score_required"`
	ValidityDays           int  `gorm:"not null" json:"validity_days"`
	RenewalRequired         bool `gorm:"default:false" json:"renewal_required"`

	Course Course `gorm:"foreignKey:CourseID;references:ID" json:"course,omitempty"`
}