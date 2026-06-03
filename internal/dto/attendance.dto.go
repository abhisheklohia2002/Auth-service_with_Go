package dto

type MarkAttendanceRequest struct {
	UserID           uint    `json:"user_id" binding:"required"`
	MarkedByUserID   *uint   `json:"marked_by_user_id"`
	Status           string  `json:"status" binding:"required"`
	CheckInTime      *string `json:"check_in_time"`
	CheckOutTime     *string `json:"check_out_time"`
	AttendanceSource string  `json:"attendance_source"`
	Remarks          string  `json:"remarks"`
}

type BulkMarkAttendanceRequest struct {
	Attendances []MarkAttendanceRequest `json:"attendances" binding:"required"`
}

type UpdateAttendanceRequest struct {
	Status           *string `json:"status"`
	CheckInTime      *string `json:"check_in_time"`
	CheckOutTime     *string `json:"check_out_time"`
	AttendanceSource *string `json:"attendance_source"`
	Remarks          *string `json:"remarks"`
}

type AttendanceFilter struct {
	SessionID string
	UserID    string
	Status    string
}
