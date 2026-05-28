package dto

type CreateRoleRequest struct {
	RoleName    string `json:"role_name" binding:"required"`
	RoleType    string `json:"role_type" binding:"required"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	RoleName    *string `json:"role_name"`
	RoleType    *string `json:"role_type"`
	Description *string `json:"description"`
}
