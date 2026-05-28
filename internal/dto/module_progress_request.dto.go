package dto

type UpdateModuleProgressRequest struct {
	Status string `json:"status" binding:"required"`
}
