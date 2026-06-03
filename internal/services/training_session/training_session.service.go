package services_trainingsession

import (
	"errors"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_trainingsession "example.com/m/internal/repositories/training_session"
)

type TrainingSessionService interface {
	Create(req dto.CreateTrainingSessionRequest) (*models.TrainingSession, error)
	GetAll(filter models.TrainingSessionFilter) ([]models.TrainingSession, error)
	GetByID(id uint) (*models.TrainingSession, error)
	Update(id uint, req dto.UpdateTrainingSessionRequest) (*models.TrainingSession, error)
	Delete(id uint) error
}

type trainingSessionService struct {
	repo repositories_trainingsession.TrainingSessionRepository
}

func NewTrainingSessionService(repo repositories_trainingsession.TrainingSessionRepository) TrainingSessionService {
	return &trainingSessionService{repo: repo}
}

var allowedSessionTypes = map[string]bool{
	"online":     true,
	"offline":    true,
	"hybrid":     true,
	"self_paced": true,
}

var allowedSessionStatuses = map[string]bool{
	"scheduled": true,
	"completed": true,
	"cancelled": true,
}

func parseDateTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func (s *trainingSessionService) Create(req dto.CreateTrainingSessionRequest) (*models.TrainingSession, error) {
	if !allowedSessionTypes[req.SessionType] {
		return nil, errors.New("invalid session_type")
	}

	startTime, err := parseDateTime(req.StartTime)
	if err != nil {
		return nil, errors.New("invalid start_time, use RFC3339 format")
	}

	endTime, err := parseDateTime(req.EndTime)
	if err != nil {
		return nil, errors.New("invalid end_time, use RFC3339 format")
	}

	if !endTime.After(startTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	isMandatory := true
	if req.IsMandatory != nil {
		isMandatory = *req.IsMandatory
	}

	session := &models.TrainingSession{
		CourseID:        req.CourseID,
		ModuleID:        req.ModuleID,
		CreatedByUserID: req.CreatedByUserID,
		SessionTitle:    req.SessionTitle,
		SessionType:     req.SessionType,
		StartTime:       startTime,
		EndTime:         endTime,
		Location:        req.Location,
		MeetingLink:     req.MeetingLink,
		IsMandatory:     isMandatory,
		Status:          "scheduled",
	}

	if err := s.repo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *trainingSessionService) GetAll(filter models.TrainingSessionFilter) ([]models.TrainingSession, error) {
	return s.repo.FindAll(filter)
}

func (s *trainingSessionService) GetByID(id uint) (*models.TrainingSession, error) {
	return s.repo.FindByID(id)
}

func (s *trainingSessionService) Update(id uint, req dto.UpdateTrainingSessionRequest) (*models.TrainingSession, error) {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("training session not found")
	}

	if session.Status == "cancelled" {
		return nil, errors.New("cancelled session cannot be updated")
	}

	if req.CourseID != nil {
		session.CourseID = *req.CourseID
	}

	if req.ModuleID != nil {
		session.ModuleID = req.ModuleID
	}

	if req.SessionTitle != nil {
		session.SessionTitle = *req.SessionTitle
	}

	if req.SessionType != nil {
		if !allowedSessionTypes[*req.SessionType] {
			return nil, errors.New("invalid session_type")
		}
		session.SessionType = *req.SessionType
	}

	if req.StartTime != nil {
		startTime, err := parseDateTime(*req.StartTime)
		if err != nil {
			return nil, errors.New("invalid start_time, use RFC3339 format")
		}
		session.StartTime = startTime
	}

	if req.EndTime != nil {
		endTime, err := parseDateTime(*req.EndTime)
		if err != nil {
			return nil, errors.New("invalid end_time, use RFC3339 format")
		}
		session.EndTime = endTime
	}

	if !session.EndTime.After(session.StartTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	if req.Location != nil {
		session.Location = *req.Location
	}

	if req.MeetingLink != nil {
		session.MeetingLink = *req.MeetingLink
	}

	if req.IsMandatory != nil {
		session.IsMandatory = *req.IsMandatory
	}

	if req.Status != nil {
		if !allowedSessionStatuses[*req.Status] {
			return nil, errors.New("invalid status")
		}
		session.Status = *req.Status
	}

	if err := s.repo.Update(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *trainingSessionService) Delete(id uint) error {
	session, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("training session not found")
	}

	return s.repo.Delete(session)
}