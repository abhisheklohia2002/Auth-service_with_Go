package dto

type CreateCourseRequest struct {
	CourseTitle          string `form:"course_title" json:"course_title" binding:"required"`
	CourseDescription    string `form:"course_description" json:"course_description"`
	CourseType           string `form:"course_type" json:"course_type" binding:"required"`
	TotalDurationMinutes uint   `form:"total_duration_minutes" json:"total_duration_minutes" binding:"required"`
}

type UpdateCourseRequest struct {
	CourseTitle          string `form:"course_title" json:"course_title"`
	CourseDescription    string `form:"course_description" json:"course_description"`
	CourseType           string `form:"course_type" json:"course_type"`
	IsActive             *bool  `form:"is_active" json:"is_active"`
	TotalDurationMinutes *uint  `form:"total_duration_minutes" json:"total_duration_minutes"`
}
