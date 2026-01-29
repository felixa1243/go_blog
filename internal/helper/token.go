package helper

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

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
	idVal, exists := claims["user_id"]
	if !exists || idVal == nil {
		return "", false
	}

	return fmt.Sprintf("%v", idVal), true
}
