package publisher

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type NotificationPublisher interface {
	Publish(ctx context.Context, channel string, payload interface{}) error
}

type redisNotificationPublisher struct {
	redis *redis.Client
}

func NewRedisNotificationPublisher(redis *redis.Client) NotificationPublisher {
	return &redisNotificationPublisher{redis: redis}
}

func (p *redisNotificationPublisher) Publish(ctx context.Context, channel string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.redis.Publish(ctx, channel, data).Err()
}