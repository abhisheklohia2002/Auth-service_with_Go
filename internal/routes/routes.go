package routes

import (
	"example.com/m/internal/enums"
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
	handlerNotification "example.com/m/internal/handlers/notifications"
	handlers_role "example.com/m/internal/handlers/role"
	handlerSSE "example.com/m/internal/handlers/sse"
	handlers_trainingassignment "example.com/m/internal/handlers/training_assignment"
	handlers_trainingmapping "example.com/m/internal/handlers/training_mapping"
	handlers_trainingsession "example.com/m/internal/handlers/training_session"
	"example.com/m/internal/middleware"
	"example.com/m/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler,
	tokenService *services.TokenService,
	authMiddleware *middleware.AuthMiddleware,
	courseHandler *handlers_course.CourseHandler,
	trainingMappingHandler *handlers_trainingmapping.TrainingMappingHandler,
	trainingAssignmentHandler *handlers_trainingassignment.TrainingAssignmentHandler,
	moduleProgressHandler *handlers_moduleprogress.ModuleProgressHandler,
	moduleHandler *handlers_module.ModuleHandler,
	assessmentRuleHandler *handlers_assessmentrule.AssessmentRuleHandler,
	assessmentHandler *handlers_assessment.AssessmentHandler,
	assessmentAttemptHandler *handlers_assessmentattempt.AssessmentAttemptHandler,
	certificationRuleHandler *handlers_certificationrule.CertificationRuleHandler,
	certificationHandler *handlers_certification.CertificationHandler,
	certificateIssueHandler *handlers_certificateIssue.CertificateIssueHandler,
	roleHandler *handlers_role.RoleHandler,
	assessmentQuestionHandler *handlers_assessmentquestion.AssessmentQuestionHandler,
	assignmentRuleHandler *handlers_assignments.AssignmentRuleHandler,
	trainingsession *handlers_trainingsession.TrainingSessionHandler,
	attendanceHandler *handlers_attendance.AttendanceHandler,
	departmentHandler *handlers_department.DepartmentHandler,
	departmentTrainingMappingHandler *handlers_department_training_mapping.DepartmentTrainingMappingHandler,
	entityHandler *handlers_entity.EntityHandler,
	sseHandler *handlerSSE.SSEHandler,
	notificationHandler *handlerNotification.NotificationHandler,
) {
	api := router.Group("/api")
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST(
			"/create",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			authHandler.CreateUser,
		)
		auth.POST("/login", authHandler.Login)
		auth.GET("/users", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), authHandler.UsersList)
		auth.DELETE("/users/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin)), authHandler.DeleteUserById)
		auth.PUT("/users/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), authHandler.UpdateUserById)
		auth.GET("self", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Employee), string(enums.Manager)), authHandler.Self)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Employee), string(enums.Manager)), authHandler.Logout)
		auth.POST("/bulk-users", authHandler.CreateBulkUsers)
	}

	courses := api.Group("/courses")
	{
		courses.POST(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			courseHandler.CreateCourse,
		)

		courses.GET(
			"",
			authMiddleware.IsAuthMiddleware(),
			courseHandler.GetCourses,
		)

		courses.GET(
			"/:id",
			authMiddleware.IsAuthMiddleware(),
			courseHandler.GetCourseByID,
		)

		courses.PUT(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			courseHandler.UpdateCourse,
		)

		courses.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			courseHandler.DeleteCourse,
		)
	}

	trainingMappings := api.Group("/training-mappings")
	trainingMappings.Use(authMiddleware.IsAuthMiddleware())
	{
		trainingMappings.POST("", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), trainingMappingHandler.Create)
		trainingMappings.GET("", trainingMappingHandler.FindAll)
		trainingMappings.GET("/:id", trainingMappingHandler.FindByID)
		trainingMappings.GET("/role/:roleId", trainingMappingHandler.FindByRoleID)
		trainingMappings.PUT("/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), trainingMappingHandler.Update)
		trainingMappings.DELETE("/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin)), trainingMappingHandler.Delete)
	}

	trainingAssignments := api.Group("/training-assignments")
	{
		trainingAssignments.POST(
			"/manual",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			trainingAssignmentHandler.CreateManual,
		)

		trainingAssignments.POST(
			"/auto",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			trainingAssignmentHandler.AutoAssignByUserRole,
		)

		trainingAssignments.GET(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			trainingAssignmentHandler.FindAll,
		)

		trainingAssignments.GET(
			"/:id",
			authMiddleware.IsAuthMiddleware(),
			trainingAssignmentHandler.FindByID,
		)

		trainingAssignments.GET(
			"/user/:userId",
			authMiddleware.IsAuthMiddleware(),
			trainingAssignmentHandler.FindByUserID,
		)

		trainingAssignments.PATCH(
			"/:id/status",
			authMiddleware.IsAuthMiddleware(),
			trainingAssignmentHandler.UpdateStatus,
		)
		trainingAssignments.PATCH(
			"/:id/reactivate",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			trainingAssignmentHandler.Reactivate,
		)
		trainingAssignments.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			trainingAssignmentHandler.Delete,
		)
		trainingAssignments.POST("/department",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), trainingAssignmentHandler.AssignCourseToDepartment)
		trainingAssignments.GET("/department",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),

			trainingAssignmentHandler.FindDepartmentAssignments)
	}

	moduleProgress := api.Group("/module-progress")
	moduleProgress.Use(authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager), string(enums.Employee)))
	{
		moduleProgress.GET(
			"/assignment/:assignmentId",
			authMiddleware.IsAuthMiddleware(),
			moduleProgressHandler.FindByAssignmentID,
		)

		moduleProgress.PATCH(
			"/:id/status",
			authMiddleware.IsAuthMiddleware(),
			moduleProgressHandler.UpdateStatus,
		)
	}

	modules := api.Group("/modules")
	{
		modules.POST(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			moduleHandler.Create,
		)

		modules.GET(
			"",
			authMiddleware.IsAuthMiddleware(),
			moduleHandler.FindAll,
		)

		modules.GET(
			"/:id",
			authMiddleware.IsAuthMiddleware(),
			moduleHandler.FindByID,
		)

		modules.GET(
			"/course/:courseId",
			authMiddleware.IsAuthMiddleware(),
			moduleHandler.FindByCourseID,
		)

		modules.PUT(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			moduleHandler.Update,
		)

		modules.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			moduleHandler.Delete,
		)
	}

	assessmentRules := api.Group("/assessment-rules")
	{
		assessmentRules.POST(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentRuleHandler.Create,
		)

		assessmentRules.GET(
			"",
			authMiddleware.IsAuthMiddleware(),
			assessmentRuleHandler.FindAll,
		)

		assessmentRules.GET(
			"/:id",
			authMiddleware.IsAuthMiddleware(),
			assessmentRuleHandler.FindByID,
		)

		assessmentRules.PUT(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentRuleHandler.Update,
		)

		assessmentRules.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			assessmentRuleHandler.Delete,
		)
		assessmentRules.POST(
			"/course/:courseId",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			assessmentRuleHandler.CreateByCourseID,
		)
	}

	assessments := api.Group("/assessments")
	{
		assessments.POST(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentHandler.Create,
		)

		assessments.GET(
			"",
			authMiddleware.IsAuthMiddleware(),
			assessmentHandler.FindAll,
		)

		assessments.GET(
			"/:id",
			authMiddleware.IsAuthMiddleware(),
			assessmentHandler.FindByID,
		)

		assessments.GET(
			"/course/:courseId",
			authMiddleware.IsAuthMiddleware(),
			assessmentHandler.FindByCourseID,
		)

		assessments.GET(
			"/module/:moduleId",
			authMiddleware.IsAuthMiddleware(),
			assessmentHandler.FindByModuleID,
		)

		assessments.PUT(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentHandler.Update,
		)

		assessments.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			assessmentHandler.Delete,
		)
	}

	assessmentAttempts := api.Group("/assessment-attempts")
	{
		assessmentAttempts.POST(
			"",
			authMiddleware.IsAuthMiddleware(),
			assessmentAttemptHandler.Create,
		)

		assessmentAttempts.POST(
			"/submit",
			authMiddleware.IsAuthMiddleware(),
			assessmentAttemptHandler.Submit,
		)
		assessmentAttempts.GET(
			"/user/:userId",
			authMiddleware.IsAuthMiddleware(),
			assessmentAttemptHandler.FindByUserID,
		)

		assessmentAttempts.GET(
			"/assessment/:assessmentId",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentAttemptHandler.FindByAssessmentID,
		)

		assessmentAttempts.GET(
			"/user/:userId/assessment/:assessmentId",
			authMiddleware.IsAuthMiddleware(),
			assessmentAttemptHandler.FindByUserAndAssessment,
		)

	}

	certificationRules := api.Group("/certification-rules")
	{
		certificationRules.POST(
			"/course/:courseId",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			certificationRuleHandler.CreateByCourseID,
		)
		certificationRules.POST("", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), certificationRuleHandler.Create)
		certificationRules.GET("", authMiddleware.IsAuthMiddleware(), certificationRuleHandler.FindAll)
		certificationRules.GET("/:id", authMiddleware.IsAuthMiddleware(), certificationRuleHandler.FindByID)
		certificationRules.PUT("/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), certificationRuleHandler.Update)
		certificationRules.DELETE("/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin)), certificationRuleHandler.Delete)
	}

	certifications := api.Group("/certifications")
	{
		certifications.POST("", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), certificationHandler.Create)
		certifications.GET("", authMiddleware.IsAuthMiddleware(), certificationHandler.FindAll)
		certifications.GET("/:id", authMiddleware.IsAuthMiddleware(), certificationHandler.FindByID)
		certifications.GET("/course/:courseId", authMiddleware.IsAuthMiddleware(), certificationHandler.FindByCourseID)
		certifications.PUT("/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), certificationHandler.Update)
		certifications.DELETE("/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin)), certificationHandler.Delete)
	}

	certificateIssues := api.Group("/certificate-issues")
	{
		certificateIssues.POST("/issue", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager), string(enums.Employee)), certificateIssueHandler.Issue)
		certificateIssues.GET("/:id", authMiddleware.IsAuthMiddleware(), certificateIssueHandler.FindByID)
		certificateIssues.GET("/user/:userId", authMiddleware.IsAuthMiddleware(), certificateIssueHandler.FindByUserID)
		certificateIssues.GET("/:id/download", authMiddleware.IsAuthMiddleware(), certificateIssueHandler.DownloadPDF)
	}
	certificates := api.Group("/certificates")
	{
		certificates.GET("/verify/:certificateNumber", certificateIssueHandler.VerifyCertificate)
	}

	roles := api.Group("/roles")
	{
		roles.GET(
			"",
			authMiddleware.IsAuthMiddleware(),
			roleHandler.FindAll,
		)

		roles.GET(
			"/:id",
			authMiddleware.IsAuthMiddleware(),
			roleHandler.FindByID,
		)

		roles.POST(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			roleHandler.Create,
		)

		roles.PUT(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			roleHandler.Update,
		)

		roles.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			roleHandler.Delete,
		)
	}

	assessmentQuestions := api.Group("/assessment-questions")
	{
		assessmentQuestions.POST(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentQuestionHandler.Create,
		)

		assessmentQuestions.POST(
			"/:assessmentId/questions/bulk-upload",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentQuestionHandler.CreateBulkQuestions,
		)

		assessmentQuestions.GET(
			"/assessment/:assessmentId",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentQuestionHandler.FindByAssessmentID,
		)

		assessmentQuestions.GET(
			"/assessment/:assessmentId/learner",
			authMiddleware.IsAuthMiddleware(),
			assessmentQuestionHandler.FindLearnerQuestions,
		)

		assessmentQuestions.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assessmentQuestionHandler.Delete,
		)
	}

	assignmentRules := api.Group("/assignment-rules")
	{
		assignmentRules.POST(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assignmentRuleHandler.Create,
		)

		assignmentRules.GET(
			"",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assignmentRuleHandler.FindAll,
		)

		assignmentRules.GET(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager), string(enums.Employee)),
			assignmentRuleHandler.FindByID,
		)

		assignmentRules.PUT(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)),
			assignmentRuleHandler.Update,
		)

		assignmentRules.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			assignmentRuleHandler.Delete,
		)
	}

	trainingSessions := api.Group("/training-sessions")
	{
		trainingSessions.POST("", trainingsession.Create)
		trainingSessions.GET("", trainingsession.GetAll)
		trainingSessions.GET("/:id", trainingsession.GetByID)
		trainingSessions.PUT("/:id", trainingsession.Update)
		trainingSessions.DELETE("/:id", trainingsession.Delete)
	}

	attendences := api.Group("/attendances")
	{
		attendences.POST("/training-sessions/:sessionId/attendance", attendanceHandler.Mark)
		attendences.POST("/training-sessions/:sessionId/attendance/bulk", attendanceHandler.BulkMark)
		attendences.GET("/training-sessions/:sessionId/attendance", attendanceHandler.GetBySession)
		attendences.GET("/users/:userId/attendance", attendanceHandler.GetByUser)
		attendences.PUT("/attendance/:attendanceId", attendanceHandler.Update)
	}

	moduleDocuments := api.Group("/module-documents")
	{
		moduleDocuments.POST("/upload", moduleHandler.UploadPDF)
		moduleDocuments.GET("/module/:moduleId", moduleHandler.FindByModuleID)
		moduleDocuments.DELETE("/:documentId", moduleHandler.DeleteByIdDocument)
		moduleDocuments.POST("/:moduleId/video", moduleHandler.UploadModuleVideo)
		moduleDocuments.GET("/:moduleId/video", moduleHandler.GetModuleVideo)
		moduleDocuments.GET("/video-upload-tasks/:taskId", moduleHandler.GetVideoUploadTaskStatus)
		moduleDocuments.GET("/pdf-upload-tasks/:taskId", moduleHandler.GetPDFUploadTaskStatus)
	}

	departments := api.Group("/departments")
	{
		departments.POST("", departmentHandler.Create)
		departments.GET("", departmentHandler.GetAll)
		departments.GET("/:id", departmentHandler.GetByID)
		departments.PUT("/:id", departmentHandler.Update)
		departments.DELETE("/:id", departmentHandler.Delete)
		departments.GET("/:id/training-mappings", departmentTrainingMappingHandler.GetByDepartmentID)
		departments.POST(
			"/:department_id/users/bulk-upload",
			departmentTrainingMappingHandler.BulkUploadUsersToDepartment,
		)
	}

	departmentMappings := api.Group("/department-training-mappings")
	{
		departmentMappings.POST("", departmentTrainingMappingHandler.Create)
		departmentMappings.GET("", departmentTrainingMappingHandler.GetAll)
	}

	entities := api.Group("/entities")
	{
		entities.POST("", entityHandler.CreateEntity)
		entities.GET("", entityHandler.GetEntities)
		entities.GET("/:id", entityHandler.GetEntityByID)
		entities.PUT("/:id", entityHandler.UpdateEntity)
		entities.DELETE("/:id", entityHandler.DeleteEntity)
	}

	admin := api.Group("/admin")
	{
		admin.POST("/notifications", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), notificationHandler.CreateNotification)
	}

	notificationUser := api.Group("")
	notificationUser.Use(authMiddleware.IsAuthMiddleware(
		string(enums.Admin),
		string(enums.Manager),
		string(enums.Employee),
	))
	{
		notificationUser.GET("/notifications", notificationHandler.GetMyNotifications)
		notificationUser.PATCH("/notifications/:notificationID/read", notificationHandler.MarkAsRead)
		notificationUser.GET("/notifications/stream", sseHandler.StreamNotifications)
	}

}
