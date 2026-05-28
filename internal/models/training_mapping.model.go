package models

type TrainingMapping struct {
	ID                uint   `gorm:"primaryKey" json:"mapping_id"`
	RoleID            uint   `gorm:"not null;index" json:"role_id"`
	CourseID          uint   `gorm:"not null;index" json:"course_id"`
	IsMandatory       bool   `gorm:"default:false" json:"is_mandatory"`
	AssignmentTrigger string `gorm:"size:100" json:"assignment_trigger"`
	ActiveFlag        bool   `gorm:"default:true" json:"active_flag"`

	Role   Role   `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Course Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}
