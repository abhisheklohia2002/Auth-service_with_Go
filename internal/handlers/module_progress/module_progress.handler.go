package handlers_moduleprogress

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_moduleprogress "example.com/m/internal/services/module_progress"

	"github.com/gin-gonic/gin"
)

type ModuleProgressHandler struct {
	moduleProgressService *services_moduleprogress.ModuleProgressService
}

func NewModuleProgressHandler(moduleProgressService *services_moduleprogress.ModuleProgressService) *ModuleProgressHandler {
	return &ModuleProgressHandler{
		moduleProgressService: moduleProgressService,
	}
}

func (h *ModuleProgressHandler) FindByAssignmentID(c *gin.Context) {
	assignmentID, ok := helper.ParseUintParam(c, "assignmentId")
	if !ok {
		return
	}

	progresses, err := h.moduleProgressService.FindByAssignmentID(assignmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "module progress fetched successfully",
		"data":    progresses,
	})
}

func (h *ModuleProgressHandler) UpdateStatus(c *gin.Context) {
	progressID, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateModuleProgressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	progress, err := h.moduleProgressService.UpdateStatus(progressID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "module progress updated successfully",
		"data":    progress,
	})
}
