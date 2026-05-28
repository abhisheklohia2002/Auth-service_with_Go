package routes

import (
	"example.com/m/internal/enums"
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
) {
	api := router.Group("/api/")
	auth := api.Group("/auth")
	{
		auth.POST("/create", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/users", authMiddleware.IsAuthMiddleware(string(enums.Admin)), authHandler.UsersList)
		auth.DELETE("/users/:id", authMiddleware.IsAuthMiddleware(string(enums.Admin)), authHandler.DeleteUserById)
		auth.PUT("/users/:id", authHandler.UpdateUserById)
		auth.GET("self", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Employee), string(enums.Manager)), authHandler.Self)
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

		trainingAssignments.DELETE(
			"/:id",
			authMiddleware.IsAuthMiddleware(string(enums.Admin)),
			trainingAssignmentHandler.Delete,
		)
	}

	moduleProgress := api.Group("/module-progress")
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
}
