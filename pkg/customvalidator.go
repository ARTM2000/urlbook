package pkg

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpperCase := false
	hasLowerCase := false
	hasDigit := false
	hasSymbol := false

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsLower(c):
			hasLowerCase = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSymbol = true
		}
	}

	return hasUpperCase && hasLowerCase && hasDigit && hasSymbol
}

func validateFilename(fl validator.FieldLevel) bool {
	malformedChars := []string{"^", "<", ">", ";", "|", "'", "/", ",", "\\", ":", "=", "?", "\"", "*"}

	filename := fl.Field().String()
	if filename != "" {
		return true
	}

	for _, ch := range malformedChars {
		if strings.Contains(filename, ch) {
			return false
		}
	}

	return true
}
