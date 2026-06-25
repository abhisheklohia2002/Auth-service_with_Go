package handlers_module

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_module "example.com/m/internal/services/module"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VideoUploadStatus string

const (
	VideoUploadQueued     VideoUploadStatus = "queued"
	VideoUploadProcessing VideoUploadStatus = "processing"
	VideoUploadCompleted  VideoUploadStatus = "completed"
	VideoUploadFailed     VideoUploadStatus = "failed"
)

type VideoUploadTask struct {
	TaskID    string            `json:"task_id"`
	Status    VideoUploadStatus `json:"status"`
	Progress  int               `json:"progress"`
	Message   string            `json:"message"`
	Error     string            `json:"error,omitempty"`
	Video     any               `json:"video,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type PDFUploadStatus string

const (
	PDFUploadQueued     PDFUploadStatus = "queued"
	PDFUploadProcessing PDFUploadStatus = "processing"
	PDFUploadCompleted  PDFUploadStatus = "completed"
	PDFUploadFailed     PDFUploadStatus = "failed"
)

type PDFUploadTask struct {
	TaskID    string          `json:"task_id"`
	Status    PDFUploadStatus `json:"status"`
	Progress  int             `json:"progress"`
	Message   string          `json:"message"`
	Error     string          `json:"error,omitempty"`
	Document  any             `json:"document,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ModuleHandler struct {
	moduleService *services_module.ModuleService
	videoTasks    map[string]*VideoUploadTask
	videoTaskMu   sync.RWMutex
	pdfTasks      map[string]*PDFUploadTask
	pdfTaskMu     sync.RWMutex
}

func NewModuleHandler(moduleService *services_module.ModuleService) *ModuleHandler {
	return &ModuleHandler{
		moduleService: moduleService,
		videoTasks:    make(map[string]*VideoUploadTask),
		pdfTasks:      make(map[string]*PDFUploadTask),
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

	thumbnailHeader, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Thumbnail image is required"})
		return
	}

	title := c.PostForm("title")
	publicId := c.PostForm("publicId")
	oldThumbnailPublicID := c.PostForm("oldThumbnailPublicId")

	taskID := uuid.NewString()

	tempDir := "./tmp/module-pdfs"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create temp directory"})
		return
	}

	pdfTempPath := filepath.Join(tempDir, taskID+"_pdf_"+filepath.Base(fileHeader.Filename))
	thumbnailTempPath := filepath.Join(tempDir, taskID+"_thumb_"+filepath.Base(thumbnailHeader.Filename))

	if err := c.SaveUploadedFile(fileHeader, pdfTempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save PDF temporarily"})
		return
	}

	if err := c.SaveUploadedFile(thumbnailHeader, thumbnailTempPath); err != nil {
		_ = os.Remove(pdfTempPath)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save thumbnail temporarily"})
		return
	}

	now := time.Now()

	task := &PDFUploadTask{
		TaskID:    taskID,
		Status:    PDFUploadQueued,
		Progress:  0,
		Message:   "PDF upload accepted",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctrl.pdfTaskMu.Lock()
	ctrl.pdfTasks[taskID] = task
	ctrl.pdfTaskMu.Unlock()

	go ctrl.processPDFUpload(
		taskID,
		pdfTempPath,
		fileHeader.Filename,
		fileHeader.Size,
		thumbnailTempPath,
		thumbnailHeader.Filename,
		thumbnailHeader.Size,
		uint(moduleID64),
		title,
		publicId,
		oldThumbnailPublicID,
	)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "PDF upload accepted",
		"task_id": taskID,
		"status":  "queued",
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

func (ctrl *ModuleHandler) UploadModuleVideo(c *gin.Context) {
	moduleIDParam := c.Param("moduleId")

	moduleID64, err := strconv.ParseUint(moduleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid module id"})
		return
	}

	courseIDParam := c.PostForm("course_id")

	courseID64, err := strconv.ParseUint(courseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid course id"})
		return
	}

	videoHeader, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Video file is required"})
		return
	}

	title := c.PostForm("title")
	oldVideoPublicID := c.PostForm("oldVideoPublicId")

	taskID := uuid.NewString()

	tempDir := "./tmp/module-videos"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create temp directory"})
		return
	}

	tempPath := filepath.Join(tempDir, taskID+"_"+filepath.Base(videoHeader.Filename))

	if err := c.SaveUploadedFile(videoHeader, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save video temporarily"})
		return
	}

	now := time.Now()

	task := &VideoUploadTask{
		TaskID:    taskID,
		Status:    VideoUploadQueued,
		Progress:  0,
		Message:   "Video upload accepted",
		CreatedAt: now,
		UpdatedAt: now,
	}

	ctrl.videoTaskMu.Lock()
	ctrl.videoTasks[taskID] = task
	ctrl.videoTaskMu.Unlock()

	go ctrl.processModuleVideoUpload(
		taskID,
		tempPath,
		videoHeader.Filename,
		videoHeader.Size,
		uint(courseID64),
		uint(moduleID64),
		title,
		oldVideoPublicID,
	)

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Video upload accepted",
		"task_id": taskID,
		"status":  VideoUploadQueued,
	})
}

func (ctrl *ModuleHandler) GetModuleVideo(c *gin.Context) {
	moduleIDParam := c.Param("moduleId")

	moduleID64, err := strconv.ParseUint(moduleIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid module id"})
		return
	}

	video, err := ctrl.moduleService.GetModuleVideo(
		c.Request.Context(),
		uint(moduleID64),
	)
	if err != nil {
		status := http.StatusInternalServerError

		switch err.Error() {
		case "module not found":
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"video": video,
	})
}

func (ctrl *ModuleHandler) processModuleVideoUpload(
	taskID string,
	tempPath string,
	originalFileName string,
	fileSize int64,
	courseID uint,
	moduleID uint,
	title string,
	oldVideoPublicID string,
) {
	defer func() {
		if err := os.Remove(tempPath); err != nil {
			log.Printf("failed to remove temp video file: %v", err)
		}
	}()

	ctrl.updateVideoTask(taskID, VideoUploadProcessing, 10, "Processing video upload", "", nil)

	file, err := os.Open(tempPath)
	if err != nil {
		ctrl.updateVideoTask(taskID, VideoUploadFailed, 10, "Failed to open temp video", err.Error(), nil)
		return
	}
	defer file.Close()

	fileHeader := &multipart.FileHeader{
		Filename: originalFileName,
		Size:     fileSize,
	}

	ctrl.updateVideoTask(taskID, VideoUploadProcessing, 35, "Uploading video to storage", "", nil)

	video, err := ctrl.moduleService.UploadModuleVideo(
		context.Background(),
		courseID,
		moduleID,
		title,
		file,
		fileHeader,
		oldVideoPublicID,
	)

	if err != nil {
		ctrl.updateVideoTask(taskID, VideoUploadFailed, 70, "Video upload failed", err.Error(), nil)
		return
	}

	ctrl.updateVideoTask(taskID, VideoUploadCompleted, 100, "Video uploaded successfully", "", video)
}

func (ctrl *ModuleHandler) updateVideoTask(
	taskID string,
	status VideoUploadStatus,
	progress int,
	message string,
	errorMessage string,
	video any,
) {
	ctrl.videoTaskMu.Lock()
	defer ctrl.videoTaskMu.Unlock()

	task, ok := ctrl.videoTasks[taskID]
	if !ok {
		return
	}

	task.Status = status
	task.Progress = progress
	task.Message = message
	task.Error = errorMessage
	task.UpdatedAt = time.Now()

	if video != nil {
		task.Video = video
	}
}

func (ctrl *ModuleHandler) GetVideoUploadTaskStatus(c *gin.Context) {
	taskID := c.Param("taskId")

	ctrl.videoTaskMu.RLock()
	task, ok := ctrl.videoTasks[taskID]
	ctrl.videoTaskMu.RUnlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Upload task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ctrl *ModuleHandler) processPDFUpload(
	taskID string,
	pdfTempPath string,
	pdfFileName string,
	pdfSize int64,
	thumbnailTempPath string,
	thumbnailFileName string,
	thumbnailSize int64,
	moduleID uint,
	title string,
	publicId string,
	oldThumbnailPublicID string,
) {
	defer func() {
		if err := os.Remove(pdfTempPath); err != nil {
			log.Printf("failed to remove temp PDF file: %v", err)
		}

		if err := os.Remove(thumbnailTempPath); err != nil {
			log.Printf("failed to remove temp thumbnail file: %v", err)
		}
	}()

	ctrl.updatePDFTask(taskID, PDFUploadProcessing, 10, "Processing PDF upload", "", nil)

	pdfFile, err := os.Open(pdfTempPath)
	if err != nil {
		ctrl.updatePDFTask(taskID, PDFUploadFailed, 10, "Failed to open temp PDF", err.Error(), nil)
		return
	}
	defer pdfFile.Close()

	thumbnailFile, err := os.Open(thumbnailTempPath)
	if err != nil {
		ctrl.updatePDFTask(taskID, PDFUploadFailed, 10, "Failed to open temp thumbnail", err.Error(), nil)
		return
	}
	defer thumbnailFile.Close()

	pdfHeader := &multipart.FileHeader{
		Filename: pdfFileName,
		Size:     pdfSize,
	}

	thumbnailHeader := &multipart.FileHeader{
		Filename: thumbnailFileName,
		Size:     thumbnailSize,
	}

	ctrl.updatePDFTask(taskID, PDFUploadProcessing, 35, "Uploading PDF to storage", "", nil)

	document, err := ctrl.moduleService.UploadPDF(
		context.Background(),
		moduleID,
		title,
		pdfFile,
		pdfHeader,
		publicId,
		thumbnailFile,
		thumbnailHeader,
		oldThumbnailPublicID,
	)

	if err != nil {
		ctrl.updatePDFTask(taskID, PDFUploadFailed, 70, "PDF upload failed", err.Error(), nil)
		return
	}

	ctrl.updatePDFTask(taskID, PDFUploadCompleted, 100, "PDF uploaded successfully", "", document)
}

func (ctrl *ModuleHandler) updatePDFTask(
	taskID string,
	status PDFUploadStatus,
	progress int,
	message string,
	errorMessage string,
	document any,
) {
	ctrl.pdfTaskMu.Lock()
	defer ctrl.pdfTaskMu.Unlock()

	task, ok := ctrl.pdfTasks[taskID]
	if !ok {
		return
	}

	task.Status = status
	task.Progress = progress
	task.Message = message
	task.Error = errorMessage
	task.UpdatedAt = time.Now()

	if document != nil {
		task.Document = document
	}
}

func (ctrl *ModuleHandler) GetPDFUploadTaskStatus(c *gin.Context) {
	taskID := c.Param("taskId")

	ctrl.pdfTaskMu.RLock()
	task, ok := ctrl.pdfTasks[taskID]
	ctrl.pdfTaskMu.RUnlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "PDF upload task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}
