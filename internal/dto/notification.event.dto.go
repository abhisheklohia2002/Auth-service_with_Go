package dto

type NotificationEvent struct {
	NotificationID   uint   `json:"notification_id"`
	NotificationType string `json:"notification_type"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	TargetAudience   string `json:"target_audience"`
	CreatedAt        string `json:"created_at"`
}