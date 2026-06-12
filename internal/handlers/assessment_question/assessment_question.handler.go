package handlers_assessmentquestion

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_assessmentquestion "example.com/m/internal/services/assessment_question"

	"github.com/gin-gonic/gin"
)

type AssessmentQuestionHandler struct {
	service *services_assessmentquestion.AssessmentQuestionService
}

func NewAssessmentQuestionHandler(service *services_assessmentquestion.AssessmentQuestionService) *AssessmentQuestionHandler {
	return &AssessmentQuestionHandler{service: service}
}

func (h *AssessmentQuestionHandler) Create(c *gin.Context) {
	var req dto.CreateAssessmentQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	question, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "assessment question created successfully",
		"data":    question,
	})
}

func (h *AssessmentQuestionHandler) FindByAssessmentID(c *gin.Context) {
	assessmentID, ok := helper.ParseUintParam(c, "assessmentId")
	if !ok {
		return
	}

	questions, err := h.service.FindByAssessmentID(assessmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "assessment questions fetched successfully",
		"data":    questions,
	})
}

func (h *AssessmentQuestionHandler) FindLearnerQuestions(c *gin.Context) {
	assessmentID, ok := helper.ParseUintParam(c, "assessmentId")
	if !ok {
		return
	}

	questions, err := h.service.FindLearnerQuestions(assessmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "learner assessment questions fetched successfully",
		"data":    questions,
	})
}

func (h *AssessmentQuestionHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "assessment question deleted successfully",
	})
}

func (h *AssessmentQuestionHandler) CreateBulkQuestions(c *gin.Context) {
	assessmentIDParam := c.Param("assessmentId")
	assessmentID64, err := strconv.ParseUint(assessmentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid assessment id"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "excel file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to open file"})
		return
	}
	defer file.Close()

	result, err := h.service.BulkUploadQuestions(
		c.Request.Context(),
		uint(assessmentID64),
		file,
		fileHeader,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusOK, result)
}
