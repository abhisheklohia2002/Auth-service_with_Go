package services_course

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"example.com/m/internal/dto"
	interfaces "example.com/m/internal/interface"
	"example.com/m/internal/models"
	repositories_course "example.com/m/internal/repositories/course"
	repositories_module "example.com/m/internal/repositories/module"

	"gorm.io/gorm"
)

type CourseService struct {
	courseRepo *repositories_course.CourseRepository
	moduleRepo *repositories_module.ModuleRepository
	uploader   interfaces.FileUploader
}

func NewCourseService(courseRepo *repositories_course.CourseRepository, moduleRepo *repositories_module.ModuleRepository, uploader interfaces.FileUploader) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		moduleRepo: moduleRepo,
		uploader:   uploader,
	}
}

func (s *CourseService) CreateCourse(
	ctx context.Context,
	req dto.CreateCourseRequest,
	createdByUserID uint,
	thumbnailFile multipart.File,
	thumbnailHeader *multipart.FileHeader,
) (*models.Course, error) {
	req.CourseTitle = strings.TrimSpace(req.CourseTitle)
	req.CourseDescription = strings.TrimSpace(req.CourseDescription)
	req.CourseType = strings.TrimSpace(req.CourseType)

	if req.CourseTitle == "" {
		return nil, errors.New("course title is required")
	}

	if req.CourseType == "" {
		return nil, errors.New("course type is required")
	}

	if req.TotalDurationMinutes == 0 {
		return nil, errors.New("total duration minutes must be greater than 0")
	}

	var thumbnailURL string
	var thumbnailPublicID string

	if thumbnailFile != nil && thumbnailHeader != nil {
		if thumbnailHeader.Size > 5*1024*1024 {
			return nil, errors.New("thumbnail size must be less than 5MB")
		}

		contentType := thumbnailHeader.Header.Get("Content-Type")
		if contentType != "image/jpeg" &&
			contentType != "image/png" &&
			contentType != "image/webp" {
			return nil, errors.New("thumbnail must be jpeg, png or webp")
		}

		uploadResult, err := s.uploader.UploadImage(
			ctx,
			thumbnailFile,
			thumbnailHeader.Filename,
			"course-thumbnails",
		)
		if err != nil {
			return nil, errors.New("failed to upload thumbnail")
		}

		thumbnailURL = uploadResult.URL
		thumbnailPublicID = uploadResult.PublicID
	}

	course := &models.Course{
		CourseTitle:          req.CourseTitle,
		CourseDescription:    req.CourseDescription,
		CourseType:           req.CourseType,
		IsActive:             true,
		CreatedByUserID:      createdByUserID,
		TotalDurationMinutes: req.TotalDurationMinutes,
		ThumbnailURL:         thumbnailURL,
		ThumbnailPublicID:    thumbnailPublicID,
	}

	return s.courseRepo.Create(course)
}

func (s *CourseService) GetCourses() ([]models.Course, error) {
	return s.courseRepo.FindAll()
}

func (s *CourseService) GetCourseByID(id uint) (*models.Course, error) {
	return s.courseRepo.FindByID(id)
}

func (s *CourseService) UpdateCourse(
	ctx context.Context,
	id uint,
	req dto.UpdateCourseRequest,
	thumbnailFile multipart.File,
	thumbnailHeader *multipart.FileHeader,
) (*models.Course, error) {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}

	if strings.TrimSpace(req.CourseTitle) != "" {
		course.CourseTitle = strings.TrimSpace(req.CourseTitle)
	}

	if strings.TrimSpace(req.CourseDescription) != "" {
		course.CourseDescription = strings.TrimSpace(req.CourseDescription)
	}

	if strings.TrimSpace(req.CourseType) != "" {
		course.CourseType = strings.TrimSpace(req.CourseType)
	}

	if req.TotalDurationMinutes != nil {
		if *req.TotalDurationMinutes == 0 {
			return nil, errors.New("total duration minutes must be greater than 0")
		}

		totalModuleDuration, err := s.moduleRepo.SumDurationByCourseID(id)
		if err != nil {
			return nil, err
		}

		if *req.TotalDurationMinutes < totalModuleDuration {
			return nil, errors.New("course total duration cannot be less than total module duration")
		}

		course.TotalDurationMinutes = *req.TotalDurationMinutes
	}

	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}
	if s.uploader == nil {
		return nil, errors.New("file uploader is not configured")
	}
	if thumbnailFile != nil && thumbnailHeader != nil {
		if thumbnailHeader.Size > 5*1024*1024 {
			return nil, errors.New("thumbnail size must be less than 5MB")
		}

		contentType := thumbnailHeader.Header.Get("Content-Type")
		if contentType != "image/jpeg" &&
			contentType != "image/png" &&
			contentType != "image/webp" {
			return nil, errors.New("thumbnail must be jpeg, png or webp")
		}

		uploadResult, err := s.uploader.UploadImage(
			ctx,
			thumbnailFile,
			thumbnailHeader.Filename,
			"course-thumbnails",
		)
		if err != nil {
			return nil, errors.New("failed to upload thumbnail")
		}

		if course.ThumbnailPublicID != "" {
			_ = s.uploader.Delete(ctx, course.ThumbnailPublicID)
		}

		course.ThumbnailURL = uploadResult.URL
		course.ThumbnailPublicID = uploadResult.PublicID
	}

	return s.courseRepo.Update(course)
}

func (s *CourseService) DeleteCourse(id uint) error {
	_, err := s.courseRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return err
	}

	return s.courseRepo.Delete(id)
}
