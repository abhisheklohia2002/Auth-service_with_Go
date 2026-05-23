package services

import (
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
)

type DepartmentService struct {
	departmentRepo *repositories.DepartmentRepository
}

func NewDepartmentService(departmentRepo *repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{departmentRepo: departmentRepo}
}

func (d *DepartmentService) CreateDepartment(req *models.CreateCompanyRequest) (models.Department, error) {
	department := models.Department{
		Name:   req.Name,
		UserID: req.UserID,
	}

	createDepartment, err := d.departmentRepo.CreateDepartment(&department)
	if err != nil {
		return models.Department{}, err
	}
	return createDepartment, nil
}

func (d *DepartmentService) ListDepartment() ([]models.Department, error) {
	return d.departmentRepo.ListDepartments()
}

func (d *DepartmentService) UserIdByDepartment(userId int) ([]models.Department, error) {
	return d.departmentRepo.UserIdByDepartment(userId)
}
