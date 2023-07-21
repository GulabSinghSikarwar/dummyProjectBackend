package authController

import "github.com/gofiber/fiber/v2"

// api/user/login

func Login(c *fiber.Ctx) error {
	return c.SendString("Login route")
}
