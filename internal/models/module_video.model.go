package models

import "time"
type ModuleVideo struct {
	VideoID  uint `gorm:"primaryKey;column:video_id" json:"video_id"`
	CourseID uint `gorm:"column:course_id;not null" json:"course_id"`
	ModuleID uint `gorm:"column:module_id;not null" json:"module_id"`

	Title         string `gorm:"column:title" json:"title"`
	VideoName     string `gorm:"column:video_name" json:"video_name"`
	VideoURL      string `gorm:"column:video_url;type:text;not null" json:"video_url"`
	VideoPublicID string `gorm:"column:video_public_id;not null" json:"video_public_id"`
	VideoSize     int64  `gorm:"column:video_size" json:"video_size"`
	VideoType     string `gorm:"column:video_type" json:"video_type"`

	IsActive  bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ModuleVideo) TableName() string {
	return "module_videos"
}