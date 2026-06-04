package handlers_trainingsession

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	services_trainingsession "example.com/m/internal/services/training_session"
	"github.com/gin-gonic/gin"
)

type TrainingSessionHandler struct {
	service services_trainingsession.TrainingSessionService
}

func NewTrainingSessionHandler(service services_trainingsession.TrainingSessionService) *TrainingSessionHandler {
	return &TrainingSessionHandler{service: service}
}

func (h *TrainingSessionHandler) Create(c *gin.Context) {
	var req dto.CreateTrainingSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	session, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Training session created successfully",
		"data":    session,
	})
}

func (h *TrainingSessionHandler) GetAll(c *gin.Context) {
	filter := models.TrainingSessionFilter{
		CourseID:    c.Query("course_id"),
		ModuleID:    c.Query("module_id"),
		Status:      c.Query("status"),
		SessionType: c.Query("session_type"),
	}

	sessions, err := h.service.GetAll(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch training sessions",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Training sessions fetched successfully",
		"data":    sessions,
	})
}

func (h *TrainingSessionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid session id",
		})
		return
	}

	session, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Training session not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Training session fetched successfully",
		"data":    session,
	})
}

func (h *TrainingSessionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid session id",
		})
		return
	}

	var req dto.UpdateTrainingSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	session, err := h.service.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Training session updated successfully",
		"data":    session,
	})
}

func (h *TrainingSessionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid session id",
		})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Training session deleted successfully",
	})
}
