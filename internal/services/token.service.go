package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"example.com/m/internal/common"
	"example.com/m/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// const jwtKeyID = "access-key-1"

type TokenService struct {
	config     config.Config
	privateKey *rsa.PrivateKey
	jwksURL    string
	issuer     string
}

type TokenClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenService(cfg config.Config) (*TokenService, error) {
	privateKey, err := common.LoadRSAPrivateKeyFromEnv("JWT_PRIVATE_KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	return &TokenService{
		config:     cfg,
		privateKey: privateKey,
	}, nil
}

func (s *TokenService) GenerateAccessToken(userID uint, email string, role string) (string, error) {
	now := time.Now()
	expiry := now.Add(time.Duration(s.config.AccessTokenExpiryMins) * time.Minute)

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.JWTIssuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{"access"},
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.config.JWTKEYID

	return token.SignedString(s.privateKey)
}

func (s *TokenService) GenerateRefreshToken(userID uint, email string, role string) (string, error) {
	now := time.Now()
	expiry := now.Add(time.Duration(s.config.RefreshTokenExpiryHours) * time.Hour)

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.JWTIssuer,
			Subject:   fmt.Sprintf("%d", userID),
			Audience:  []string{"refresh"},
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.config.JWTKEYID

	return token.SignedString(s.privateKey)
}

func LoadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)

	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		rsaPrivateKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}

		return rsaPrivateKey, nil

	default:
		return nil, errors.New("unsupported private key type: " + block.Type)
	}
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Make sure token uses HMAC signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.LoadDotenv().REFRESH_TOKEN_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
