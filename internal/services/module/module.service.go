package services_module

import (
	"context"
	"errors"
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
) (*models.ModuleDocument, error) {
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

	document := models.ModuleDocument{
		ModuleID: moduleID,
		Title:    title,
		FileName: fileHeader.Filename,
		FileURL:  uploadResult.URL,
		PublicID: uploadResult.PublicID,
		FileType: "pdf",
		FileSize: fileHeader.Size,
		IsActive: true,
	}

	createdDocument, err := s.documentRepo.Create(&document)
	if err != nil {
		// Cleanup Cloudinary if DB insert fails
		_ = s.uploader.Delete(ctx, uploadResult.PublicID)
		return nil, err
	}

	return createdDocument, nil
}

func (s *ModuleService) Create(req dto.CreateModuleRequest) (*models.Module, error) {
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
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

	if req.CourseID != nil {
		course, err := s.courseRepo.FindByID(*req.CourseID)
		if err != nil {
			return nil, err
		}
		if course == nil {
			return nil, errors.New("course not found")
		}

		module.CourseID = *req.CourseID
	}

	if req.ModuleTitle != nil {
		module.ModuleTitle = *req.ModuleTitle
	}

	if req.ModuleDescription != nil {
		module.ModuleDescription = *req.ModuleDescription
	}

	if req.SequenceNo != nil {
		module.SequenceNo = *req.SequenceNo
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
