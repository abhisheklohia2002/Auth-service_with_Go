package main

import (
	"log"
	"time"

	"example.com/m/internal/config"
	"example.com/m/internal/db"

	"example.com/m/internal/handlers"
	handlers_assessment "example.com/m/internal/handlers/assessment"
	handlers_assessmentattempt "example.com/m/internal/handlers/assessment_attempt"
	handlers_assessmentquestion "example.com/m/internal/handlers/assessment_question"
	handlers_assessmentrule "example.com/m/internal/handlers/assessment_rule"
	handlers_assignments "example.com/m/internal/handlers/assignments"
	handlers_certificateIssue "example.com/m/internal/handlers/certificate_issue"
	handlers_certification "example.com/m/internal/handlers/certification"
	handlers_certificationrule "example.com/m/internal/handlers/certification_rule"
	handlers_course "example.com/m/internal/handlers/course"
	handlers_module "example.com/m/internal/handlers/module"
	handlers_moduleprogress "example.com/m/internal/handlers/module_progress"
	handlers_role "example.com/m/internal/handlers/role"
	handlers_trainingassignment "example.com/m/internal/handlers/training_assignment"
	handlers_trainingmapping "example.com/m/internal/handlers/training_mapping"

	"example.com/m/internal/middleware"
	"example.com/m/internal/models"

	"example.com/m/internal/repositories"
	repositories_role "example.com/m/internal/repositories/Role"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentattempt "example.com/m/internal/repositories/assessment_attempt"
	repositories_assessmentattemptanswer "example.com/m/internal/repositories/assessment_attempt_answer"
	repositories_assessmentquestion "example.com/m/internal/repositories/assessment_question"
	repositories_assessmentquestionoption "example.com/m/internal/repositories/assessment_question_option"
	repositories_assessmentrule "example.com/m/internal/repositories/assessment_rule"
	repositories_assignments "example.com/m/internal/repositories/assignments"
	repositories_certificateIssue "example.com/m/internal/repositories/certificate_issue"
	repositories_certification "example.com/m/internal/repositories/certification"
	repositories_certificationrule "example.com/m/internal/repositories/certification_rule"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_module "example.com/m/internal/repositories/module"
	repositories_moduleprogress "example.com/m/internal/repositories/module_progress"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
	repositories_trainingmapping "example.com/m/internal/repositories/training_mapping"

	"example.com/m/internal/routes"
	"example.com/m/internal/seeders"

	"example.com/m/internal/services"
	services_assessment "example.com/m/internal/services/assessment"
	services_assessmentattempt "example.com/m/internal/services/assessment_attempt"
	services_assessmentquestion "example.com/m/internal/services/assessment_question"
	services_assessmentrule "example.com/m/internal/services/assessment_rule"
	services_assignments "example.com/m/internal/services/assignments"
	services_certificateissue "example.com/m/internal/services/certificate_issue"
	services_certificatepdf "example.com/m/internal/services/certificate_pdf"
	services_certification "example.com/m/internal/services/certification"
	services_certificationrule "example.com/m/internal/services/certification_rule"
	services_course "example.com/m/internal/services/course"
	services_module "example.com/m/internal/services/module"
	services_moduleprogress "example.com/m/internal/services/module_progress"
	services_role "example.com/m/internal/services/role"
	services_trainingassignment "example.com/m/internal/services/training_assignment"
	services_trainingmapping "example.com/m/internal/services/training_mapping"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadDotenv()

	database := db.SetupDB()

	err := database.AutoMigrate(

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
	)
	if err != nil {
		log.Fatal("migration failed: ", err)
	}

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

	tokenService, err := services.NewTokenService(cfg)
	if err != nil {
		log.Fatal("failed to initialize token service: ", err)
	}

	userRepo := repositories.NewUserRepository(database)
	roleRepo := repositories.NewRoleRepository(database)

	courseRepo := repositories_course.NewCourseRepository(database)
	moduleRepo := repositories_module.NewModuleRepository(database)

	trainingMappingRepo := repositories_trainingmapping.NewTrainingMappingRepository(database)
	trainingAssignmentRepo := repositories_trainingassignment.NewTrainingAssignmentRepository(database)
	moduleProgressRepo := repositories_moduleprogress.NewModuleProgressRepository(database)

	assignmentRuleRepo := repositories_assignments.NewAssignmentRuleRepository(database)

	assignmentRuleService := services_assignments.NewAssignmentRuleService(
		assignmentRuleRepo,
		roleRepo,
	)

	assignmentRuleHandler := handlers_assignments.NewAssignmentRuleHandler(assignmentRuleService)
	assessmentRuleRepo := repositories_assessmentrule.NewAssessmentRuleRepository(database)
	assessmentRepo := repositories_assessment.NewAssessmentRepository(database)
	assessmentQuestionRepo := repositories_assessmentquestion.NewAssessmentQuestionRepository(database)
	assessmentQuestionOptionRepo := repositories_assessmentquestionoption.NewAssessmentQuestionOptionRepository(database)
	assessmentAttemptRepo := repositories_assessmentattempt.NewAssessmentAttemptRepository(database)
	assessmentAttemptAnswerRepo := repositories_assessmentattemptanswer.NewAssessmentAttemptAnswerRepository(database)

	certificationRuleRepo := repositories_certificationrule.NewCertificationRuleRepository(database)
	certificationRepo := repositories_certification.NewCertificationRepository(database)
	certificateIssueRepo := repositories_certificateIssue.NewCertificateIssueRepository(database)

	roleModuleRepo := repositories_role.NewRoleRepository(database)

	authService := services.NewAuthService(userRepo, tokenService, roleRepo)
	authHandler := handlers.NewAuthHandler(authService, tokenService)

	roleService := services_role.NewRoleService(roleModuleRepo)
	roleHandler := handlers_role.NewRoleHandler(roleService)

	courseService := services_course.NewCourseService(courseRepo)
	courseHandler := handlers_course.NewCourseHandler(courseService)

	moduleService := services_module.NewModuleService(moduleRepo, courseRepo)
	moduleHandler := handlers_module.NewModuleHandler(moduleService)

	trainingMappingService := services_trainingmapping.NewTrainingMappingService(
		trainingMappingRepo,
		roleRepo,
		courseRepo,
	)

	trainingMappingHandler := handlers_trainingmapping.NewTrainingMappingHandler(
		trainingMappingService,
	)

	moduleProgressService := services_moduleprogress.NewModuleProgressService(
		moduleProgressRepo,
		trainingAssignmentRepo,
	)

	moduleProgressHandler := handlers_moduleprogress.NewModuleProgressHandler(
		moduleProgressService,
	)

	trainingAssignmentService := services_trainingassignment.NewTrainingAssignmentService(
		trainingAssignmentRepo,
		trainingMappingRepo,
		userRepo,
		courseRepo,
		moduleRepo,
		moduleProgressRepo,
	)

	trainingAssignmentHandler := handlers_trainingassignment.NewTrainingAssignmentHandler(
		trainingAssignmentService,
	)

	assessmentRuleService := services_assessmentrule.NewAssessmentRuleService(
		assessmentRuleRepo,
	)

	assessmentRuleHandler := handlers_assessmentrule.NewAssessmentRuleHandler(
		assessmentRuleService,
	)

	assessmentService := services_assessment.NewAssessmentService(
		assessmentRepo,
		assessmentRuleRepo,
		courseRepo,
		moduleRepo,
	)

	assessmentHandler := handlers_assessment.NewAssessmentHandler(
		assessmentService,
	)

	assessmentQuestionService := services_assessmentquestion.NewAssessmentQuestionService(
		assessmentQuestionRepo,
		assessmentQuestionOptionRepo,
		assessmentRepo,
	)

	assessmentQuestionHandler := handlers_assessmentquestion.NewAssessmentQuestionHandler(
		assessmentQuestionService,
	)

	assessmentAttemptService := services_assessmentattempt.NewAssessmentAttemptService(
		assessmentAttemptRepo,
		assessmentRepo,
		userRepo,
		moduleProgressRepo,
		trainingAssignmentRepo,
		assessmentQuestionRepo,
		assessmentAttemptAnswerRepo,
	)

	assessmentAttemptHandler := handlers_assessmentattempt.NewAssessmentAttemptHandler(
		assessmentAttemptService,
	)

	certificationRuleService := services_certificationrule.NewCertificationRuleService(
		certificationRuleRepo,
	)

	certificationRuleHandler := handlers_certificationrule.NewCertificationRuleHandler(
		certificationRuleService,
	)

	certificationService := services_certification.NewCertificationService(
		certificationRepo,
		certificationRuleRepo,
		courseRepo,
	)

	certificationHandler := handlers_certification.NewCertificationHandler(
		certificationService,
	)

	certificatePDFService := services_certificatepdf.NewCertificatePDFService(
		"./storage/certificates",
	)

	certificateIssueService := services_certificateissue.NewCertificateIssueService(
		certificateIssueRepo,
		certificationRepo,
		userRepo,
		trainingAssignmentRepo,
		certificatePDFService,
	)

	certificateIssueHandler := handlers_certificateIssue.NewCertificateIssueHandler(
		certificateIssueService,
	)

	routes.SetupRoutes(
		router,
		authHandler,
		tokenService,
		authMiddleware,
		courseHandler,
		trainingMappingHandler,
		trainingAssignmentHandler,
		moduleProgressHandler,
		moduleHandler,
		assessmentRuleHandler,
		assessmentHandler,
		assessmentAttemptHandler,
		certificationRuleHandler,
		certificationHandler,
		certificateIssueHandler,
		roleHandler,
		assessmentQuestionHandler,
		assignmentRuleHandler,
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	log.Println("server running on port:", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
