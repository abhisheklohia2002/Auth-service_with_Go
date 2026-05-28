package handlers_trainingmapping

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	services_trainingmapping "example.com/m/internal/services/training_mapping"

	"github.com/gin-gonic/gin"
)

type TrainingMappingHandler struct {
	trainingMappingService *services_trainingmapping.TrainingMappingService
}

func NewTrainingMappingHandler(trainingMappingService *services_trainingmapping.TrainingMappingService) *TrainingMappingHandler {
	return &TrainingMappingHandler{
		trainingMappingService: trainingMappingService,
	}
}

func (h *TrainingMappingHandler) Create(c *gin.Context) {
	var req dto.CreateTrainingMappingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	mapping, err := h.trainingMappingService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "training mapping created successfully",
		"data":    mapping,
	})
}

func (h *TrainingMappingHandler) FindAll(c *gin.Context) {
	mappings, err := h.trainingMappingService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch training mappings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training mappings fetched successfully",
		"data":    mappings,
	})
}

func (h *TrainingMappingHandler) FindByID(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	mapping, err := h.trainingMappingService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch training mapping",
		})
		return
	}

	if mapping == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "training mapping not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training mapping fetched successfully",
		"data":    mapping,
	})
}

func (h *TrainingMappingHandler) FindByRoleID(c *gin.Context) {
	roleID, ok := parseUintParam(c, "roleId")
	if !ok {
		return
	}

	mappings, err := h.trainingMappingService.FindByRoleID(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training mappings fetched successfully",
		"data":    mappings,
	})
}

func (h *TrainingMappingHandler) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateTrainingMappingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	mapping, err := h.trainingMappingService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training mapping updated successfully",
		"data":    mapping,
	})
}

func (h *TrainingMappingHandler) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		return
	}

	err := h.trainingMappingService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "training mapping deleted successfully",
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
