package authRoute

import (
	authController "github.com/dummyProjectBackend/controllers/auth"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoute(authGroup fiber.Router) {

	authGroup.Post("/login", authController.Login)
	authGroup.Post("/signup", authController.SignUp)
}
