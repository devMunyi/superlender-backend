package controllers

import (
	"fmt"
	"super-lender/inits"
	"super-lender/models"
	"time"
)

// create a function named ValidateToken that takes token as a parameter and returns a boolean
func ValidateToken(token string) int {
	var tokenCount int64
	expiryDate := time.Now()

	result := inits.CurrentDB.Model(&models.OToken{}).Where("token = ? AND expiry_date > ?", token, expiryDate).Count(&tokenCount)

	if result.Error != nil {
		fmt.Println("Error validating token: ", result.Error)
		return 0
	}

	if tokenCount > 0 {
		return 1
	} else {
		return 0
	}
}
