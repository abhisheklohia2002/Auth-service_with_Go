package uploadtaskstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UploadTaskStatus string

const (
	UploadQueued     UploadTaskStatus = "queued"
	UploadProcessing UploadTaskStatus = "processing"
	UploadCompleted  UploadTaskStatus = "completed"
	UploadFailed     UploadTaskStatus = "failed"
)

type UploadTask struct {
	TaskID    string           `json:"task_id"`
	Kind      string           `json:"kind"` 
	Status    UploadTaskStatus `json:"status"`
	Progress  int              `json:"progress"`
	Message   string           `json:"message"`
	Error     string           `json:"error,omitempty"`
	Result    any              `json:"result,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type UploadTaskStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewUploadTaskStore(client *redis.Client) *UploadTaskStore {
	return &UploadTaskStore{
		client: client,
		ttl:    24 * time.Hour,
	}
}

func uploadTaskKey(taskID string) string {
	return fmt.Sprintf("upload_task:%s", taskID)
}

func (s *UploadTaskStore) Create(ctx context.Context, task *UploadTask) error {
	now := time.Now()

	task.CreatedAt = now
	task.UpdatedAt = now

	return s.Save(ctx, task)
}

func (s *UploadTaskStore) Save(ctx context.Context, task *UploadTask) error {
	task.UpdatedAt = time.Now()

	payload, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, uploadTaskKey(task.TaskID), payload, s.ttl).Err()
}

func (s *UploadTaskStore) Get(ctx context.Context, taskID string) (*UploadTask, error) {
	value, err := s.client.Get(ctx, uploadTaskKey(taskID)).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var task UploadTask
	if err := json.Unmarshal([]byte(value), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *UploadTaskStore) Update(
	ctx context.Context,
	taskID string,
	status UploadTaskStatus,
	progress int,
	message string,
	errorMessage string,
	result any,
) error {
	task, err := s.Get(ctx, taskID)
	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("upload task not found")
	}

	task.Status = status
	task.Progress = progress
	task.Message = message
	task.Error = errorMessage

	if result != nil {
		task.Result = result
	}

	return s.Save(ctx, task)
}