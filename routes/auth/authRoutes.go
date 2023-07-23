package authRoute

import (
	authController "github.com/dummyProjectBackend/controllers/auth"
	"github.com/gofiber/fiber/v2"
)

// RegisterAuthRoute registers the authentication routes under the given authGroup
func RegisterAuthRoute(authGroup fiber.Router) {
	// POST /login route for user login
	authGroup.Post("/login", authController.SignInController)

	// POST /signup route for user signup
	authGroup.Post("/signup", authController.SignupController)
}
