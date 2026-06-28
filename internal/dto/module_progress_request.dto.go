package dto

type UpdateModuleProgressRequest struct {
	Status string `json:"status" binding:"required"`
}

type UpdateVideoProgressRequest struct {
	WatchedSeconds  int `json:"watched_seconds" binding:"required,min=0"`
	DurationSeconds int `json:"duration_seconds" binding:"required,min=1"`
}