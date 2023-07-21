package routes

import (
	authRoute "github.com/dummyProjectBackend/routes/auth"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	authRoute.RegisterAuthRoute(api.Group("/auth"))

}
