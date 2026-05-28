package handlers_role

import (
	"net/http"

	"example.com/m/internal/dto"
	"example.com/m/internal/helper"

	services_role "example.com/m/internal/services/role"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *services_role.RoleService
}

func NewRoleHandler(roleService *services_role.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	role, err := h.roleService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "role created successfully",
		"data":    role,
	})
}

func (h *RoleHandler) FindAll(c *gin.Context) {
	roles, err := h.roleService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch roles",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "roles fetched successfully",
		"data":    roles,
	})
}

func (h *RoleHandler) FindByID(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	role, err := h.roleService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch role",
		})
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "role not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role fetched successfully",
		"data":    role,
	})
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
		})
		return
	}

	role, err := h.roleService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role updated successfully",
		"data":    role,
	})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, ok := helper.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.roleService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "role deleted successfully",
	})
}
