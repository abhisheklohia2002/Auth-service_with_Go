package repositories

import (
	"context"
	"errors"
	"time"

	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) (models.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return models.User{}, err
	}

	return *user, nil
}

func (r *UserRepository) UsersList() ([]models.User, error) {
	var users []models.User

	err := r.db.
		Preload("Role").
		Preload("Department").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) DeleteUserById(id int) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).Delete(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, err
	}

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUserById(ctx context.Context, id uint, req models.UpdateUserRequest) (*models.User, error) {
	var user models.User

	fullName := req.FullName
	if fullName == "" {
		fullName = req.Name
	}

	updates := map[string]interface{}{
		"full_name":     fullName,
		"email":         req.Email,
		"employee_code": req.EmployeeCode,
		"role_id":       req.RoleID,
		"manager_id":    req.ManagerID,
		"department_id": req.DepartmentID,
		"status":        req.Status,
	}

	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Department").
		First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (r *UserRepository) ExistsByEmail(email string) (bool, *models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil, nil
	}

	if err != nil {
		return false, nil, err
	}

	return true, &user, nil
}

func (r *UserRepository) PersistRefreshToken(
	userID uint,
	tokenHash string,
	expiresAt time.Time,
) (*models.RefreshToken, error) {
	now := time.Now()

	err := r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
	if err != nil {
		return nil, err
	}

	refreshToken := models.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *UserRepository) FindByIDWithRole(userID uint) (models.User, error) {
	var user models.User

	err := r.db.
		Preload("Role").
		Preload("Department").
		First(&user, userID).Error

	return user, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User

	err := r.db.
		Preload("Role").
		Preload("Department").
		First(&user, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindValidRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	err := r.db.
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&refreshToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *UserRepository) RevokeRefreshToken(tokenHash string) error {
	now := time.Now()

	return r.db.Model(&models.RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", &now).Error
}

func (r *UserRepository) FindActiveByDepartmentID(departmentID uint) ([]models.User, error) {
	var users []models.User

	err := r.db.
		Preload("Role").
		Preload("Department").
		Where("department_id = ? AND status = ?", departmentID, "active").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) FindExistingUsersByEmailOrEmployeeID(
	ctx context.Context,
	emails []string,
	employeeIDs []string,
) ([]models.User, error) {
	var users []models.User

	if len(emails) == 0 && len(employeeIDs) == 0 {
		return users, nil
	}

	err := r.db.WithContext(ctx).
		Where("LOWER(email) IN ? OR UPPER(employee_code) IN ?", emails, employeeIDs).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CreateUsersInBatches(
	ctx context.Context,
	users []models.User,
	batchSize int,
) error {
	if len(users) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		CreateInBatches(&users, batchSize).
		Error
}

func (r *UserRepository) GetAllUserIDs(ctx context.Context) ([]uint, error) {
	var userIDs []uint

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Pluck("id", &userIDs).Error

	return userIDs, err
}
