package models

type AssessmentAttemptAnswer struct {
	ID         uint `gorm:"primaryKey" json:"answer_id"`
	AttemptID  uint `gorm:"not null;index" json:"attempt_id"`
	QuestionID uint `gorm:"not null;index" json:"question_id"`

	SelectedOptionIDs string `gorm:"type:text" json:"selected_option_ids"`
	TextAnswer        string `gorm:"type:text" json:"text_answer"`

	IsCorrect    bool `gorm:"default:false" json:"is_correct"`
	MarksAwarded int  `gorm:"default:0" json:"marks_awarded"`

	Attempt  AssessmentAttempt  `gorm:"foreignKey:AttemptID" json:"attempt,omitempty"`
	Question AssessmentQuestion `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}
