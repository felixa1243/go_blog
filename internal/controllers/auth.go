package controllers

import (
	"fmt"
	"go_blog/internal/database"
	"go_blog/internal/dto"
	"go_blog/internal/models"
	"go_blog/internal/utils"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := utils.Validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
func Register(c *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body: " + err.Error(),
		})
	}

	myValidator := &XValidator{validator: utils.Validate}
	if errs := myValidator.Validate(req); len(errs) > 0 {
		errMsgs := make([]string, 0)
		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf("[%s]: '%v' | Failed on condition: %s", err.FailedField, err.Value, err.Tag))
		}
		return c.Status(400).JSON(fiber.Map{
			"message": strings.Join(errMsgs, ", "),
		})
	}
	foundUser := database.New().GetDB().Where("email = ?", req.Email).First(&models.User{})
	if foundUser.RowsAffected > 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": fmt.Sprint("User with email ", req.Email, " already exists"),
		})
	}

	user := models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: utils.GeneratePassword(req.Password),
	}

	res := database.New().GetDB().Create(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Could not create user: " + res.Error.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "user created",
		"user":    user,
	})
}
func Login(c *fiber.Ctx) error {
	var req dto.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var user models.User
	res := database.New().GetDB().Where("email = ?", req.Email).First(&user)
	if res.Error != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "user not found",
		})
	}
	if !utils.ComparePassword(user.Password, req.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response := dto.LoginResponse{
		Message: "User Logged in Successfully",
		Token:   token,
		User: dto.UserLoginResponse{
			Name:  user.Name,
			Email: user.Email,
		},
	}
	return c.JSON(&response)
}
