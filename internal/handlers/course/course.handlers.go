package handlers_course

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	services_course "example.com/m/internal/services/course"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService *services_course.CourseService
}

func NewCourseHandler(courseService *services_course.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req dto.CreateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id missing from context",
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id type in context",
		})
		return
	}

	course, err := h.courseService.CreateCourse(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "course created successfully",
		"data":    course,
	})
}

func (h *CourseHandler) GetCourses(c *gin.Context) {
	courses, err := h.courseService.GetCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch courses",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "courses fetched successfully",
		"data":    courses,
	})
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	course, err := h.courseService.GetCourseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "course not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "course fetched successfully",
		"data":    course,
	})
}

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var req dto.UpdateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	course, err := h.courseService.UpdateCourse(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "course updated successfully",
		"data":    course,
	})
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	if err := h.courseService.DeleteCourse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "course deleted successfully",
	})
}

func parseUintParam(c *gin.Context, key string) (uint, error) {
	value := c.Param(key)

	num, err := strconv.Atoi(value)
	if err != nil || num <= 0 {
		return 0, err
	}

	return uint(num), nil
}
