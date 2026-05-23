package repositories

import (
	"example.com/m/internal/models"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{
		db: db,
	}
}

func (r *CompanyRepository) CompanyCreate(company *models.Company) (models.Company, error) {
	if err := r.db.Create(company).Error; err != nil {
		return models.Company{}, err
	}

	return *company, nil
}
