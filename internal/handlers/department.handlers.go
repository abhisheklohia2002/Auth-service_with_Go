package handlers

import (
	"strconv"

	"example.com/m/internal/models"
	"example.com/m/internal/services"
	"github.com/gin-gonic/gin"
)

type DepartmentHandlers struct {
	departmentService *services.DepartmentService
}

func NewDepartmentService(departmentService *services.DepartmentService) *DepartmentHandlers {
	return &DepartmentHandlers{departmentService: departmentService}
}

func (d *DepartmentHandlers) CreateDepartment(c *gin.Context) {
	var req models.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	department, err := d.departmentService.CreateDepartment(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "department created successfully",
		"data":    department,
	})
}

func (d *DepartmentHandlers) ListDepartment(c *gin.Context) {
	departments, err := d.departmentService.ListDepartment()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch users",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Departments fetched successfully",
		"data":    departments,
	})
}

func (d *DepartmentHandlers) UserIdByDepartment(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		panic(err)
	}
	departments, err := d.departmentService.UserIdByDepartment(userId)
	if err != nil {
		c.JSON(404, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"departments": departments,
		"userId":      userId,
		"message":     "You have got departments by user Id",
	})
	
}
