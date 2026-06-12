package models

import "time"

type Department struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EntityID       uint      `gorm:"uniqueIndex:idx_entity_department_name" json:"entity_id"`
	DepartmentName string    `gorm:"column:name;not null;uniqueIndex:idx_entity_department_name" json:"department_name"`
	Description    string    `json:"description"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Entity Entity `gorm:"foreignKey:EntityID" json:"entity,omitempty"`
	Users  []User `gorm:"foreignKey:DepartmentID" json:"users,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}
