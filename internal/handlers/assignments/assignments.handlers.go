package handlers_assignments

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_assignments "example.com/m/internal/services/assignments"
	"github.com/gin-gonic/gin"
)

type AssignmentRuleHandler struct {
	service *services_assignments.AssignmentRuleService
}

func NewAssignmentRuleHandler(service *services_assignments.AssignmentRuleService) *AssignmentRuleHandler {
	return &AssignmentRuleHandler{service: service}
}

func (h *AssignmentRuleHandler) Create(c *gin.Context) {
	var req dto.CreateAssignmentRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	rule, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "assignment rule created successfully",
		"data":    rule,
	})
}

func (h *AssignmentRuleHandler) FindAll(c *gin.Context) {
	rules, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assignment rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "assignment rules fetched successfully",
		"data":    rules,
	})
}

func (h *AssignmentRuleHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	rule, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assignment rule"})
		return
	}

	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assignment rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "assignment rule fetched successfully",
		"data":    rule,
	})
}

func (h *AssignmentRuleHandler) Update(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateAssignmentRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	rule, err := h.service.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "assignment rule updated successfully",
		"data":    rule,
	})
}

func (h *AssignmentRuleHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "assignment rule deleted successfully",
	})
}

func (h *AssignmentRuleHandler) FindByCourseID(c *gin.Context) {
	courseID, ok := helper.ParseUintParam(c, "courseId")
	if !ok {
		return
	}

	rules, err := h.service.FindByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch course assignment rules",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "course assignment rules fetched successfully",
		"data":    rules,
	})
}
