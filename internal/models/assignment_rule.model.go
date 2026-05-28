package models

type AssignmentRule struct {
	ID           uint   `gorm:"primaryKey" json:"assignment_rule_id"`
	RuleName     string `gorm:"size:150;not null" json:"rule_name"`
	TriggerEvent string `gorm:"size:100;not null" json:"trigger_event"`
	RoleID       uint   `gorm:"not null;index" json:"role_id"`
	IsActive     bool   `gorm:"default:true" json:"is_active"`

	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}
