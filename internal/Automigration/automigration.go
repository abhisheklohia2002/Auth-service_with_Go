package automigration

import (
	"example.com/m/internal/models"
	"gorm.io/gorm"
)

func RunAutoMigration(database *gorm.DB, env string) error {
	if env == "production" {
		return nil
	}

	return database.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},

		&models.Course{},
		&models.Module{},

		&models.AssignmentRule{},
		&models.TrainingMapping{},
		&models.TrainingAssignment{},
		&models.ModuleProgress{},

		&models.AssessmentRule{},
		&models.Assessment{},
		&models.AssessmentQuestion{},
		&models.AssessmentQuestionOption{},
		&models.AssessmentAttempt{},
		&models.AssessmentAttemptAnswer{},

		&models.CertificationRule{},
		&models.Certification{},
		&models.CertificateIssue{},

		&models.Notification{},
		&models.TrainingSession{},
		&models.Attendance{},
		&models.ModuleDocument{},

		&models.Department{},
		&models.DepartmentTrainingMapping{},
		&models.Entity{},
		&models.Notification{},
		&models.NotificationRecipient{},
	)
}
