package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"example.com/m/internal/models"
	"example.com/m/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService  *services.AuthService
	tokenService *services.TokenService
}

func NewAuthHandler(authService *services.AuthService, tokenService *services.TokenService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
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
	setAuthCookies(c, tokens.AccessToken, tokens.RefreshToken)

	// oneHour := 60 * 60
	// oneYear := 365 * 24 * 60 * 60

	// secure := os.Getenv("GIN_MODE") == "release"
	// httpOnly := true
	// c.SetSameSite(http.SameSiteNoneMode)

	// c.SetCookie(
	// 	"access_token",
	// 	tokens.AccessToken,
	// 	oneHour,
	// 	"/",
	// 	"",
	// 	secure,
	// 	httpOnly,
	// )

	// c.SetCookie(
	// 	"refresh_token",
	// 	tokens.RefreshToken,
	// 	oneYear,
	// 	"/",
	// 	"",
	// 	secure,
	// 	httpOnly,
	// )

	c.JSON(http.StatusCreated, gin.H{
		"message": "register successful",
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
	setAuthCookies(c, tokens.AccessToken, tokens.RefreshToken)

	// oneHour := 60
	// oneYear := 365 * 24 * 60 * 60

	// secure := os.Getenv("GIN_MODE") == "release"
	// httpOnly := true

	// c.SetSameSite(http.SameSiteNoneMode)

	// c.SetCookie(
	// 	"access_token",
	// 	tokens.AccessToken,
	// 	oneHour,
	// 	"/",
	// 	"",
	// 	secure,
	// 	httpOnly,
	// )

	// c.SetCookie(
	// 	"refresh_token",
	// 	tokens.RefreshToken,
	// 	oneYear,
	// 	"/",
	// 	"",
	// 	secure,
	// 	httpOnly,
	// )

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    existingUser,
	})
}

func (h *AuthHandler) UsersList(c *gin.Context) {
	departmentIDQuery := c.Query("department_id")
	var departmentID uint

	if departmentIDQuery != "" {
		parsedDepartmentID, err := strconv.ParseUint(departmentIDQuery, 10, 64)
		if err != nil || parsedDepartmentID == 0 {
			c.JSON(400, gin.H{
				"error": "department_id must be a valid number greater than 0",
			})
			return
		}

		departmentID = uint(parsedDepartmentID)
	}
	users, err := h.authService.UsersList(departmentID)
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
	var req models.UpdateUserRequest
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

func (h *AuthHandler) Self(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user id missing from context",
		})
		return
	}

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "email missing from context",
		})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "role missing from context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    userID,
			"email": email,
			"role":  role,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token missing",
		})
		return
	}

	tokens, user, err := h.authService.RefreshTokens(refreshToken)
	if err != nil {
		clearAuthCookies(c)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired refresh token",
		})
		return
	}
	log.Println("refresh hit")
	log.Println("new access token empty?", tokens.AccessToken == "")
	log.Println("new refresh token empty?", tokens.RefreshToken == "")
	setAuthCookies(c, tokens.AccessToken, tokens.RefreshToken)

	// accessTokenAge := 60 * 60
	// refreshTokenAge := 365 * 24 * 60 * 60

	// secure := os.Getenv("GIN_MODE") == "release"
	// httpOnly := true

	// c.SetSameSite(http.SameSiteNoneMode)

	// c.SetCookie(
	// 	"access_token",
	// 	tokens.AccessToken,
	// 	accessTokenAge,
	// 	"/",
	// 	"",
	// 	secure,
	// 	httpOnly,
	// )

	// c.SetCookie(
	// 	"refresh_token",
	// 	tokens.RefreshToken,
	// 	refreshTokenAge,
	// 	"/",
	// 	"",
	// 	secure,
	// 	httpOnly,
	// )

	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		tokenHash := services.HashToken(refreshToken)
		_ = h.authService.RevokeRefreshToken(tokenHash)
	}

	clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

func setAuthCookies(c *gin.Context, accessToken string, refreshToken string) {
	isProd := os.Getenv("GIN_MODE") == "release"
	httpOnly := true

	accessTokenAge := 60 * 60
	refreshTokenAge := 365 * 24 * 60 * 60

	if isProd {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie(
		"access_token",
		accessToken,
		accessTokenAge,
		"/",
		"",
		isProd,
		httpOnly,
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		refreshTokenAge,
		"/",
		"",
		isProd,
		httpOnly,
	)
}

func clearAuthCookies(c *gin.Context) {
	isProd := os.Getenv("GIN_MODE") == "release"
	httpOnly := true

	if isProd {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie("access_token", "", -1, "/", "", isProd, httpOnly)
	c.SetCookie("refresh_token", "", -1, "/", "", isProd, httpOnly)
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to bind request",
			"details": err.Error(),
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

	user, err := h.authService.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "user creation failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user":    user,
	})
}

func (h *AuthHandler) CreateBulkUsers(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "excel file is required",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to open file",
		})
		return
	}
	defer file.Close()

	result, err := h.authService.CreateBulkUsers(
		c.Request.Context(),
		file,
		fileHeader,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bulk users upload processed",
		"data":    result,
	})
}
