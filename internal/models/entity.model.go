package models

import "time"

type Entity struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EntityName string    `gorm:"column:name;not null;unique" json:"entity_name"`
	EntityType string    `json:"entity_type"`
	Description string   `json:"description"`
	IsActive    bool     `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Departments []Department `gorm:"foreignKey:EntityID" json:"departments,omitempty"`
}

func (Entity) TableName() string {
	return "entities"
}