package routes

import (
	"example.com/m/internal/enums"
	"example.com/m/internal/handlers"
	"example.com/m/internal/middleware"
	"example.com/m/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler,
	tokenService *services.TokenService,
	companyHandlers *handlers.CompanyHandlers,
	departmentHandlers *handlers.DepartmentHandlers,
	authMiddleware *middleware.AuthMiddleware,
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
	company := api.Group("/company")

	{
		company.POST("/", companyHandlers.CompanyCreate)
	}

	department := api.Group("/departmemt")
	{
		department.POST("/", authMiddleware.IsAuthMiddleware(string(enums.Admin)), departmentHandlers.CreateDepartment)
		department.GET("/", authMiddleware.IsAuthMiddleware(string(enums.Admin), string(enums.Manager)), departmentHandlers.ListDepartment)
		department.GET("/:userId", departmentHandlers.UserIdByDepartment)
	}
}
