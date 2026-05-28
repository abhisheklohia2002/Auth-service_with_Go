package dto

type CreateCourseRequest struct {
	CourseTitle       string `json:"course_title" binding:"required"`
	CourseDescription string `json:"course_description"`
	CourseType        string `json:"course_type" binding:"required"`
}
type UpdateCourseRequest struct {
	CourseTitle       string `json:"course_title"`
	CourseDescription string `json:"course_description"`
	CourseType        string `json:"course_type"`
	IsActive          *bool  `json:"is_active"`
}
