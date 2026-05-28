package models

type Role struct {
	ID          uint   `gorm:"primaryKey" json:"role_id"`
	RoleName    string `gorm:"size:100;not null;uniqueIndex" json:"role_name"`
	RoleType    string `gorm:"size:50;not null" json:"role_type"`
	Description string `gorm:"type:text" json:"description"`

	Users            []User            `gorm:"foreignKey:RoleID" json:"users,omitempty"`
	TrainingMappings []TrainingMapping `gorm:"foreignKey:RoleID" json:"training_mappings,omitempty"`
	AssignmentRules  []AssignmentRule  `gorm:"foreignKey:RoleID" json:"assignment_rules,omitempty"`
}
