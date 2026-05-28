package services_certification

import (
	"errors"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_certification "example.com/m/internal/repositories/certification"
	repositories_certificationrule "example.com/m/internal/repositories/certification_rule"
	repositories_course "example.com/m/internal/repositories/course"
)

type CertificationService struct {
	certificationRepo     *repositories_certification.CertificationRepository
	certificationRuleRepo *repositories_certificationrule.CertificationRuleRepository
	courseRepo            *repositories_course.CourseRepository
}

func NewCertificationService(
	certificationRepo *repositories_certification.CertificationRepository,
	certificationRuleRepo *repositories_certificationrule.CertificationRuleRepository,
	courseRepo *repositories_course.CourseRepository,
) *CertificationService {
	return &CertificationService{
		certificationRepo:     certificationRepo,
		certificationRuleRepo: certificationRuleRepo,
		courseRepo:            courseRepo,
	}
}

func (s *CertificationService) Create(req dto.CreateCertificationRequest) (*models.Certification, error) {
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	if req.RuleID != nil {
		rule, err := s.certificationRuleRepo.FindByID(*req.RuleID)
		if err != nil {
			return nil, err
		}
		if rule == nil {
			return nil, errors.New("certification rule not found")
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	cert := models.Certification{
		CourseID:          req.CourseID,
		CertificationName: req.CertificationName,
		ValidityDays:      req.ValidityDays,
		RuleID:            req.RuleID,
		IsActive:          isActive,
	}

	return s.certificationRepo.Create(&cert)
}

func (s *CertificationService) FindAll() ([]models.Certification, error) {
	return s.certificationRepo.FindAll()
}

func (s *CertificationService) FindByID(id uint) (*models.Certification, error) {
	return s.certificationRepo.FindByID(id)
}

func (s *CertificationService) FindByCourseID(courseID uint) ([]models.Certification, error) {
	return s.certificationRepo.FindByCourseID(courseID)
}

func (s *CertificationService) Update(id uint, req dto.UpdateCertificationRequest) (*models.Certification, error) {
	cert, err := s.certificationRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, errors.New("certification not found")
	}

	if req.CourseID != nil {
		course, err := s.courseRepo.FindByID(*req.CourseID)
		if err != nil {
			return nil, err
		}
		if course == nil {
			return nil, errors.New("course not found")
		}
		cert.CourseID = *req.CourseID
	}

	if req.CertificationName != nil {
		cert.CertificationName = *req.CertificationName
	}

	if req.ValidityDays != nil {
		if *req.ValidityDays <= 0 {
			return nil, errors.New("validity days must be greater than zero")
		}
		cert.ValidityDays = *req.ValidityDays
	}

	if req.RuleID != nil {
		rule, err := s.certificationRuleRepo.FindByID(*req.RuleID)
		if err != nil {
			return nil, err
		}
		if rule == nil {
			return nil, errors.New("certification rule not found")
		}
		cert.RuleID = req.RuleID
	}

	if req.IsActive != nil {
		cert.IsActive = *req.IsActive
	}

	return s.certificationRepo.Update(cert)
}

func (s *CertificationService) Delete(id uint) error {
	cert, err := s.certificationRepo.FindByID(id)
	if err != nil {
		return err
	}
	if cert == nil {
		return errors.New("certification not found")
	}

	return s.certificationRepo.Delete(id)
}
