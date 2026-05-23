package services

import (
	"example.com/m/internal/models"
	"example.com/m/internal/repositories"
)

type CompanyService struct {
	companyRepo *repositories.CompanyRepository
}

func NewCompanyService(
	companyRepo *repositories.CompanyRepository,

) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
	}
}

func (s *CompanyService) CompanyCreate(req models.CreateCompanyRequest) (models.Company, error) {
	company := models.Company{
		Name:   req.Name,
		UserID: req.UserID,
	}

	createdCompany, err := s.companyRepo.CompanyCreate(&company)
	if err != nil {
		return models.Company{}, err
	}

	return createdCompany, nil
}
