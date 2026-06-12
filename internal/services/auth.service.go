package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"net/mail"
	"path/filepath"
	"strings"
	"time"

	"example.com/m/internal/config"
	"example.com/m/internal/dto"
	"example.com/m/internal/helper"
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
	"github.com/xuri/excelize/v2"
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

func (s *AuthService) CreateBulkUsers(
	ctx context.Context,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (*dto.BulkUsersUploadResponse, error) {

	if file == nil || fileHeader == nil {
		return nil, errors.New("file is required")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".xlsx" {
		return nil, errors.New("only .xlsx files are allowed")
	}

	if fileHeader.Size > 10*1024*1024 {
		return nil, errors.New("file size must be less than 10MB")
	}

	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("failed to read excel file")
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, errors.New("excel sheet is empty")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, errors.New("failed to read excel rows")
	}

	if len(rows) <= 1 {
		return nil, errors.New("excel file has no user data")
	}

	response := &dto.BulkUsersUploadResponse{
		Errors: []dto.BulkUsersError{},
	}

	seenEmails := make(map[string]int)
	seenEmpIDs := make(map[string]int)

	type excelUser struct {
		Row        int
		Name       string
		Email      string
		EmployeeID string
		Password   string
	}

	validRows := []excelUser{}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if helper.IsEmptyRow(row) {
			continue
		}

		rowNumber := i + 1

		name := strings.TrimSpace(helper.GetCell(row, 0))
		email := strings.ToLower(strings.TrimSpace(helper.GetCell(row, 1)))
		employeeID := strings.ToUpper(strings.TrimSpace(helper.GetCell(row, 2)))
		password := strings.TrimSpace(helper.GetCell(row, 3))

		if helper.IsHeaderRow(name, email, employeeID, password) {
			continue
		}

		response.TotalRows++

		if name == "" || email == "" || employeeID == "" || password == "" {
			response.Errors = append(response.Errors, dto.BulkUsersError{
				Row:     rowNumber,
				Email:   email,
				EmpID:   employeeID,
				Message: "name, email, employee code and password are required",
			})
			continue
		}

		if _, err := mail.ParseAddress(email); err != nil {
			response.Errors = append(response.Errors, dto.BulkUsersError{
				Row:     rowNumber,
				Email:   email,
				EmpID:   employeeID,
				Message: "invalid email format",
			})
			continue
		}

		if previousRow, exists := seenEmails[email]; exists {
			response.Errors = append(response.Errors, dto.BulkUsersError{
				Row:     rowNumber,
				Email:   email,
				EmpID:   employeeID,
				Message: fmt.Sprintf("duplicate email in excel, already used at row %d", previousRow),
			})
			continue
		}

		if previousRow, exists := seenEmpIDs[employeeID]; exists {
			response.Errors = append(response.Errors, dto.BulkUsersError{
				Row:     rowNumber,
				Email:   email,
				EmpID:   employeeID,
				Message: fmt.Sprintf("duplicate employeeId in excel, already used at row %d", previousRow),
			})
			continue
		}

		seenEmails[email] = rowNumber
		seenEmpIDs[employeeID] = rowNumber

		validRows = append(validRows, excelUser{
			Row:        rowNumber,
			Name:       name,
			Email:      email,
			EmployeeID: employeeID,
			Password:   password,
		})
	}

	if len(validRows) == 0 {
		response.FailedCount = len(response.Errors)
		return response, nil
	}

	emails := make([]string, 0, len(validRows))
	empIDs := make([]string, 0, len(validRows))

	for _, user := range validRows {
		emails = append(emails, user.Email)
		empIDs = append(empIDs, user.EmployeeID)
	}

	existingUsers, err := s.userRepo.FindExistingUsersByEmailOrEmployeeID(
		ctx,
		emails,
		empIDs,
	)
	if err != nil {
		return nil, errors.New("failed to check existing users")
	}

	existingEmailMap := make(map[string]bool)
	existingEmpIDMap := make(map[string]bool)

	for _, user := range existingUsers {
		existingEmailMap[strings.ToLower(user.Email)] = true
		existingEmpIDMap[strings.ToUpper(user.EmployeeCode)] = true
	}

	usersToCreate := []models.User{}

	for _, rowUser := range validRows {
		if existingEmailMap[rowUser.Email] {
			response.Errors = append(response.Errors, dto.BulkUsersError{
				Row:     rowUser.Row,
				Email:   rowUser.Email,
				EmpID:   rowUser.EmployeeID,
				Message: "email already exists",
			})
			continue
		}

		if existingEmpIDMap[rowUser.EmployeeID] {
			response.Errors = append(response.Errors, dto.BulkUsersError{
				Row:     rowUser.Row,
				Email:   rowUser.Email,
				EmpID:   rowUser.EmployeeID,
				Message: "employeeId already exists",
			})
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(rowUser.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			response.Errors = append(response.Errors, dto.BulkUsersError{
				Row:     rowUser.Row,
				Email:   rowUser.Email,
				EmpID:   rowUser.EmployeeID,
				Message: "failed to hash password",
			})
			continue
		}

		usersToCreate = append(usersToCreate, models.User{
			FullName:     rowUser.Name,
			Email:        rowUser.Email,
			EmployeeCode: rowUser.EmployeeID,
			Password:     string(hashedPassword),

			RoleID:      3,
			Status:      "active",
			JoiningDate: time.Now(),
		})
	}

	if len(usersToCreate) > 0 {
		err := s.userRepo.CreateUsersInBatches(ctx, usersToCreate, 100)
		if err != nil {
			return nil, errors.New("failed to create users, please check duplicate email or employeeId")
		}
	}

	response.SuccessCount = len(usersToCreate)
	response.FailedCount = len(response.Errors)

	return response, nil
}
