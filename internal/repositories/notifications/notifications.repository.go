package repositoryNotification

import (
	"context"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *models.Notification) error
	CreateRecipients(ctx context.Context, recipients []models.NotificationRecipient) error
	GetUserNotifications(ctx context.Context, userID uint) ([]models.NotificationRecipient, error)
	MarkAsRead(ctx context.Context, notificationID uint, userID uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepository) CreateRecipients(ctx context.Context, recipients []models.NotificationRecipient) error {
	if len(recipients) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Create(&recipients).Error
}

func (r *notificationRepository) GetUserNotifications(
	ctx context.Context,
	userID uint,
) ([]models.NotificationRecipient, error) {
	var recipients []models.NotificationRecipient

	err := r.db.WithContext(ctx).
		Preload("Notification").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&recipients).Error

	return recipients, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notificationID uint, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.NotificationRecipient{}).
		Where("notification_id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{
			"read_status": true,
			"read_at":     gorm.Expr("NOW()"),
		}).Error
}
