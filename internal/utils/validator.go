package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Exported variable (Capitalized) so other packages can use it
var Validate = validator.New()

func init() {
	// Register the custom tag here once
	Validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
		hasNumber := strings.ContainsAny(password, "0123456789")
		return len(password) >= 8 && hasUpper && hasLower && hasNumber
	})
}
