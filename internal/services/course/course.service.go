package services_course

import (
	"errors"
	"strings"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_course "example.com/m/internal/repositories/course"

	"gorm.io/gorm"
)

type CourseService struct {
	courseRepo *repositories_course.CourseRepository
}

func NewCourseService(courseRepo *repositories_course.CourseRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
	}
}

func (s *CourseService) CreateCourse(req dto.CreateCourseRequest, createdByUserID uint) (*models.Course, error) {
	req.CourseTitle = strings.TrimSpace(req.CourseTitle)
	req.CourseDescription = strings.TrimSpace(req.CourseDescription)
	req.CourseType = strings.TrimSpace(req.CourseType)

	if req.CourseTitle == "" {
		return nil, errors.New("course title is required")
	}

	if req.CourseType == "" {
		return nil, errors.New("course type is required")
	}

	course := &models.Course{
		CourseTitle:       req.CourseTitle,
		CourseDescription: req.CourseDescription,
		CourseType:        req.CourseType,
		IsActive:          true,
		CreatedByUserID:   createdByUserID,
	}

	return s.courseRepo.Create(course)
}

func (s *CourseService) GetCourses() ([]models.Course, error) {
	return s.courseRepo.FindAll()
}

func (s *CourseService) GetCourseByID(id uint) (*models.Course, error) {
	return s.courseRepo.FindByID(id)
}

func (s *CourseService) UpdateCourse(id uint, req dto.UpdateCourseRequest) (*models.Course, error) {
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

	if req.IsActive != nil {
		course.IsActive = *req.IsActive
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
