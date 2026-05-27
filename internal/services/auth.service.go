package services

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"example.com/m/internal/config"
	"example.com/m/internal/helper"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *repositories.UserRepository
	tokenService *TokenService
}
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}
type JWKS struct {
	Keys []JWK `json:"keys"`
}

func (s *AuthService) UsersList() ([]models.User, error) {
	return s.userRepo.UsersList()
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *AuthService) UserExistsByEmail(email string) (bool, *models.User, error) {
	return s.userRepo.ExistsByEmail(email)
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) Register(req models.RegisterRequest) (models.AuthResponse, models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	createdUser, err := s.userRepo.Create(&user)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	tokens, err := s.generateTokens(&createdUser)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	now := time.Now()
	expiry := now.Add(time.Duration(config.LoadDotenv().RefreshTokenExpiryHours) * time.Hour)

	refreshTokenHash := HashToken(tokens.RefreshToken)

	_, err = s.userRepo.PersistRefreshToken(createdUser.ID, refreshTokenHash, expiry)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	return tokens, createdUser, nil
}

func (s *AuthService) Login(user *models.User, password string) (models.AuthResponse, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	tokens, err := s.generateTokens(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	now := time.Now()
	expiry := now.Add(time.Duration(config.LoadDotenv().RefreshTokenExpiryHours) * time.Hour)

	refreshTokenHash := HashToken(tokens.RefreshToken)

	_, err = s.userRepo.PersistRefreshToken(user.ID, refreshTokenHash, expiry)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return tokens, nil
}

func (s *AuthService) Userslist() ([]models.User, error) {
	return s.userRepo.UsersList()
}
func (s *AuthService) DeleteUserById(id int) (models.User, error) {
	return s.userRepo.DeleteUserById(id)
}

func (s *AuthService) UpdateUserById(c context.Context, id uint, req models.RegisterRequest) (models.User, error) {
	user, err := s.userRepo.UpdateUserById(c, id, req)
	if err != nil {
		return models.User{}, err
	}
	if user == nil {
		return models.User{}, fmt.Errorf("user not found")
	}
	return *user, nil
}

func (s *AuthService) generateTokens(user *models.User) (models.AuthResponse, error) {
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.AuthResponse{}, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthService) ValidateAccessToken(tokenString string) (*models.User, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, errors.New("missing kid in token header")
		}
		publicKey, err := s.tokenService.getPublicKeyFromJWKS(kid)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if !helper.HasAudience(claims.Audience, "access") {
		return nil, errors.New("invalid audience")
	}
	_, user, err := s.userRepo.ExistsByEmail(claims.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *TokenService) getPublicKeyFromJWKS(kid string) (*rsa.PublicKey, error) {

	resp, err := http.Get(config.LoadDotenv().JWTIssuer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {

			if key.Kty != "RSA" {
				return nil, errors.New("invalid key type")
			}

			if key.Alg != "RS256" {
				return nil, errors.New("invalid key algorithm")
			}

			return jwkToRSAPublicKey(key)
		}
	}

	return nil, errors.New("matching public key not found")
}

func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)

	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	if e == 0 {
		return nil, errors.New("invalid exponent")
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}
