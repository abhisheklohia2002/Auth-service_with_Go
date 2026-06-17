package serviceNotification

import (
	"context"
	"strconv"
	"time"

	"example.com/m/internal/common/publisher"
	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositoryNotification "example.com/m/internal/repositories/notifications"
)

type NotificationService interface {
	CreateNotification(ctx context.Context, adminID uint, req dto.CreateNotificationRequest) (*models.Notification, error)
	GetUserNotifications(ctx context.Context, userID uint) ([]models.NotificationRecipient, error)
	MarkAsRead(ctx context.Context, notificationID uint, userID uint) error
}

type notificationService struct {
	repo      repositoryNotification.NotificationRepository
	publisher publisher.NotificationPublisher
	userRepo  repositories.UserRepository
}

func NewNotificationService(
	repo repositoryNotification.NotificationRepository,
	publisher publisher.NotificationPublisher,
	userRepo repositories.UserRepository,
) NotificationService {
	return &notificationService{
		repo:      repo,
		publisher: publisher,
		userRepo:  userRepo,
	}
}

func (s *notificationService) CreateNotification(
	ctx context.Context,
	adminID uint,
	req dto.CreateNotificationRequest,
) (*models.Notification, error) {
	now := time.Now()

	notification := &models.Notification{
		CreatedByID:        adminID,
		AssignmentID:       req.AssignmentID,
		CertificateIssueID: req.CertificateIssueID,
		NotificationType:   req.NotificationType,
		Title:              req.Title,
		Message:            req.Message,
		TargetAudience:     req.TargetAudience,
		SentAt:             &now,
	}

	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return nil, err
	}

	var userIDs []uint
	var err error

	if req.TargetAudience == "all" {
		userIDs, err = s.userRepo.GetAllUserIDs(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		userIDs = req.UserIDs
	}

	recipients := make([]models.NotificationRecipient, 0, len(userIDs))

	for _, userID := range userIDs {
		recipients = append(recipients, models.NotificationRecipient{
			NotificationID: notification.ID,
			UserID:         userID,
		})
	}

	if err := s.repo.CreateRecipients(ctx, recipients); err != nil {
		return nil, err
	}

	event := dto.NotificationEvent{
		NotificationID:   notification.ID,
		NotificationType: notification.NotificationType,
		Title:            notification.Title,
		Message:          notification.Message,
		TargetAudience:   notification.TargetAudience,
		CreatedAt:        notification.CreatedAt.Format(time.RFC3339),
	}

	_ = s.publisher.Publish(ctx, "notifications:global", event)

	return notification, nil
}

func (s *notificationService) GetUserNotifications(ctx context.Context, userID uint) ([]models.NotificationRecipient, error) {
	return s.repo.GetUserNotifications(ctx, userID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, notificationID uint, userID uint) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

func uintToString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
