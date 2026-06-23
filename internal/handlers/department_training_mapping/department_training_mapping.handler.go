package handlers_department_training_mapping

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/m/internal/services"
	services_department_training_mapping "example.com/m/internal/services/department_training_mapping"

	"github.com/gin-gonic/gin"
)

type DepartmentTrainingMappingHandler struct {
	service     services_department_training_mapping.DepartmentTrainingMappingService
	authService *services.AuthService
}

func NewDepartmentTrainingMappingHandler(
	service services_department_training_mapping.DepartmentTrainingMappingService,
	authService *services.AuthService,
) *DepartmentTrainingMappingHandler {
	return &DepartmentTrainingMappingHandler{service: service, authService: authService}
}

func (h *DepartmentTrainingMappingHandler) Create(c *gin.Context) {
	var req services_department_training_mapping.CreateDepartmentTrainingMappingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	mapping, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "department training mapping created successfully",
		"mapping": mapping,
	})
}

func (h *DepartmentTrainingMappingHandler) GetAll(c *gin.Context) {
	mappings, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mappings": mappings,
	})
}

func (h *DepartmentTrainingMappingHandler) GetByDepartmentID(c *gin.Context) {
	departmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid department id",
		})
		return
	}

	mappings, err := h.service.GetByDepartmentID(uint(departmentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mappings": mappings,
	})
}

func (h *DepartmentTrainingMappingHandler) BulkUploadUsersToDepartment(c *gin.Context) {
	departmentIDParam := strings.TrimSpace(c.Param("department_id"))

	parsedDepartmentID, err := strconv.ParseUint(departmentIDParam, 10, 64)
	if err != nil || parsedDepartmentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "valid department_id is required",
		})
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "file is required",
		})
		return
	}
	defer file.Close()

	departmentID := uint(parsedDepartmentID)

	response, err := h.authService.CreateBulkUsers(
		c.Request.Context(),
		file,
		fileHeader,
		&departmentID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
