package handlers_certificationrule

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	services_certificationrule "example.com/m/internal/services/certification_rule"
	"github.com/gin-gonic/gin"
)

type CertificationRuleHandler struct {
	service *services_certificationrule.CertificationRuleService
}

func NewCertificationRuleHandler(service *services_certificationrule.CertificationRuleService) *CertificationRuleHandler {
	return &CertificationRuleHandler{service: service}
}

func (h *CertificationRuleHandler) Create(c *gin.Context) {
	var req dto.CreateCertificationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	rule, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "certification rule created successfully", "data": rule})
}

func (h *CertificationRuleHandler) FindAll(c *gin.Context) {
	rules, err := h.service.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certification rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification rules fetched successfully", "data": rules})
}

func (h *CertificationRuleHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	rule, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch certification rule"})
		return
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "certification rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification rule fetched successfully", "data": rule})
}

func (h *CertificationRuleHandler) Update(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateCertificationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind request", "details": err.Error()})
		return
	}

	rule, err := h.service.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification rule updated successfully", "data": rule})
}

func (h *CertificationRuleHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "certification rule deleted successfully"})
}

func (h *CertificationRuleHandler) CreateByCourseID(c *gin.Context) {
	courseID, ok := helper.ParseUintParam(c, "courseId")
	if !ok {
		return
	}

	var req dto.CreateCertificationRuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	rule, err := h.service.CreateByCourseID(courseID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "certification rule created successfully",
		"data":    rule,
	})
}