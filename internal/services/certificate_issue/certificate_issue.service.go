package services_certificateissue

import (
	"errors"
	"fmt"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	repositories_CertificateIssue "example.com/m/internal/repositories/certificate_issue"
	repositories_certification "example.com/m/internal/repositories/certification"
	repositories_trainingassignment "example.com/m/internal/repositories/training_assignment"
	services_certificatepdf "example.com/m/internal/services/certificate_pdf"
)

type CertificateIssueService struct {
	certificateIssueRepo   *repositories_CertificateIssue.CertificateIssueRepository
	certificationRepo      *repositories_certification.CertificationRepository
	userRepo               *repositories.UserRepository
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository
	pdfService             *services_certificatepdf.CertificatePDFService
}

func NewCertificateIssueService(
	certificateIssueRepo *repositories_CertificateIssue.CertificateIssueRepository,
	certificationRepo *repositories_certification.CertificationRepository,
	userRepo *repositories.UserRepository,
	trainingAssignmentRepo *repositories_trainingassignment.TrainingAssignmentRepository,
	pdfService *services_certificatepdf.CertificatePDFService,
) *CertificateIssueService {
	return &CertificateIssueService{
		certificateIssueRepo:   certificateIssueRepo,
		certificationRepo:      certificationRepo,
		userRepo:               userRepo,
		trainingAssignmentRepo: trainingAssignmentRepo,
		pdfService:             pdfService,
	}
}

func (s *CertificateIssueService) Issue(req dto.IssueCertificateRequest) (*models.CertificateIssue, error) {
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	cert, err := s.certificationRepo.FindByID(req.CertificationID)
	if err != nil {
		return nil, err
	}
	if cert == nil {
		return nil, errors.New("certification not found")
	}
	if !cert.IsActive {
		return nil, errors.New("certification is inactive")
	}

	assignment, err := s.trainingAssignmentRepo.FindByID(req.TrainingAssignmentID)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, errors.New("training assignment not found")
	}

	if assignment.UserID != req.UserID {
		return nil, errors.New("training assignment does not belong to this user")
	}

	if assignment.CourseID != cert.CourseID {
		return nil, errors.New("training assignment course does not match certification course")
	}

	if assignment.Status != "completed" {
		return nil, errors.New("training assignment is not completed")
	}

	exists, err := s.certificateIssueRepo.ExistsByUserAndCertification(req.UserID, req.CertificationID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("certificate already issued for this user and certification")
	}

	validityDays := cert.ValidityDays
	if cert.Rule != nil && cert.Rule.ValidityDays > 0 {
		validityDays = cert.Rule.ValidityDays
	}

	now := time.Now()
	expiry := now.AddDate(0, 0, validityDays)

	assignmentID := req.TrainingAssignmentID

	issue := models.CertificateIssue{
		CertificationID:      req.CertificationID,
		UserID:               req.UserID,
		TrainingAssignmentID: &assignmentID,
		IssuedOn:             now,
		ExpiryDate:           expiry,
		CertificateNumber:    generateCertificateNumber(req.UserID, req.CertificationID),
		IssueStatus:          "issued",
	}

	created, err := s.certificateIssueRepo.Create(&issue)
	if err != nil {
		return nil, err
	}

	pdfPath, err := s.pdfService.Generate(created)
	if err != nil {
		return nil, err
	}

	created.PDFPath = pdfPath

	return s.certificateIssueRepo.Update(created)
}

func (s *CertificateIssueService) FindByID(id uint) (*models.CertificateIssue, error) {
	return s.certificateIssueRepo.FindByID(id)
}

func (s *CertificateIssueService) FindByUserID(userID uint) ([]models.CertificateIssue, error) {
	return s.certificateIssueRepo.FindByUserID(userID)
}

func generateCertificateNumber(userID uint, certificationID uint) string {
	return fmt.Sprintf("CERT-%d-%d-%d", userID, certificationID, time.Now().UnixNano())
}

func (s *CertificateIssueService) VerifyCertificate(certificateNumber string) (*models.CertificateIssue, error) {
	if certificateNumber == "" {
		return nil, errors.New("certificate number is required")
	}

	issue, err := s.certificateIssueRepo.FindByCertificateNumber(certificateNumber)
	if err != nil {
		return nil, err
	}

	if issue == nil {
		return nil, errors.New("certificate not found")
	}

	if issue.IssueStatus != "issued" {
		return nil, errors.New("certificate is not active")
	}

	return issue, nil
}
