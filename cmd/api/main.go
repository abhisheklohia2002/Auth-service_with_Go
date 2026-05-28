package main

import (
	"log"
	"time"

	"example.com/m/internal/config"
	"example.com/m/internal/db"
	"example.com/m/internal/handlers"
	handlers_course "example.com/m/internal/handlers/course"
	handlers_trainingmapping "example.com/m/internal/handlers/training_mapping"
	"example.com/m/internal/middleware"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_course "example.com/m/internal/repositories/course"
	repository_trainingmapping "example.com/m/internal/repositories/training_mapping"
	"example.com/m/internal/routes"
	"example.com/m/internal/seeders"
	"example.com/m/internal/services"
	services_course "example.com/m/internal/services/course"
	services_trainingmapping "example.com/m/internal/services/training_mapping"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database := db.SetupDB()
	err := database.AutoMigrate(
		&models.Company{},
		&models.Department{},
		&models.RefreshToken{},
		&models.Role{},
		&models.User{},
		&models.Course{},
		&models.Module{},
		&models.AssignmentRule{},
		&models.TrainingMapping{},
		&models.TrainingAssignment{},
		&models.AssessmentRule{},
		&models.Assessment{},
		&models.AssessmentAttempt{},
		&models.CertificationRule{},
		&models.Certification{},
		&models.CertificateIssue{},
		&models.Notification{},
	)
	if err != nil {
		panic(err)
	}
	cfg := config.LoadDotenv()
	seeders.SeedRoles(database)
	seeders.SeedTestUsers(database)
	router := gin.Default()
	router.StaticFile(
		"/.well-known/jwks.json",
		"./public/.well-known/jwks.json",
	)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	authMiddleware, err := middleware.NewAuthMiddleware(cfg)
	if err != nil {
		log.Fatal("failed to initialize auth middleware: ", err)
	}
	//auth routes paths
	tokenService, err := services.NewTokenService(cfg)
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repositories.NewUserRepository(database)
	roleRepo := repositories.NewRoleRepository(database)
	authService := services.NewAuthService(userRepo, tokenService, roleRepo)
	authHandler := handlers.NewAuthHandler(authService, tokenService)

	//company  path
	companyRepo := repositories.NewCompanyRepository(database)
	companyService := services.NewCompanyService(companyRepo)
	companyHandler := handlers.NewCompanyHandler(companyService)

	// department  paths
	departmentRepo := repositories.NewDepartmentRepository(database)
	departmentService := services.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentService(departmentService)

	// course  path
	courseRepo := repositories_course.NewCourseRepository(database)
	courseService := services_course.NewCourseService(courseRepo)
	courseHandler := handlers_course.NewCourseHandler(courseService)

	//training Mapping path
	trainingMappingRepo := repository_trainingmapping.NewTrainingMappingRepository(database)

	trainingMappingService := services_trainingmapping.NewTrainingMappingService(
		trainingMappingRepo,
		roleRepo,
		courseRepo,
	)

	trainingMappingHandler := handlers_trainingmapping.NewTrainingMappingHandler(trainingMappingService)

	routes.SetupRoutes(router, authHandler, tokenService, companyHandler, departmentHandler, authMiddleware, courseHandler, trainingMappingHandler)
	router.Run(":" + cfg.Port)
}
