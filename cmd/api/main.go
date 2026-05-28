package main

import (
	"log"
	"time"

	"example.com/m/internal/config"
	"example.com/m/internal/db"
	"example.com/m/internal/handlers"
	handlers_assessment "example.com/m/internal/handlers/assessment"
	handlers_assessmentattempt "example.com/m/internal/handlers/assessment_attempt"
	handlers_assessmentrule "example.com/m/internal/handlers/assessment_rule"
	handlers_course "example.com/m/internal/handlers/course"
	handlers_module "example.com/m/internal/handlers/module"
	handlers_moduleprogress "example.com/m/internal/handlers/module_progress"
	handlers_trainingassignment "example.com/m/internal/handlers/training_assignment"
	handlers_trainingmapping "example.com/m/internal/handlers/training_mapping"
	"example.com/m/internal/middleware"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentattempt "example.com/m/internal/repositories/assessment_attempt"
	repositories_assessmentrule "example.com/m/internal/repositories/assessment_rule"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_module "example.com/m/internal/repositories/module"
	repositories_moduleprogress "example.com/m/internal/repositories/module_progress"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
	repository_trainingmapping "example.com/m/internal/repositories/training_mapping"
	"example.com/m/internal/routes"
	"example.com/m/internal/seeders"
	"example.com/m/internal/services"
	services_assessment "example.com/m/internal/services/assessment"
	services_assessmentattempt "example.com/m/internal/services/assessment_attempt"
	services_assessmentrule "example.com/m/internal/services/assessment_rule"
	services_course "example.com/m/internal/services/course"
	services_module "example.com/m/internal/services/module"
	services_moduleprogress "example.com/m/internal/services/module_progress"
	services_trainingassignment "example.com/m/internal/services/training_assignment"
	services_trainingmapping "example.com/m/internal/services/training_mapping"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database := db.SetupDB()
	err := database.AutoMigrate(
		&models.RefreshToken{},
		&models.Module{},
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
		&models.ModuleProgress{},
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

	// course  path
	courseRepo := repositories_course.NewCourseRepository(database)
	courseService := services_course.NewCourseService(courseRepo)
	courseHandler := handlers_course.NewCourseHandler(courseService)

	moduleRepo := repositories_module.NewModuleRepository(database)
	moduleService := services_module.NewModuleService(moduleRepo, courseRepo)
	moduleHandler := handlers_module.NewModuleHandler(moduleService)
	//training Mapping path
	trainingMappingRepo := repository_trainingmapping.NewTrainingMappingRepository(database)

	trainingMappingService := services_trainingmapping.NewTrainingMappingService(
		trainingMappingRepo,
		roleRepo,
		courseRepo,
	)

	//training Assignments Path
	// moduleRepo := repositories_module.NewModuleRepository(database)
	moduleProgressRepo := repositories_moduleprogress.NewModuleProgressRepository(database)

	trainingAssignmentRepo := repositories_trainingassignment.NewTrainingAssignmentRepository(database)

	moduleProgressService := services_moduleprogress.NewModuleProgressService(
		moduleProgressRepo,
		trainingAssignmentRepo,
	)

	moduleProgressHandler := handlers_moduleprogress.NewModuleProgressHandler(moduleProgressService)

	trainingAssignmentService := services_trainingassignment.NewTrainingAssignmentService(
		trainingAssignmentRepo,
		trainingMappingRepo,
		userRepo,
		courseRepo,
		moduleRepo,
		moduleProgressRepo,
	)

	trainingAssignmentHandler := handlers_trainingassignment.NewTrainingAssignmentHandler(trainingAssignmentService)

	trainingMappingHandler := handlers_trainingmapping.NewTrainingMappingHandler(trainingMappingService)

	assessmentRuleRepo := repositories_assessmentrule.NewAssessmentRuleRepository(database)
	assessmentRepo := repositories_assessment.NewAssessmentRepository(database)
	assessmentAttemptRepo := repositories_assessmentattempt.NewAssessmentAttemptRepository(database)

	assessmentRuleService := services_assessmentrule.NewAssessmentRuleService(assessmentRuleRepo)

	assessmentService := services_assessment.NewAssessmentService(
		assessmentRepo,
		assessmentRuleRepo,
		courseRepo,
		moduleRepo,
	)

	assessmentAttemptService := services_assessmentattempt.NewAssessmentAttemptService(
		assessmentAttemptRepo,
		assessmentRepo,
		userRepo,
		moduleProgressRepo,
		trainingAssignmentRepo,
	)
	assessmentRuleHandler := handlers_assessmentrule.NewAssessmentRuleHandler(assessmentRuleService)
	assessmentHandler := handlers_assessment.NewAssessmentHandler(assessmentService)
	assessmentAttemptHandler := handlers_assessmentattempt.NewAssessmentAttemptHandler(assessmentAttemptService)

	routes.SetupRoutes(router, authHandler, tokenService, authMiddleware, courseHandler, trainingMappingHandler, trainingAssignmentHandler, moduleProgressHandler, moduleHandler, assessmentRuleHandler, assessmentHandler, assessmentAttemptHandler)
	router.Run(":" + cfg.Port)
}
