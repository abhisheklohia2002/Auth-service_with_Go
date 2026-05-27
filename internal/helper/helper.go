package helper

import (
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

