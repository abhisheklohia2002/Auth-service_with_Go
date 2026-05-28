package handlers_certificateIssue

import (
	"net/http"
	"time"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"

	services_certificateissue "example.com/m/internal/services/certificate_issue"
	"github.com/gin-gonic/gin"
)

type CertificateIssueHandler struct {
	service *services_certificateissue.CertificateIssueService
}

func NewCertificateIssueHandler(service *services_certificateissue.CertificateIssueService) *CertificateIssueHandler {
	return &CertificateIssueHandler{service: service}
}

func (h *CertificateIssueHandler) Issue(c *gin.Context) {
	var req dto.IssueCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	issue, err := h.service.Issue(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "certificate issued successfully",
		"data":    issue,
	})
}

func (h *CertificateIssueHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	issue, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certificate issue"})
		return
	}
	if issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate issue not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certificate issue fetched successfully", "data": issue})
}

func (h *CertificateIssueHandler) FindByUserID(c *gin.Context) {
	userID, ok := helper.ParseUintParam(c, "userId")
	if !ok {
		return
	}

	issues, err := h.service.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certificate issues fetched successfully", "data": issues})
}

func (h *CertificateIssueHandler) DownloadPDF(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	issue, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certificate issue"})
		return
	}
	if issue == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate issue not found"})
		return
	}

	if issue.PDFPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "certificate pdf not found"})
		return
	}

	c.FileAttachment(issue.PDFPath, "certificate.pdf")
}

func (h *CertificateIssueHandler) VerifyCertificate(c *gin.Context) {
	certificateNumber := c.Param("certificateNumber")

	issue, err := h.service.VerifyCertificate(certificateNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	isExpired := time.Now().After(issue.ExpiryDate)

	c.JSON(http.StatusOK, gin.H{
		"valid":      !isExpired,
		"is_expired": isExpired,
		"certificate": gin.H{
			"certificate_number": issue.CertificateNumber,
			"issue_status":       issue.IssueStatus,
			"issued_on":          issue.IssuedOn,
			"expiry_date":        issue.ExpiryDate,
			"user": gin.H{
				"id":    issue.User.ID,
				"name":  issue.User.FullName,
				"email": issue.User.Email,
			},
			"course": gin.H{
				"id":    issue.Certification.Course.ID,
				"title": issue.Certification.Course.CourseTitle,
			},
			"certification": gin.H{
				"id":   issue.Certification.ID,
				"name": issue.Certification.CertificationName,
			},
		},
	})
}
