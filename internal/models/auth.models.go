package models

import "time"

type RegisterRequest struct {
	Name         string `json:"name"`
	FullName     string `json:"full_name"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	RoleID       uint   `json:"role_id"`
	EmployeeCode string `json:"employee_code"`
	ManagerID    *uint  `json:"manager_id"`
	DepartmentID *uint  `json:"department_id"`
	Status       string `json:"status"`
}
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshToken struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	TokenHash string     `json:"-" gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt *time.Time `json:"revoked_at" gorm:"index"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type AuthResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type UpdateUserRequest struct {
	FullName     string `json:"full_name"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	EmployeeCode string `json:"employee_code"`
	RoleID       uint   `json:"role_id"`
	ManagerID    *uint  `json:"manager_id"`
	DepartmentID *uint  `json:"department_id"`
	Status       string `json:"status"`
}
