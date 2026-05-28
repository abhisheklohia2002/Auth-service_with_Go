package models

type Certification struct {
	ID                uint   `gorm:"primaryKey" json:"certification_id"`
	CourseID          uint   `gorm:"not null;index" json:"course_id"`
	CertificationName string `gorm:"size:200;not null" json:"certification_name"`
	ValidityDays      int    `gorm:"not null" json:"validity_days"`
	RuleID            *uint  `gorm:"index" json:"rule_id"`
	IsActive          bool   `gorm:"default:true" json:"is_active"`

	Course Course             `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Rule   *CertificationRule `gorm:"foreignKey:RuleID" json:"rule,omitempty"`

	CertificateIssues []CertificateIssue `gorm:"foreignKey:CertificationID" json:"certificate_issues,omitempty"`
}