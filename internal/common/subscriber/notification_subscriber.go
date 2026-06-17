package subscriber

import (
	"context"
	"log"

	"example.com/m/internal/common/sse"
	"github.com/redis/go-redis/v9"
)

type NotificationSubscriber struct {
	redis *redis.Client
	hub   *sse.Hub
}

func NewNotificationSubscriber(redis *redis.Client, hub *sse.Hub) *NotificationSubscriber {
	return &NotificationSubscriber{
		redis: redis,
		hub:   hub,
	}
}

func (s *NotificationSubscriber) Subscribe(ctx context.Context) {
	pubsub := s.redis.Subscribe(ctx, "notifications:global")

	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Println("Subscribed to Redis channel: notifications:global")

	for msg := range ch {
		s.hub.Broadcast([]byte(msg.Payload))
	}
}