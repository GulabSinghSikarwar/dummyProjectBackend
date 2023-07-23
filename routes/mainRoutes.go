package routes

import (
	authRoute "github.com/dummyProjectBackend/routes/auth"
	watchlistRoute "github.com/dummyProjectBackend/routes/watchlist"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	authRoute.RegisterAuthRoute(api.Group("/auth"))

	watchlistRoute.RegisterWatchlistRoute(api.Group("/watchlist"))
}
