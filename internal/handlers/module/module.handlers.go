package handlers_module

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_module "example.com/m/internal/services/module"

	"github.com/gin-gonic/gin"
)

type ModuleHandler struct {
	moduleService *services_module.ModuleService
}

func NewModuleHandler(moduleService *services_module.ModuleService) *ModuleHandler {
	return &ModuleHandler{
		moduleService: moduleService,
	}
}

func (h *ModuleHandler) Create(c *gin.Context) {
	var req dto.CreateModuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	module, err := h.moduleService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "module created successfully",
		"data":    module,
	})
}

func (h *ModuleHandler) FindAll(c *gin.Context) {
	modules, err := h.moduleService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch modules",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "modules fetched successfully",
		"data":    modules,
	})
}

func (h *ModuleHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	module, err := h.moduleService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch module",
		})
		return
	}

	if module == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "module not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "module fetched successfully",
		"data":    module,
	})
}

func (h *ModuleHandler) FindByCourseID(c *gin.Context) {
	courseID, ok := helper.ParseUintParam(c, "courseId")
	if !ok {
		return
	}

	modules, err := h.moduleService.FindByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "modules fetched successfully",
		"data":    modules,
	})
}

func (h *ModuleHandler) Update(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateModuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	module, err := h.moduleService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "module updated successfully",
		"data":    module,
	})
}

func (h *ModuleHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.moduleService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "module deleted successfully",
	})
}
