package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	FullName       string    `gorm:"size:150;not null" json:"name"`
	Email          string    `gorm:"size:150;not null;uniqueIndex" json:"email"`
	Password       string    `json:"-"`
	RoleID         uint      `gorm:"not null" json:"role_id"`
	Role           Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Status         string    `gorm:"size:30;default:'active'" json:"status"`
	ManagerID      *uint     `json:"manager_id"`
	Manager        *User     `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	CreatedCourses []Course  `gorm:"foreignKey:CreatedByUserID" json:"created_courses,omitempty"`
	JoiningDate    time.Time `json:"joining_date"`
	EmployeeCode   string    `gorm:"size:50;not null;uniqueIndex" json:"employee_code"`

	TrainingAssignments []TrainingAssignment `gorm:"foreignKey:UserID" json:"training_assignments,omitempty"`
	AssignedTrainings   []TrainingAssignment `gorm:"foreignKey:AssignedByUserID" json:"assigned_trainings,omitempty"`
	AssessmentAttempts  []AssessmentAttempt  `gorm:"foreignKey:UserID" json:"assessment_attempts,omitempty"`
	CertificateIssues   []CertificateIssue   `gorm:"foreignKey:UserID" json:"certificate_issues,omitempty"`
	Notifications       []Notification       `gorm:"foreignKey:UserID" json:"notifications,omitempty"`

	DepartmentID *uint       `json:"department_id"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}
