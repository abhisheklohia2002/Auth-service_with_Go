package seeders

import (
	"log"
	"time"

	"example.com/m/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedTestUsers(db *gorm.DB) {
	users := []struct {
		Name         string
		Email        string
		Password     string
		EmployeeCode string
		RoleName     string
	}{
		{
			Name:         "Manager Test",
			Email:        "manager@test.com",
			Password:     "12345678",
			EmployeeCode: "MGR001",
			RoleName:     "manager",
		},
		{
			Name:         "Admin Test",
			Email:        "admin@test.com",
			Password:     "12345678",
			EmployeeCode: "ADM001",
			RoleName:     "admin",
		},
	}

	for _, item := range users {
		var role models.Role

		if err := db.Where("role_name = ?", item.RoleName).First(&role).Error; err != nil {
			log.Printf("role not found %s: %v", item.RoleName, err)
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to hash password: %v", err)
			continue
		}

		user := models.User{
			FullName:     item.Name,
			Email:        item.Email,
			Password:     string(hashedPassword),
			EmployeeCode: item.EmployeeCode,
			RoleID:       role.ID,
			Status:       "active",
			JoiningDate:  time.Now(),
		}

		err = db.Where("email = ?", item.Email).FirstOrCreate(&user).Error
		if err != nil {
			log.Printf("failed to seed user %s: %v", item.Email, err)
			continue
		}

		log.Printf("test user ready: %s role=%s", item.Email, item.RoleName)
	}
}