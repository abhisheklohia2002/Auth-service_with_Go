package handlers_assessmentrule

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_assessmentrule "example.com/m/internal/services/assessment_rule"

	"github.com/gin-gonic/gin"
)

type AssessmentRuleHandler struct {
	assessmentRuleService *services_assessmentrule.AssessmentRuleService
}

func NewAssessmentRuleHandler(assessmentRuleService *services_assessmentrule.AssessmentRuleService) *AssessmentRuleHandler {
	return &AssessmentRuleHandler{
		assessmentRuleService: assessmentRuleService,
	}
}

func (h *AssessmentRuleHandler) Create(c *gin.Context) {
	var req dto.CreateAssessmentRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	rule, err := h.assessmentRuleService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "assessment rule created successfully", "data": rule})
}

func (h *AssessmentRuleHandler) FindAll(c *gin.Context) {
	rules, err := h.assessmentRuleService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assessment rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment rules fetched successfully", "data": rules})
}

func (h *AssessmentRuleHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	rule, err := h.assessmentRuleService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assessment rule"})
		return
	}

	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment rule fetched successfully", "data": rule})
}

func (h *AssessmentRuleHandler) Update(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateAssessmentRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	rule, err := h.assessmentRuleService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment rule updated successfully", "data": rule})
}

func (h *AssessmentRuleHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.assessmentRuleService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment rule deleted successfully"})
}

func (h *AssessmentRuleHandler) CreateByCourseID(c *gin.Context) {
	courseID, ok := helper.ParseUintParam(c, "courseId")
	if !ok {
		return
	}

	var req dto.CreateAssessmentRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	rule, err := h.assessmentRuleService.CreateByCourseID(courseID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "assessment rule created successfully",
		"data":    rule,
	})
}
