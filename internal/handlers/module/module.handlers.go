package handlers_module

import (
	"net/http"
	"strconv"

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

func (ctrl *ModuleHandler) UploadPDF(c *gin.Context) {
	moduleIDParam := c.PostForm("module_id")

	moduleID64, err := strconv.ParseUint(moduleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid module id"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "PDF file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to open uploaded file"})
		return
	}
	defer file.Close()

	thumbnailHeader, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Thumbnail image is required"})
		return
	}

	thumbnailFile, err := thumbnailHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to open thumbnail image"})
		return
	}
	defer thumbnailFile.Close()

	title := c.PostForm("title")
	publicId := c.PostForm("publicId")
	oldThumbnailPublicID := c.PostForm("oldThumbnailPublicId")

	document, err := ctrl.moduleService.UploadPDF(
		c.Request.Context(),
		uint(moduleID64),
		title,
		file,
		fileHeader,
		publicId,
		thumbnailFile,
		thumbnailHeader,
		oldThumbnailPublicID,
	)

	if err != nil {
		status := http.StatusInternalServerError

		switch err.Error() {
		case "module not found", "document not found":
			status = http.StatusNotFound
		case "only PDF files are allowed", "PDF size must be less than 25MB":
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "PDF uploaded successfully",
		"document": document,
	})
}

func (ctrl *ModuleHandler) FindByModuleID(c *gin.Context) {
	moduleIDParam := c.Param("moduleId")

	moduleID64, err := strconv.ParseUint(moduleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid module id"})
		return
	}

	documents, err := ctrl.moduleService.FindByModuleID(uint(moduleID64))
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == "module not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
	})
}

func (ctrl *ModuleHandler) DeleteByIdDocument(c *gin.Context) {
	documentIDParam := c.Param("documentId")

	documentID64, err := strconv.ParseUint(documentIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid document id"})
		return
	}

	err = ctrl.moduleService.DeleteByIdDocument(uint(documentID64))
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == "document not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document deleted successfully",
	})
}
