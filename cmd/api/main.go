package main

import (
	"log"
	"os"
	"time"

	automigration "example.com/m/internal/Automigration"
	"example.com/m/internal/common/storage"
	"example.com/m/internal/config"
	"example.com/m/internal/db"

	"example.com/m/internal/handlers"
	handlers_assessment "example.com/m/internal/handlers/assessment"
	handlers_assessmentattempt "example.com/m/internal/handlers/assessment_attempt"
	handlers_assessmentquestion "example.com/m/internal/handlers/assessment_question"
	handlers_assessmentrule "example.com/m/internal/handlers/assessment_rule"
	handlers_assignments "example.com/m/internal/handlers/assignments"
	handlers_attendance "example.com/m/internal/handlers/attendance"
	handlers_certificateIssue "example.com/m/internal/handlers/certificate_issue"
	handlers_certification "example.com/m/internal/handlers/certification"
	handlers_certificationrule "example.com/m/internal/handlers/certification_rule"
	handlers_course "example.com/m/internal/handlers/course"
	handlers_department "example.com/m/internal/handlers/department"
	handlers_department_training_mapping "example.com/m/internal/handlers/department_training_mapping"
	handlers_entity "example.com/m/internal/handlers/entity"
	handlers_module "example.com/m/internal/handlers/module"
	handlers_moduleprogress "example.com/m/internal/handlers/module_progress"
	handlers_role "example.com/m/internal/handlers/role"
	handlers_trainingassignment "example.com/m/internal/handlers/training_assignment"
	handlers_trainingmapping "example.com/m/internal/handlers/training_mapping"
	handlers_trainingsession "example.com/m/internal/handlers/training_session"

	"example.com/m/internal/middleware"

	"example.com/m/internal/repositories"
	repositories_role "example.com/m/internal/repositories/Role"
	repositories_assessment "example.com/m/internal/repositories/assessment"
	repositories_assessmentattempt "example.com/m/internal/repositories/assessment_attempt"
	repositories_assessmentattemptanswer "example.com/m/internal/repositories/assessment_attempt_answer"
	repositories_assessmentquestion "example.com/m/internal/repositories/assessment_question"
	repositories_assessmentquestionoption "example.com/m/internal/repositories/assessment_question_option"
	repositories_assessmentrule "example.com/m/internal/repositories/assessment_rule"
	repositories_assignments "example.com/m/internal/repositories/assignments"
	repositories_attendance "example.com/m/internal/repositories/attendance"
	repositories_certificateIssue "example.com/m/internal/repositories/certificate_issue"
	repositories_certification "example.com/m/internal/repositories/certification"
	repositories_certificationrule "example.com/m/internal/repositories/certification_rule"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_department "example.com/m/internal/repositories/department"
	repositories_department_training_mapping "example.com/m/internal/repositories/department_training_mapping"
	repositories_entity "example.com/m/internal/repositories/entity"
	repositories_module "example.com/m/internal/repositories/module"
	repositories_moduleDocument "example.com/m/internal/repositories/module_document"
	repositories_moduleprogress "example.com/m/internal/repositories/module_progress"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
	repositories_trainingmapping "example.com/m/internal/repositories/training_mapping"
	repositories_trainingsession "example.com/m/internal/repositories/training_session"

	"example.com/m/internal/routes"
	"example.com/m/internal/seeders"

	"example.com/m/internal/services"
	services_assessment "example.com/m/internal/services/assessment"
	services_assessmentattempt "example.com/m/internal/services/assessment_attempt"
	services_assessmentquestion "example.com/m/internal/services/assessment_question"
	services_assessmentrule "example.com/m/internal/services/assessment_rule"
	services_assignments "example.com/m/internal/services/assignments"
	services_attendence "example.com/m/internal/services/attendence"
	services_certificateissue "example.com/m/internal/services/certificate_issue"
	services_certificatepdf "example.com/m/internal/services/certificate_pdf"
	services_certification "example.com/m/internal/services/certification"
	services_certificationrule "example.com/m/internal/services/certification_rule"
	services_course "example.com/m/internal/services/course"
	services_department "example.com/m/internal/services/department"
	services_department_training_mapping "example.com/m/internal/services/department_training_mapping"
	services_entity "example.com/m/internal/services/entity"
	services_module "example.com/m/internal/services/module"
	services_moduleprogress "example.com/m/internal/services/module_progress"
	services_role "example.com/m/internal/services/role"
	services_trainingassignment "example.com/m/internal/services/training_assignment"
	services_trainingmapping "example.com/m/internal/services/training_mapping"
	services_trainingsession "example.com/m/internal/services/training_session"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadDotenv()
	// fmt.Printf("config: %+v\n", cfg)
	database := db.SetupDB(cfg)

	if err := automigration.RunAutoMigration(database, cfg.APP_ENV); err != nil {
		log.Fatal("migration failed:", err)
	}

	seeders.SeedRoles(database)
	seeders.SeedTestUsers(database)

	router := gin.Default()

	router.StaticFile(
		"/.well-known/jwks.json",
		"./public/.well-known/jwks.json",
	)

	allowedOrigins := []string{"http://localhost:5173"}

	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}
	log.Printf("Allowed CORS origins: %v\n", allowedOrigins)
	router.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
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
		courseRepo,
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
	cld := config.NewCloudinary()
	fileUploader := storage.NewCloudinaryUploader(cld)
	moduleDocumentRepo := repositories_moduleDocument.NewModuleDocumentRepository(database)
	moduleService := services_module.NewModuleService(moduleRepo, courseRepo, fileUploader, moduleDocumentRepo)
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
		assessmentAttemptRepo,
	)

	trainingAssignmentHandler := handlers_trainingassignment.NewTrainingAssignmentHandler(
		trainingAssignmentService,
	)

	assessmentRuleService := services_assessmentrule.NewAssessmentRuleService(
		assessmentRuleRepo,
		courseRepo,
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
		courseRepo,
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

	trainingSessionRepo := repositories_trainingsession.NewTrainingSessionRepository(database)
	trainingSessionService := services_trainingsession.NewTrainingSessionService(trainingSessionRepo)
	trainingSessionHandler := handlers_trainingsession.NewTrainingSessionHandler(trainingSessionService)

	attendanceRepo := repositories_attendance.NewAttendanceRepository(database)
	attendanceService := services_attendence.NewAttendanceService(attendanceRepo, trainingSessionRepo)
	attendanceHandler := handlers_attendance.NewAttendanceHandler(attendanceService)

	departmentRepo := repositories_department.NewDepartmentRepository(database)
	departmentService := services_department.NewDepartmentService(departmentRepo)
	departmentHandler := handlers_department.NewDepartmentHandler(departmentService)

	departmentTrainingMappingRepo := repositories_department_training_mapping.NewDepartmentTrainingMappingRepository(database)

	departmentTrainingMappingService := services_department_training_mapping.NewDepartmentTrainingMappingService(
		departmentTrainingMappingRepo,
		departmentRepo,
		*courseRepo,
	)

	departmentTrainingMappingHandler := handlers_department_training_mapping.NewDepartmentTrainingMappingHandler(
		departmentTrainingMappingService,
	)

	entityRepo := repositories_entity.NewEntityRepository(database)
	entityService := services_entity.NewEntityService(entityRepo)
	entityHandler := handlers_entity.NewEntityHandler(entityService)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "LMS backend is running",
		})
	})

	router.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})
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
		trainingSessionHandler,
		attendanceHandler,
		departmentHandler,
		departmentTrainingMappingHandler,
		entityHandler,
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
