package models

import "time"

type Company struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"type:varchar(150);not null" json:"name"`
	UserID uint   `gorm:"not null;index" json:"userId"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateCompanyRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID uint   `json:"userId" binding:"required"`
}
