package dto

type CreateDepartmentRequest struct {
	DepartmentName string `json:"department_name" binding:"required"`
	Description    string `json:"description"`
	IsActive       *bool  `json:"is_active"`
	EntityID       uint   `json:"entity_id" binding:"required"`
}

type UpdateDepartmentRequest struct {
	DepartmentName string `json:"department_name"`
	Description    string `json:"description"`
	IsActive       *bool  `json:"is_active"`
	EntityID       *uint  `json:"entity_id"`
}
