package middleware

import (
	"os"
	"strings"

	"github.com/G1GACHADS/stashable-backend/token/jwt"
	"github.com/gofiber/fiber/v2"
)

func Authenticated(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	if len(strings.Split(authorization, " ")) < 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please login to continue",
		})
	}
	accessToken := strings.Split(authorization, " ")[1]

	_, claims, err := jwt.Verify(accessToken, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please login to continue",
		})
	}

	userID, ok := claims["userID"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please try and login again",
		})
	}

	email, ok := claims["email"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Please try and login again",
		})
	}

	c.Locals("userID", userID)
	c.Locals("email", email)
	c.Locals("accessToken", accessToken)

	return c.Next()
}
