package controllers_moduledocument

import (

	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	interfaces "example.com/m/internal/interface"
	"example.com/m/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleDocumentController struct {
	DB       *gorm.DB
	Uploader interfaces.FileUploader
}

func NewModuleDocumentController(
	db *gorm.DB,
	uploader interfaces.FileUploader,
) *ModuleDocumentController {
	return &ModuleDocumentController{
		DB:       db,
		Uploader: uploader,
	}
}

func (ctrl *ModuleDocumentController) UploadPDF(c *gin.Context) {
	moduleIDParam := c.Param("moduleId")
	log.Println("Content-Type:", c.GetHeader("Content-Type"))
	log.Println("module_id:", moduleIDParam)
	moduleID, err := strconv.ParseUint(moduleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid module id"})
		return
	}

	var module models.Module
	if err := ctrl.DB.First(&module, moduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Module not found"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "PDF file is required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Only PDF files are allowed"})
		return
	}

	if fileHeader.Size > 25*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "PDF size must be less than 25MB"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to open uploaded file"})
		return
	}
	defer file.Close()

	title := c.PostForm("title")
	if title == "" {
		title = strings.TrimSuffix(fileHeader.Filename, ext)
	}

	uploadResult, err := ctrl.Uploader.UploadPDF(
		c.Request.Context(),
		file,
		fileHeader.Filename,
		"lms/module-documents",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to upload PDF",
			"error":   err.Error(),
		})
		return
	}

	document := models.ModuleDocument{
		ModuleID: uint(moduleID),
		Title:    title,
		FileName: fileHeader.Filename,
		FileURL:  uploadResult.URL,
		PublicID: uploadResult.PublicID,
		FileType: "pdf",
		FileSize: fileHeader.Size,
		IsActive: true,
	}

	if err := ctrl.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "PDF uploaded but failed to save document",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "PDF uploaded successfully",
		"document": document,
	})
}

func (ctrl *ModuleDocumentController) GetDocumentsByModule(c *gin.Context) {
	moduleIDParam := c.Param("moduleId")

	moduleID, err := strconv.ParseUint(moduleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid module id"})
		return
	}

	var documents []models.ModuleDocument

	if err := ctrl.DB.
		Where("module_id = ? AND is_active = ?", moduleID, true).
		Order("created_at ASC").
		Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
	})
}
