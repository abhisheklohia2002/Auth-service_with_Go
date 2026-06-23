package services_module

import (
	"context"
	"errors"
	"log"

	"mime/multipart"
	"path/filepath"
	"strings"

	"example.com/m/internal/dto"
	interfaces "example.com/m/internal/interface"
	"example.com/m/internal/models"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_module "example.com/m/internal/repositories/module"
	repositories_moduleDocument "example.com/m/internal/repositories/module_document"
)

type ModuleService struct {
	moduleRepo   *repositories_module.ModuleRepository
	courseRepo   *repositories_course.CourseRepository
	uploader     interfaces.FileUploader
	documentRepo *repositories_moduleDocument.ModuleDocumentRepository
}

func NewModuleService(
	moduleRepo *repositories_module.ModuleRepository,
	courseRepo *repositories_course.CourseRepository,
	uploader interfaces.FileUploader,
	documentRepo *repositories_moduleDocument.ModuleDocumentRepository,
) *ModuleService {
	return &ModuleService{
		moduleRepo:   moduleRepo,
		courseRepo:   courseRepo,
		uploader:     uploader,
		documentRepo: documentRepo,
	}
}

func (s *ModuleService) UploadPDF(
	ctx context.Context,
	moduleID uint,
	title string,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	oldPublicID string,
	thumbnail multipart.File,
	thumbnailHeader *multipart.FileHeader,
	oldThumbnailPublicID string,
) (*models.ModuleDocument, error) {
	if file == nil || fileHeader == nil {
		return nil, errors.New("file is required")
	}

	if thumbnail == nil || thumbnailHeader == nil {
		return nil, errors.New("thumbnail image is required")
	}

	module, err := s.moduleRepo.FindByID(moduleID)
	if err != nil {
		return nil, err
	}

	if module == nil {
		return nil, errors.New("module not found")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".pdf" {
		return nil, errors.New("only PDF files are allowed")
	}

	if fileHeader.Size > 25*1024*1024 {
		return nil, errors.New("PDF size must be less than 25MB")
	}

	thumbnailExt := strings.ToLower(filepath.Ext(thumbnailHeader.Filename))
	if thumbnailExt != ".jpg" &&
		thumbnailExt != ".jpeg" &&
		thumbnailExt != ".png" &&
		thumbnailExt != ".webp" {
		return nil, errors.New("only JPG, PNG, or WEBP thumbnail images are allowed")
	}

	if thumbnailHeader.Size > 5*1024*1024 {
		return nil, errors.New("thumbnail size must be less than 5MB")
	}

	if title == "" {
		title = strings.TrimSuffix(fileHeader.Filename, ext)
	}

	uploadResult, err := s.uploader.UploadPDF(
		ctx,
		file,
		fileHeader.Filename,
		"lms/module-documents",
	)
	if err != nil {
		return nil, err
	}

	thumbnailUploadResult, err := s.uploader.UploadImage(
		ctx,
		thumbnail,
		thumbnailHeader.Filename,
		"lms/module-thumbnails",
	)
	if err != nil {
		_ = s.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	document := models.ModuleDocument{
		ModuleID:          moduleID,
		Title:             title,
		FileName:          fileHeader.Filename,
		FileURL:           uploadResult.URL,
		PublicID:          uploadResult.PublicID,
		FileType:          "pdf",
		FileSize:          fileHeader.Size,
		ThumbnailName:     thumbnailHeader.Filename,
		ThumbnailURL:      thumbnailUploadResult.URL,
		ThumbnailPublicID: thumbnailUploadResult.PublicID,
		ThumbnailSize:     thumbnailHeader.Size,
		IsActive:          true,
	}

	createdDocument, err := s.documentRepo.Create(&document)
	if err != nil {
		_ = s.uploader.Delete(ctx, uploadResult.PublicID)
		_ = s.uploader.Delete(ctx, thumbnailUploadResult.PublicID)
		return nil, err
	}

	if oldPublicID != "" {
		if err := s.documentRepo.DeleteDocumentByPublicId(oldPublicID); err != nil {
			log.Printf("failed to delete old PDF from database: %v", err)
		}

		if err := s.uploader.Delete(ctx, oldPublicID); err != nil {
			log.Printf("failed to delete old PDF from Cloudinary: %v", err)
		}
	}

	if oldThumbnailPublicID != "" {
		if err := s.uploader.Delete(ctx, oldThumbnailPublicID); err != nil {
			log.Printf("failed to delete old thumbnail from Cloudinary: %v", err)
		}
	}

	return createdDocument, nil
}

func (s *ModuleService) Create(req dto.CreateModuleRequest) (*models.Module, error) {
	req.ModuleTitle = strings.TrimSpace(req.ModuleTitle)
	req.ModuleDescription = strings.TrimSpace(req.ModuleDescription)

	if req.ModuleTitle == "" {
		return nil, errors.New("module title is required")
	}

	if req.DurationMinutes == 0 {
		return nil, errors.New("module duration minutes must be greater than 0")
	}

	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	existingDuration, err := s.moduleRepo.SumDurationByCourseID(req.CourseID)
	if err != nil {
		return nil, err
	}

	newTotalDuration := existingDuration + req.DurationMinutes

	if newTotalDuration > course.TotalDurationMinutes {
		return nil, errors.New("total module duration cannot exceed course total duration")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	module := models.Module{
		CourseID:          req.CourseID,
		ModuleTitle:       req.ModuleTitle,
		ModuleDescription: req.ModuleDescription,
		SequenceNo:        req.SequenceNo,
		DurationMinutes:   req.DurationMinutes,
		DueDays:           req.DueDays,
		IsActive:          isActive,
	}

	return s.moduleRepo.Create(&module)
}
func (s *ModuleService) FindAll() ([]models.Module, error) {
	return s.moduleRepo.FindAll()
}

func (s *ModuleService) FindByID(id uint) (*models.Module, error) {
	return s.moduleRepo.FindByID(id)
}

func (s *ModuleService) FindByCourseID(courseID uint) ([]models.Module, error) {
	course, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	return s.moduleRepo.FindByCourseID(courseID)
}

func (s *ModuleService) Update(id uint, req dto.UpdateModuleRequest) (*models.Module, error) {
	module, err := s.moduleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, errors.New("module not found")
	}

	targetCourseID := module.CourseID

	if req.CourseID != nil {
		course, err := s.courseRepo.FindByID(*req.CourseID)
		if err != nil {
			return nil, err
		}
		if course == nil {
			return nil, errors.New("course not found")
		}

		targetCourseID = *req.CourseID
	}

	targetDuration := module.DurationMinutes

	if req.DurationMinutes != nil {
		if *req.DurationMinutes == 0 {
			return nil, errors.New("module duration minutes must be greater than 0")
		}

		targetDuration = *req.DurationMinutes
	}

	course, err := s.courseRepo.FindByID(targetCourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	existingDuration, err := s.moduleRepo.SumDurationByCourseIDExcludingModule(targetCourseID, module.ID)
	if err != nil {
		return nil, err
	}

	newTotalDuration := existingDuration + targetDuration

	if newTotalDuration > course.TotalDurationMinutes {
		return nil, errors.New("total module duration cannot exceed course total duration")
	}

	if req.CourseID != nil {
		module.CourseID = *req.CourseID
	}

	if req.ModuleTitle != nil {
		title := strings.TrimSpace(*req.ModuleTitle)
		if title == "" {
			return nil, errors.New("module title cannot be empty")
		}
		module.ModuleTitle = title
	}

	if req.ModuleDescription != nil {
		module.ModuleDescription = strings.TrimSpace(*req.ModuleDescription)
	}

	if req.SequenceNo != nil {
		module.SequenceNo = *req.SequenceNo
	}

	if req.DurationMinutes != nil {
		module.DurationMinutes = *req.DurationMinutes
	}

	if req.DueDays != nil {
		module.DueDays = *req.DueDays
	}

	if req.IsActive != nil {
		module.IsActive = *req.IsActive
	}

	return s.moduleRepo.Update(module)
}

func (s *ModuleService) Delete(id uint) error {
	module, err := s.moduleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if module == nil {
		return errors.New("module not found")
	}

	return s.moduleRepo.Delete(id)
}

func (s *ModuleService) FindByModuleID(moduleID uint) ([]models.ModuleDocument, error) {
	module, err := s.moduleRepo.FindByID(moduleID)
	if err != nil {
		return nil, err
	}

	if module == nil {
		return nil, errors.New("module not found")
	}

	return s.documentRepo.FindByModuleID(moduleID)
}

func (s *ModuleService) FindByIDDocument(id uint) (*models.ModuleDocument, error) {
	document, err := s.documentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if document == nil {
		return nil, errors.New("document not found")
	}

	return document, nil
}

func (s *ModuleService) DeleteByIdDocument(id uint) error {
	document, err := s.documentRepo.FindByID(id)
	if err != nil {
		return err
	}

	if document == nil {
		return errors.New("document not found")
	}

	if document.PublicID != "" {
		if err := s.uploader.Delete(context.Background(), document.PublicID); err != nil {
			return err
		}
	}

	return s.documentRepo.Delete(id)
}

func (s *ModuleService) UploadModuleVideo(
	ctx context.Context,
	courseID uint,
	moduleID uint,
	title string,
	video multipart.File,
	videoHeader *multipart.FileHeader,
	oldVideoPublicID string,
) (*models.ModuleVideo, error) {
	if video == nil || videoHeader == nil {
		return nil, errors.New("video file is required")
	}

	module, err := s.moduleRepo.FindByID(moduleID)
	if err != nil {
		return nil, err
	}

	if module == nil {
		return nil, errors.New("module not found")
	}

	if module.CourseID != courseID {
		return nil, errors.New("module does not belong to course")
	}

	videoExt := strings.ToLower(filepath.Ext(videoHeader.Filename))

	if videoExt != ".mp4" &&
		videoExt != ".mov" &&
		videoExt != ".webm" &&
		videoExt != ".mkv" {
		return nil, errors.New("only MP4, MOV, WEBM, or MKV video files are allowed")
	}

	if videoHeader.Size > 100*1024*1024 {
		return nil, errors.New("video size must be less than 100MB")
	}

	if title == "" {
		title = strings.TrimSuffix(videoHeader.Filename, videoExt)
	}

	uploadResult, err := s.uploader.UploadVideo(
		ctx,
		video,
		videoHeader.Filename,
		"lms/module-videos",
	)
	if err != nil {
		return nil, err
	}

	moduleVideo := models.ModuleVideo{
		CourseID:      courseID,
		ModuleID:      moduleID,
		Title:         title,
		VideoName:     videoHeader.Filename,
		VideoURL:      uploadResult.URL,
		VideoPublicID: uploadResult.PublicID,
		VideoSize:     videoHeader.Size,
		VideoType:     "video",
		IsActive:      true,
	}

	createdVideo, err := s.moduleRepo.UpsertByModuleID(&moduleVideo)
	if err != nil {
		_ = s.uploader.DeleteVideo(ctx, uploadResult.PublicID)
		return nil, err
	}

	if oldVideoPublicID != "" {
		if err := s.uploader.DeleteVideo(ctx, oldVideoPublicID); err != nil {
			log.Printf("failed to delete old video from Cloudinary: %v", err)
		}
	}

	return createdVideo, nil
}

func (s *ModuleService) GetModuleVideo(
	ctx context.Context,
	moduleID uint,
) (*models.ModuleVideo, error) {
	module, err := s.moduleRepo.FindByID(moduleID)
	if err != nil {
		return nil, err
	}

	if module == nil {
		return nil, errors.New("module not found")
	}

	video, err := s.moduleRepo.FindByModuleID(moduleID)
	if err != nil {
		return nil, err
	}

	return video, nil
}
