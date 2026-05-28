package seeders

import (
	"log"

	"example.com/m/internal/models"

	"gorm.io/gorm"
)

func SeedRoles(database *gorm.DB) {
	roles := []models.Role{
		{
			RoleName:    "admin",
			RoleType:    "internal",
			Description: "System administrator with full LMS access",
		},
		{
			RoleName:    "manager",
			RoleType:    "internal",
			Description: "Manager who assigns and tracks trainings",
		},
		{
			RoleName:    "employee",
			RoleType:    "internal",
			Description: "Employee who completes assigned trainings",
		},
	}

	for _, role := range roles {
		err := database.
			Where(models.Role{RoleName: role.RoleName}).
			FirstOrCreate(&role).
			Error

		if err != nil {
			log.Printf("failed to seed role %s: %v", role.RoleName, err)
			continue
		}

		log.Printf("role ready: %s", role.RoleName)
	}
}