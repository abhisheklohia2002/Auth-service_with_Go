package models

type Assessment struct {
	ID             uint   `gorm:"primaryKey" json:"assessment_id"`
	CourseID       uint   `gorm:"not null;index" json:"course_id"`
	ModuleID       *uint  `gorm:"index" json:"module_id"`
	AssessmentTitle string `gorm:"size:200;not null" json:"assessment_title"`
	AssessmentType  string `gorm:"size:100;not null" json:"assessment_type"`
	MaxScore        int    `gorm:"not null" json:"max_score"`
	PassingScore   int    `gorm:"not null" json:"passing_score"`
	RuleID         *uint  `gorm:"index" json:"rule_id"`
	IsActive       bool   `gorm:"default:true" json:"is_active"`

	Course Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Module *Module `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Rule   *AssessmentRule `gorm:"foreignKey:RuleID" json:"rule,omitempty"`

	Attempts []AssessmentAttempt `gorm:"foreignKey:AssessmentID" json:"attempts,omitempty"`
}