package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/m/internal/models"
	"example.com/m/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to bind request",
		})
		return
	}
	isExist, _, err := h.authService.UserExistsByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email already exists",
		})
		return
	}
	tokens, user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "registration failed",
		})
		return
	}
	oneHour := 60 * 60
	oneYear := 365 * 24 * 60 * 60

	secure := false
	httpOnly := true

	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		oneHour,
		"/",
		"",
		secure,
		httpOnly,
	)

	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		oneYear,
		"/",
		"",
		secure,
		httpOnly,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Register successful",
		"user":    user,
	})

}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	isExist, existingUser, err := h.authService.UserExistsByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	tokens, err := h.authService.Login(existingUser, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	oneHour := 60 * 60
	oneYear := 365 * 24 * 60 * 60

	secure := false
	httpOnly := true

	c.SetCookie(
		"access_token",
		tokens.AccessToken,
		oneHour,
		"/",
		"",
		secure,
		httpOnly,
	)

	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		oneYear,
		"/",
		"",
		secure,
		httpOnly,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    existingUser,
	})
}

func (h *AuthHandler) UsersList(c *gin.Context) {
	users, err := h.authService.UsersList()
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch users",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "users fetched successfully",
		"data":    users,
	})
}

func (h *AuthHandler) DeleteUserById(c *gin.Context) {
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	_, err := h.authService.DeleteUserById(num)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch users",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("User Remove id %s", id),
	})
}

func (h *AuthHandler) UpdateUserById(c *gin.Context) {
	id := c.Param("id")
	num, _ := strconv.Atoi(id)
	var req models.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to bind request",
		})
		return
	}

	user, err := h.authService.UpdateUserById(c.Request.Context(), uint(num), req)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to fetch users",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("User update by id %s", id),
		"user":    user,
	})

}
