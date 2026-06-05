package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"example.com/m/internal/config"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *repositories.UserRepository
	tokenService *TokenService
	roleRepo     *repositories.RoleRepository
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
	roleRepo *repositories.RoleRepository,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
		roleRepo:     roleRepo,
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
		FullName:     req.Name,
		Email:        req.Email,
		Password:     string(hashedPassword),
		RoleID:       req.RoleID,
		EmployeeCode: req.EmployeeCode,
		ManagerID:    req.ManagerID,
		Status:       "active",
		JoiningDate:  time.Now(),
		DepartmentID: req.DepartmentID,
	}

	createdUser, err := s.userRepo.Create(&user)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	createdUserWithRole, err := s.userRepo.FindByIDWithRole(createdUser.ID)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	tokens, err := s.generateTokens(&createdUserWithRole)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	expiry := time.Now().Add(time.Duration(config.LoadDotenv().RefreshTokenExpiryHours) * time.Hour)
	refreshTokenHash := HashToken(tokens.RefreshToken)

	_, err = s.userRepo.PersistRefreshToken(createdUserWithRole.ID, refreshTokenHash, expiry)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	return tokens, createdUserWithRole, nil
}

func (s *AuthService) Login(user *models.User, password string) (models.AuthResponse, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.AuthResponse{}, errors.New("invalid email or password")
	}

	userWithRole, err := s.userRepo.FindByIDWithRole(user.ID)
	if err != nil {
		return models.AuthResponse{}, err
	}

	tokens, err := s.generateTokens(&userWithRole)
	if err != nil {
		return models.AuthResponse{}, err
	}

	expiry := time.Now().Add(time.Duration(config.LoadDotenv().RefreshTokenExpiryHours) * time.Hour)
	refreshTokenHash := HashToken(tokens.RefreshToken)

	_, err = s.userRepo.PersistRefreshToken(userWithRole.ID, refreshTokenHash, expiry)
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

func (s *AuthService) UpdateUserById(c context.Context, id uint, req models.UpdateUserRequest) (models.User, error) {
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
	roleName := user.Role.RoleName
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, roleName)
	if err != nil {
		return models.AuthResponse{}, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID, user.Email, roleName)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *AuthService) RefreshTokens(refreshToken string) (models.AuthResponse, models.User, error) {
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return models.AuthResponse{}, models.User{}, errors.New("invalid refresh token")
	}

	oldTokenHash := HashToken(refreshToken)

	storedToken, err := s.userRepo.FindValidRefreshToken(oldTokenHash)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	if storedToken == nil {
		return models.AuthResponse{}, models.User{}, errors.New("refresh token expired or revoked")
	}

	user, err := s.userRepo.FindByIDWithRole(claims.UserID)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	tokens, err := s.generateTokens(&user)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	if err := s.userRepo.RevokeRefreshToken(oldTokenHash); err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	expiry := time.Now().Add(time.Duration(config.LoadDotenv().RefreshTokenExpiryHours) * time.Hour)
	newRefreshTokenHash := HashToken(tokens.RefreshToken)

	_, err = s.userRepo.PersistRefreshToken(user.ID, newRefreshTokenHash, expiry)
	if err != nil {
		return models.AuthResponse{}, models.User{}, err
	}

	return tokens, user, nil
}

func (s *AuthService) RevokeRefreshToken(tokenHash string) error {
	return s.userRepo.RevokeRefreshToken(tokenHash)
}

func (s *AuthService) CreateUser(req models.RegisterRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	user := models.User{
		FullName:     req.Name,
		Email:        req.Email,
		Password:     string(hashedPassword),
		RoleID:       req.RoleID,
		EmployeeCode: req.EmployeeCode,
		Status:       req.Status,
		ManagerID:    req.ManagerID,
		DepartmentID: req.DepartmentID,
	}

	if user.Status == "" {
		user.Status = "active"
	}

	createdUser, err := s.userRepo.Create(&user)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}
