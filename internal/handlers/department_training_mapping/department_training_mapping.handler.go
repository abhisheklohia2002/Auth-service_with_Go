package handlers_department_training_mapping

import (
	"net/http"
	"strconv"

	services_department_training_mapping "example.com/m/internal/services/department_training_mapping"
	
	"github.com/gin-gonic/gin"
)

type DepartmentTrainingMappingHandler struct {
	service services_department_training_mapping.DepartmentTrainingMappingService
}

func NewDepartmentTrainingMappingHandler(
	service services_department_training_mapping.DepartmentTrainingMappingService,
) *DepartmentTrainingMappingHandler {
	return &DepartmentTrainingMappingHandler{service: service}
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
