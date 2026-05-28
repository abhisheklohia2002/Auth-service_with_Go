package handlers_assessment

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_assessment "example.com/m/internal/services/assessment"

	"github.com/gin-gonic/gin"
)

type AssessmentHandler struct {
	assessmentService *services_assessment.AssessmentService
}

func NewAssessmentHandler(assessmentService *services_assessment.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{
		assessmentService: assessmentService,
	}
}

func (h *AssessmentHandler) Create(c *gin.Context) {
	var req dto.CreateAssessmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	assessment, err := h.assessmentService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "assessment created successfully", "data": assessment})
}

func (h *AssessmentHandler) FindAll(c *gin.Context) {
	assessments, err := h.assessmentService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assessments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessments fetched successfully", "data": assessments})
}

func (h *AssessmentHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	assessment, err := h.assessmentService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assessment"})
		return
	}

	if assessment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "assessment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment fetched successfully", "data": assessment})
}

func (h *AssessmentHandler) FindByCourseID(c *gin.Context) {
	courseID, ok := helper.ParseUintParam(c, "courseId")
	if !ok {
		return
	}

	assessments, err := h.assessmentService.FindByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessments fetched successfully", "data": assessments})
}

func (h *AssessmentHandler) FindByModuleID(c *gin.Context) {
	moduleID, ok := helper.ParseUintParam(c, "moduleId")
	if !ok {
		return
	}

	assessments, err := h.assessmentService.FindByModuleID(moduleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessments fetched successfully", "data": assessments})
}

func (h *AssessmentHandler) Update(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateAssessmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	assessment, err := h.assessmentService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment updated successfully", "data": assessment})
}

func (h *AssessmentHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.assessmentService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment deleted successfully"})
}
