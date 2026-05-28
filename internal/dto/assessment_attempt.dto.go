package dto

type CreateAssessmentAttemptRequest struct {
	AssessmentID  uint `json:"assessment_id" binding:"required"`
	UserID        uint `json:"user_id" binding:"required"`
	ScoreObtained int  `json:"score_obtained" binding:"required,min=0"`
}
