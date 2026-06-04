package handlers_trainingassignment

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	services_trainingassignment "example.com/m/internal/services/training_assignment"

	"github.com/gin-gonic/gin"
)

type TrainingAssignmentHandler struct {
	trainingAssignmentService *services_trainingassignment.TrainingAssignmentService
}

func NewTrainingAssignmentHandler(trainingAssignmentService *services_trainingassignment.TrainingAssignmentService) *TrainingAssignmentHandler {
	return &TrainingAssignmentHandler{
		trainingAssignmentService: trainingAssignmentService,
	}
}

func (h *TrainingAssignmentHandler) CreateManual(c *gin.Context) {
	var req dto.CreateTrainingAssignmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	assignment, err := h.trainingAssignmentService.CreateManual(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "training assignment created successfully",
		"data":    assignment,
	})
}

func (h *TrainingAssignmentHandler) AutoAssignByUserRole(c *gin.Context) {
	var req dto.AutoAssignTrainingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	assignments, err := h.trainingAssignmentService.AutoAssignByUserRole(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "auto training assignment completed",
		"count":   len(assignments),
		"data":    assignments,
	})
}

func (h *TrainingAssignmentHandler) FindAll(c *gin.Context) {
	assignments, err := h.trainingAssignmentService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch training assignments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training assignments fetched successfully",
		"data":    assignments,
	})
}

func (h *TrainingAssignmentHandler) FindByID(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	assignment, err := h.trainingAssignmentService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch training assignment",
		})
		return
	}

	if assignment == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "training assignment not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training assignment fetched successfully",
		"data":    assignment,
	})
}

func (h *TrainingAssignmentHandler) FindByUserID(c *gin.Context) {
	userID, ok := parseUintParam(c, "userId")
	if !ok {
		return
	}

	assignments, err := h.trainingAssignmentService.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training assignments fetched successfully",
		"data":    assignments,
	})
}

func (h *TrainingAssignmentHandler) UpdateStatus(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateTrainingAssignmentStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	assignment, err := h.trainingAssignmentService.UpdateStatus(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training assignment status updated successfully",
		"data":    assignment,
	})
}

func (h *TrainingAssignmentHandler) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	err := h.trainingAssignmentService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training assignment deleted successfully",
	})
}

func parseUintParam(c *gin.Context, paramName string) (uint, bool) {
	rawID := c.Param(paramName)

	id64, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid " + paramName,
		})
		return 0, false
	}

	return uint(id64), true
}

func (h *TrainingAssignmentHandler) AssignCourseToDepartment(c *gin.Context) {
	var req services_trainingassignment.CreateDepartmentTrainingAssignmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.trainingAssignmentService.AssignCourseToDepartment(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "course assigned to department successfully",
		"result":  result,
	})
}

func (h *TrainingAssignmentHandler) FindDepartmentAssignments(c *gin.Context) {
	assignments, err := h.trainingAssignmentService.FindDepartmentAssignments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "department assignments fetched successfully",
		"assignments": assignments,
	})
}
