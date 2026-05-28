package dto

type IssueCertificateRequest struct {
	UserID                uint `json:"user_id" binding:"required"`
	CertificationID       uint `json:"certification_id" binding:"required"`
	TrainingAssignmentID uint `json:"training_assignment_id" binding:"required"`
}