package handlers

import (
	"example.com/m/internal/models"
	"example.com/m/internal/services"
	"github.com/gin-gonic/gin"
)

type CompanyHandlers struct {
	CompanyService *services.CompanyService
}

func NewCompanyHandler(CompanyService *services.CompanyService) *CompanyHandlers {
	return &CompanyHandlers{CompanyService: CompanyService}
}

func (h *CompanyHandlers) CompanyCreate(c *gin.Context) {
	var req models.CreateCompanyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	company, err := h.CompanyService.CompanyCreate(req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "company created successfully",
		"data":    company,
	})
}
