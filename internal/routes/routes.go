package routes

import (
	"example.com/m/internal/handlers"
	"example.com/m/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handlers.AuthHandler,
	tokenService *services.TokenService,
	companyHandlers *handlers.CompanyHandlers,
	departmentHandlers *handlers.DepartmentHandlers,
) {
	api := router.Group("/api/")
	auth := api.Group("/auth")
	{
		auth.POST("/create", authHandler.Register)
		auth.POST("/login",authHandler.Login)
		auth.GET("/users", authHandler.UsersList)
		auth.DELETE("/users/:id", authHandler.DeleteUserById)
		auth.PUT("/users/:id", authHandler.UpdateUserById)

	}
	company := api.Group("/company")

	{
		company.POST("/", companyHandlers.CompanyCreate)
	}

	
	department := api.Group("/departmemt")
	{
		department.POST("/", departmentHandlers.CreateDepartment)
		department.GET("/", departmentHandlers.ListDepartment)
		department.GET("/:userId", departmentHandlers.UserIdByDepartment)
	}
}
