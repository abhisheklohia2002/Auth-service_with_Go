package main

import (
	"log"
	"time"

	"example.com/m/internal/config"
	"example.com/m/internal/db"
	"example.com/m/internal/handlers"
	"example.com/m/internal/middleware"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	"example.com/m/internal/routes"
	"example.com/m/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database := db.SetupDB()
	err := database.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Department{},
		&models.RefreshToken{},
	)
	if err != nil {
		panic(err)
	}
	cfg := config.LoadDotenv()
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
	userRepo := repositories.NewUserRepository(database)
	tokenService, err := services.NewTokenService(cfg)
	if err != nil {
		log.Fatal(err)
	}
	authService := services.NewAuthService(userRepo, tokenService)
	authHandler := handlers.NewAuthHandler(authService, tokenService)

	//company routes path
	companyRepo := repositories.NewCompanyRepository(database)
	companyService := services.NewCompanyService(companyRepo)
	companyHandler := handlers.NewCompanyHandler(companyService)

	// department routes paths
	departmentRepo := repositories.NewDepartmentRepository(database)
	departmentService := services.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentService(departmentService)

	routes.SetupRoutes(router, authHandler, tokenService, companyHandler, departmentHandler, authMiddleware)
	router.Run(":" + cfg.Port)
}
