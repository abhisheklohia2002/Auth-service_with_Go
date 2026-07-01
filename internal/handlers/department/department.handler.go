package handlers_department

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	services_department "example.com/m/internal/services/department"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service services_department.DepartmentService
}

func NewDepartmentHandler(service services_department.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) Create(c *gin.Context) {
	var req dto.CreateDepartmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	department, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "department created successfully",
		"department": department,
	})
}

func (h *DepartmentHandler) GetAll(c *gin.Context) {
	departments, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"departments": departments,
	})
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid department id",
		})
		return
	}

	department, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "department not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"department": department,
	})
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid department id",
		})
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	department, err := h.service.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "department updated successfully",
		"department": department,
	})
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid department id",
		})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "department deleted successfully",
	})
}

func (h *DepartmentHandler) GetDepartmentsByEntityID(c *gin.Context) {
	entityID, err := strconv.ParseUint(c.Param("entity_id"), 10, 64)
	if err != nil || entityID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid entity id",
		})
		return
	}

	departments, err := h.service.GetDepartmentsByEntityID(uint(entityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch departments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"departments": departments,
	})
}
