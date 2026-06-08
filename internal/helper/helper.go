package helper

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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