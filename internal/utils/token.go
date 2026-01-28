package utils

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(id uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
	})

	// Convert the string to []byte here
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// You already did it correctly here!
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return false, err
	}

	return token.Valid, nil
}
func GetUserID(c *fiber.Ctx) (string, bool) {
	userToken, ok := c.Locals("jwt").(*jwt.Token)
	if !ok {
		return "", false
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	// Use fmt.Sprintf to safely convert whatever is in "user_id" to a string
	// This handles both uuid.UUID and string types during the transition
	idVal, exists := claims["user_id"]
	if !exists || idVal == nil {
		return "", false
	}

	return fmt.Sprintf("%v", idVal), true
}
