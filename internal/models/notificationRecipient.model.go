package models

import "time"

type NotificationRecipient struct {
	ID uint `gorm:"primaryKey" json:"id"`

	NotificationID uint `gorm:"not null;uniqueIndex:idx_notification_user" json:"notification_id"`
	UserID         uint `gorm:"not null;uniqueIndex:idx_notification_user" json:"user_id"`

	ReadStatus bool       `gorm:"default:false;index" json:"read_status"`
	ReadAt     *time.Time `json:"read_at"`

	DeliveredAt *time.Time `json:"delivered_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Notification Notification `gorm:"foreignKey:NotificationID" json:"notification,omitempty"`
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
}