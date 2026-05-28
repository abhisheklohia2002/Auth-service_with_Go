package handlers_certification

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_certification "example.com/m/internal/services/certification"
	"github.com/gin-gonic/gin"
)

type CertificationHandler struct {
	service *services_certification.CertificationService
}

func NewCertificationHandler(service *services_certification.CertificationService) *CertificationHandler {
	return &CertificationHandler{service: service}
}

func (h *CertificationHandler) Create(c *gin.Context) {
	var req dto.CreateCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	cert, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "certification created successfully", "data": cert})
}

func (h *CertificationHandler) FindAll(c *gin.Context) {
	certs, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certifications fetched successfully", "data": certs})
}

func (h *CertificationHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	cert, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certification"})
		return
	}
	if cert == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification fetched successfully", "data": cert})
}

func (h *CertificationHandler) FindByCourseID(c *gin.Context) {
	courseID, ok := helper.ParseUintParam(c, "courseId")
	if !ok {
		return
	}

	certs, err := h.service.FindByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certifications fetched successfully", "data": certs})
}

func (h *CertificationHandler) Update(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	cert, err := h.service.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification updated successfully", "data": cert})
}

func (h *CertificationHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification deleted successfully"})
}
