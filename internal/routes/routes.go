package routes

import (
	"example.com/m/internal/enums"
	"example.com/m/internal/handlers"
	handlers_course "example.com/m/internal/handlers/course"
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
}
