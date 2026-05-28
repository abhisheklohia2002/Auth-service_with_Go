package models

type AssessmentQuestion struct {
	ID           uint   `gorm:"primaryKey" json:"question_id"`
	AssessmentID uint   `gorm:"not null;index" json:"assessment_id"`
	QuestionText string `gorm:"type:text;not null" json:"question_text"`
	QuestionType string `gorm:"size:50;not null" json:"question_type"` // single_choice, multiple_choice, true_false, text
	Marks        int    `gorm:"not null;default:1" json:"marks"`
	SequenceNo   int    `gorm:"not null" json:"sequence_no"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`

	Assessment Assessment                 `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
	Options    []AssessmentQuestionOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
}
