package models

import "time"

type Department struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	DepartmentName string    `gorm:"column:name;not null;unique" json:"department_name"`
	Description    string    `json:"description"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Users []User `gorm:"foreignKey:DepartmentID" json:"users,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}
