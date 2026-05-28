package handlers_assessmentattempt

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_assessmentattempt "example.com/m/internal/services/assessment_attempt"

	"github.com/gin-gonic/gin"
)

type AssessmentAttemptHandler struct {
	assessmentAttemptService *services_assessmentattempt.AssessmentAttemptService
}

func NewAssessmentAttemptHandler(assessmentAttemptService *services_assessmentattempt.AssessmentAttemptService) *AssessmentAttemptHandler {
	return &AssessmentAttemptHandler{
		assessmentAttemptService: assessmentAttemptService,
	}
}

func (h *AssessmentAttemptHandler) Create(c *gin.Context) {
	var req dto.CreateAssessmentAttemptRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	attempt, err := h.assessmentAttemptService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "assessment attempt created successfully", "data": attempt})
}

func (h *AssessmentAttemptHandler) FindByUserID(c *gin.Context) {
	userID, ok := helper.ParseUintParam(c, "userId")
	if !ok {
		return
	}

	attempts, err := h.assessmentAttemptService.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment attempts fetched successfully", "data": attempts})
}

func (h *AssessmentAttemptHandler) FindByAssessmentID(c *gin.Context) {
	assessmentID, ok := helper.ParseUintParam(c, "assessmentId")
	if !ok {
		return
	}

	attempts, err := h.assessmentAttemptService.FindByAssessmentID(assessmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment attempts fetched successfully", "data": attempts})
}

func (h *AssessmentAttemptHandler) FindByUserAndAssessment(c *gin.Context) {
	userID, ok := helper.ParseUintParam(c, "userId")
	if !ok {
		return
	}

	assessmentID, ok := helper.ParseUintParam(c, "assessmentId")
	if !ok {
		return
	}

	attempts, err := h.assessmentAttemptService.FindByUserAndAssessment(userID, assessmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "assessment attempts fetched successfully", "data": attempts})
}
