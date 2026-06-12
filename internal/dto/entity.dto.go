package dto

type CreateEntityRequest struct {
	EntityName  string `json:"entity_name" binding:"required"`
	EntityType  string `json:"entity_type"`
	Description string `json:"description"`
}

type UpdateEntityRequest struct {
	EntityName  string `json:"entity_name"`
	EntityType  string `json:"entity_type"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}