package models

import "time"

type ModuleDocument struct {
	ID       uint `gorm:"primaryKey" json:"document_id"`
	ModuleID uint `gorm:"not null;index" json:"module_id"`

	Title    string `gorm:"size:200;not null" json:"title"`
	FileName string `gorm:"size:255;not null" json:"file_name"`

	FileURL  string `gorm:"type:text;not null" json:"file_url"`
	PublicID string `gorm:"size:255" json:"public_id"`

	FileType string `gorm:"size:20;not null;default:'pdf'" json:"file_type"`
	FileSize int64  `gorm:"not null" json:"file_size"`

	PageCount int  `gorm:"default:0" json:"page_count"`
	IsActive  bool `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Module Module `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
}