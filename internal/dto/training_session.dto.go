package dto

type CreateTrainingSessionRequest struct {
	CourseID        uint  `json:"course_id" binding:"required"`
	ModuleID        *uint `json:"module_id"`
	CreatedByUserID uint  `json:"created_by_user_id" binding:"required"`

	SessionTitle string `json:"session_title" binding:"required"`
	SessionType  string `json:"session_type" binding:"required"`

	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`

	Location    string `json:"location"`
	MeetingLink string `json:"meeting_link"`

	IsMandatory *bool `json:"is_mandatory"`
}

type UpdateTrainingSessionRequest struct {
	CourseID *uint `json:"course_id"`
	ModuleID *uint `json:"module_id"`

	SessionTitle *string `json:"session_title"`
	SessionType  *string `json:"session_type"`

	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`

	Location    *string `json:"location"`
	MeetingLink *string `json:"meeting_link"`

	IsMandatory *bool   `json:"is_mandatory"`
	Status      *string `json:"status"`
}
