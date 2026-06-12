package handlers_entity

import (
	"net/http"
	"strconv"

	"example.com/m/internal/dto"
	services_entity "example.com/m/internal/services/entity"
	"github.com/gin-gonic/gin"
)

type EntityHandler struct {
	entityService services_entity.EntityService
}

func NewEntityHandler(entityService services_entity.EntityService) *EntityHandler {
	return &EntityHandler{entityService: entityService}
}

func (h *EntityHandler) CreateEntity(c *gin.Context) {
	var req dto.CreateEntityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	entity, err := h.entityService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "entity created successfully",
		"data":    entity,
	})
}

func (h *EntityHandler) GetEntities(c *gin.Context) {
	entities, err := h.entityService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch entities",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    entities,
	})
}

func (h *EntityHandler) GetEntityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid entity id",
		})
		return
	}

	entity, err := h.entityService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    entity,
	})
}

func (h *EntityHandler) UpdateEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid entity id",
		})
		return
	}

	var req dto.UpdateEntityRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	entity, err := h.entityService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "entity updated successfully",
		"data":    entity,
	})
}

func (h *EntityHandler) DeleteEntity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid entity id",
		})
		return
	}

	if err := h.entityService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "entity deleted successfully",
	})
}