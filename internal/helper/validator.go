package helper

import "github.com/go-playground/validator/v10"

type Validator struct {
	Validate *validator.Validate
}

func (v *Validator) ValidateStruct(req any) map[string]string {
	err := v.Validate.Struct(req)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	// Safely assert the type
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			errors[fe.Field()] = getCustomMessage(fe)
		}
	} else {
		errors["_error"] = err.Error()
	}

	return errors
}

func getCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Too short (minimum " + fe.Param() + " characters)"
	case "max":
		return "Too long (maximum " + fe.Param() + " characters)"
	}
	return "Invalid value"
}
