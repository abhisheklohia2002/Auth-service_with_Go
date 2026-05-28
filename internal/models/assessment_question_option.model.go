package models
type AssessmentQuestionOption struct {
	ID         uint   `gorm:"primaryKey" json:"option_id"`
	QuestionID uint   `gorm:"not null;index" json:"question_id"`
	OptionText string `gorm:"type:text;not null" json:"option_text"`
	IsCorrect  bool   `gorm:"default:false" json:"is_correct"`

	Question AssessmentQuestion `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}