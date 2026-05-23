package models

type Department struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"type:varchar(150);not null" json:"name"`
	UserID uint   `gorm:"not null;index" json:"userId"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
}
type CreateDepartmentRequest struct {
	Name   string `json:"name" binding:"required"`
	UserID uint   `json:"userId"`
}
