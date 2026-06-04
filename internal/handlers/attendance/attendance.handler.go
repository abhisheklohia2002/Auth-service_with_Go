package handlers_attendance

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	services_attendence "example.com/m/internal/services/attendence"
	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service services_attendence.AttendanceService
}

func NewAttendanceHandler(service services_attendence.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

func (h *AttendanceHandler) Mark(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid session id",
		})
		return
	}

	var req dto.MarkAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	attendance, err := h.service.Mark(uint(sessionID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Attendance marked successfully",
		"data":    attendance,
	})
}

func (h *AttendanceHandler) BulkMark(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid session id",
		})
		return
	}

	var req dto.BulkMarkAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	attendances, err := h.service.BulkMark(uint(sessionID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bulk attendance marked successfully",
		"data":    attendances,
	})
}

func (h *AttendanceHandler) GetBySession(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid session id",
		})
		return
	}

	attendances, err := h.service.GetBySession(uint(sessionID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session attendance fetched successfully",
		"data":    attendances,
	})
}

func (h *AttendanceHandler) GetByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user id",
		})
		return
	}

	attendances, err := h.service.GetByUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user attendance",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User attendance fetched successfully",
		"data":    attendances,
	})
}

func (h *AttendanceHandler) Update(c *gin.Context) {
	attendanceID, err := strconv.Atoi(c.Param("attendanceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid attendance id",
		})
		return
	}

	var req dto.UpdateAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	attendance, err := h.service.Update(uint(attendanceID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance updated successfully",
		"data":    attendance,
	})
}