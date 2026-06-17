package dto

type CreateNotificationRequest struct {
	NotificationType   string `json:"notification_type" binding:"required"`
	Title              string `json:"title"`
	Message            string `json:"message" binding:"required"`
	TargetAudience    string `json:"target_audience" binding:"required"` 
	UserIDs            []uint `json:"user_ids"`                          
	AssignmentID       *uint  `json:"assignment_id"`
	CertificateIssueID *uint  `json:"certificate_issue_id"`
}