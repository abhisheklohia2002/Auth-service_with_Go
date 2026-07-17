package helper

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/m/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Pagination struct {
	page     int
	pageSize int
}

func HasAudience(audiences jwt.ClaimStrings, target string) bool {
	for _, aud := range audiences {
		if aud == target {
			return true
		}
	}

	return false
}

func ParseUintParam(c *gin.Context, paramName string) (uint, bool) {
	rawID := c.Param(paramName)

	id64, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid " + paramName,
		})
		return 0, false
	}

	return uint(id64), true
}

func IsValidQuestionType(questionType string) bool {
	switch questionType {
	case "single_choice", "multiple_choice", "true_false":
		return true
	default:
		return false
	}
}

func IsCorrectLabel(label string, correctAnswer string) bool {
	answers := strings.Split(correctAnswer, ",")
	for _, answer := range answers {
		if strings.TrimSpace(answer) == label {
			return true
		}
	}
	return false
}

func GetCell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func IsEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

func IsHeaderRow(name, email, employeeCode, password string) bool {
	return strings.EqualFold(name, "Full Name") ||
		strings.EqualFold(email, "Email") ||
		strings.EqualFold(employeeCode, "Employee Code") ||
		strings.EqualFold(password, "Password")
}

func GetEffectiveAssessmentRule(assessment *models.Assessment) (maxAttempts int, retakeAllowed bool, passingScore int) {
	maxAttempts = 1
	retakeAllowed = false
	passingScore = assessment.PassingScore

	if assessment.Rule != nil {
		if assessment.Rule.MaxAttempts > 0 {
			maxAttempts = assessment.Rule.MaxAttempts
		}

		retakeAllowed = assessment.Rule.RetakeAllowed

		if assessment.Rule.PassingScore > 0 {
			passingScore = assessment.Rule.PassingScore
		}
	}

	return maxAttempts, retakeAllowed, passingScore
}

func NewPagination(page, pageSize int) *Pagination {
	if page < 1 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	return &Pagination{
		page:     page,
		pageSize: pageSize,
	}
}

func (p *Pagination) GetPage() int {
	return p.page
}

func (p *Pagination) GetPageSize() int {
	return p.pageSize
}

func (p *Pagination) ModifyStatement(stmt *gorm.Statement) {
	stmt.AddClause(clause.Limit{
		Limit:  &p.pageSize,
		Offset: (p.page - 1) * p.pageSize,
	})
}

func (p *Pagination) Build(clause.Builder) {}
