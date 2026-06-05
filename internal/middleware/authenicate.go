package middleware

import (
	"crypto/rsa"
	"errors"
	"log"
	"net/http"

	"example.com/m/internal/common"
	"example.com/m/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	config    config.Config
	publicKey *rsa.PublicKey
}

func NewAuthMiddleware(cfg config.Config) (*AuthMiddleware, error) {
	// raw := os.Getenv("JWT_PUBLIC_KEY")
	publicKey, err := common.LoadRSAPublicKeyFromEnv("JWT_PUBLIC_KEY")
	if err != nil {
		return nil, err
	}

	return &AuthMiddleware{
		config:    cfg,
		publicKey: publicKey,
	}, nil
}

func (m *AuthMiddleware) IsAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "access token cookie is required",
			})
			c.Abort()
			return
		}

		claims := &common.Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
					return nil, errors.New("invalid signing algorithm")
				}

				return m.publicKey, nil
			},
			jwt.WithIssuer(m.config.JWTIssuer),
			jwt.WithAudience("access"),
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired access token",
			})
			c.Abort()
			return
		}

		log.Println("JWT role:", claims.Role)
		log.Println("Allowed roles:", allowedRoles)
		if len(allowedRoles) > 0 && !isRoleAllowed(claims.Role, allowedRoles) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "you do not have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func isRoleAllowed(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if userRole == role {
			return true
		}
	}

	return false
}
