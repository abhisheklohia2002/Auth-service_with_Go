package services_attendence

import (
	"errors"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	repositories_attendance "example.com/m/internal/repositories/attendance"
	repositories_trainingsession "example.com/m/internal/repositories/training_session"
)

type AttendanceService interface {
	Mark(sessionID uint, req dto.MarkAttendanceRequest) (*models.Attendance, error)
	BulkMark(sessionID uint, req dto.BulkMarkAttendanceRequest) ([]models.Attendance, error)
	GetBySession(sessionID uint) ([]models.Attendance, error)
	GetByUser(userID uint) ([]models.Attendance, error)
	Update(id uint, req dto.UpdateAttendanceRequest) (*models.Attendance, error)
}

type attendanceService struct {
	attendanceRepo      repositories_attendance.AttendanceRepository
	trainingSessionRepo repositories_trainingsession.TrainingSessionRepository
}

func NewAttendanceService(
	attendanceRepo repositories_attendance.AttendanceRepository,
	trainingSessionRepo repositories_trainingsession.TrainingSessionRepository,
) AttendanceService {
	return &attendanceService{
		attendanceRepo:      attendanceRepo,
		trainingSessionRepo: trainingSessionRepo,
	}
}

var allowedAttendanceStatuses = map[string]bool{
	"present": true,
	"absent":  true,
	"late":    true,
	"excused": true,
}

var allowedAttendanceSources = map[string]bool{
	"manual": true,
	"qr":     true,
	"geo":    true,
	"auto":   true,
}

func parseAttendanceDateTime(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func calculateDurationMinutes(checkIn *time.Time, checkOut *time.Time) int {
	if checkIn == nil || checkOut == nil {
		return 0
	}

	if !checkOut.After(*checkIn) {
		return 0
	}

	return int(checkOut.Sub(*checkIn).Minutes())
}

func (s *attendanceService) validateSession(sessionID uint) error {
	session, err := s.trainingSessionRepo.FindByID(sessionID)
	if err != nil {
		return errors.New("training session not found")
	}

	if session.Status == "cancelled" {
		return errors.New("cannot mark attendance for cancelled session")
	}

	return nil
}

func (s *attendanceService) Mark(sessionID uint, req dto.MarkAttendanceRequest) (*models.Attendance, error) {
	if err := s.validateSession(sessionID); err != nil {
		return nil, err
	}

	if !allowedAttendanceStatuses[req.Status] {
		return nil, errors.New("invalid attendance status")
	}

	source := req.AttendanceSource
	if source == "" {
		source = "manual"
	}

	if !allowedAttendanceSources[source] {
		return nil, errors.New("invalid attendance source")
	}

	existing, err := s.attendanceRepo.FindBySessionAndUser(sessionID, req.UserID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("attendance already marked for this user in this session")
	}

	checkInTime, err := parseAttendanceDateTime(req.CheckInTime)
	if err != nil {
		return nil, errors.New("invalid check_in_time, use RFC3339 format")
	}

	checkOutTime, err := parseAttendanceDateTime(req.CheckOutTime)
	if err != nil {
		return nil, errors.New("invalid check_out_time, use RFC3339 format")
	}

	if checkInTime != nil && checkOutTime != nil && !checkOutTime.After(*checkInTime) {
		return nil, errors.New("check_out_time must be after check_in_time")
	}

	attendance := &models.Attendance{
		SessionID:        sessionID,
		UserID:           req.UserID,
		MarkedByUserID:   req.MarkedByUserID,
		Status:           req.Status,
		CheckInTime:      checkInTime,
		CheckOutTime:     checkOutTime,
		DurationMinutes:  calculateDurationMinutes(checkInTime, checkOutTime),
		AttendanceSource: source,
		Remarks:          req.Remarks,
	}

	if err := s.attendanceRepo.Create(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (s *attendanceService) BulkMark(sessionID uint, req dto.BulkMarkAttendanceRequest) ([]models.Attendance, error) {
	if err := s.validateSession(sessionID); err != nil {
		return nil, err
	}

	if len(req.Attendances) == 0 {
		return nil, errors.New("attendances cannot be empty")
	}

	seenUsers := make(map[uint]bool)
	attendances := make([]models.Attendance, 0)

	for _, item := range req.Attendances {
		if seenUsers[item.UserID] {
			return nil, errors.New("duplicate user_id found in request")
		}
		seenUsers[item.UserID] = true

		if !allowedAttendanceStatuses[item.Status] {
			return nil, errors.New("invalid attendance status")
		}

		source := item.AttendanceSource
		if source == "" {
			source = "manual"
		}

		if !allowedAttendanceSources[source] {
			return nil, errors.New("invalid attendance source")
		}

		existing, err := s.attendanceRepo.FindBySessionAndUser(sessionID, item.UserID)
		if err != nil {
			return nil, err
		}

		if existing != nil {
			return nil, errors.New("attendance already marked for one or more users")
		}

		checkInTime, err := parseAttendanceDateTime(item.CheckInTime)
		if err != nil {
			return nil, errors.New("invalid check_in_time, use RFC3339 format")
		}

		checkOutTime, err := parseAttendanceDateTime(item.CheckOutTime)
		if err != nil {
			return nil, errors.New("invalid check_out_time, use RFC3339 format")
		}

		if checkInTime != nil && checkOutTime != nil && !checkOutTime.After(*checkInTime) {
			return nil, errors.New("check_out_time must be after check_in_time")
		}

		attendances = append(attendances, models.Attendance{
			SessionID:        sessionID,
			UserID:           item.UserID,
			MarkedByUserID:   item.MarkedByUserID,
			Status:           item.Status,
			CheckInTime:      checkInTime,
			CheckOutTime:     checkOutTime,
			DurationMinutes:  calculateDurationMinutes(checkInTime, checkOutTime),
			AttendanceSource: source,
			Remarks:          item.Remarks,
		})
	}

	if err := s.attendanceRepo.BulkCreate(attendances); err != nil {
		return nil, err
	}

	return attendances, nil
}

func (s *attendanceService) GetBySession(sessionID uint) ([]models.Attendance, error) {
	if err := s.validateSession(sessionID); err != nil {
		return nil, err
	}

	return s.attendanceRepo.FindBySessionID(sessionID)
}

func (s *attendanceService) GetByUser(userID uint) ([]models.Attendance, error) {
	return s.attendanceRepo.FindByUserID(userID)
}

func (s *attendanceService) Update(id uint, req dto.UpdateAttendanceRequest) (*models.Attendance, error) {
	attendance, err := s.attendanceRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("attendance not found")
	}

	if req.Status != nil {
		if !allowedAttendanceStatuses[*req.Status] {
			return nil, errors.New("invalid attendance status")
		}
		attendance.Status = *req.Status
	}

	if req.AttendanceSource != nil {
		if !allowedAttendanceSources[*req.AttendanceSource] {
			return nil, errors.New("invalid attendance source")
		}
		attendance.AttendanceSource = *req.AttendanceSource
	}

	if req.CheckInTime != nil {
		checkInTime, err := parseAttendanceDateTime(req.CheckInTime)
		if err != nil {
			return nil, errors.New("invalid check_in_time, use RFC3339 format")
		}
		attendance.CheckInTime = checkInTime
	}

	if req.CheckOutTime != nil {
		checkOutTime, err := parseAttendanceDateTime(req.CheckOutTime)
		if err != nil {
			return nil, errors.New("invalid check_out_time, use RFC3339 format")
		}
		attendance.CheckOutTime = checkOutTime
	}

	if attendance.CheckInTime != nil && attendance.CheckOutTime != nil && !attendance.CheckOutTime.After(*attendance.CheckInTime) {
		return nil, errors.New("check_out_time must be after check_in_time")
	}

	attendance.DurationMinutes = calculateDurationMinutes(attendance.CheckInTime, attendance.CheckOutTime)

	if req.Remarks != nil {
		attendance.Remarks = *req.Remarks
	}

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}
