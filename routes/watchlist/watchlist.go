package watchlistRoute

import (
	watchlistController "github.com/dummyProjectBackend/controllers/watchlist"
	"github.com/dummyProjectBackend/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterWatchlistRoute(watchlistGroup fiber.Router) {

	watchlistGroup.Delete("/:watchlistId/:stockId", middleware.AuthenticateUser, watchlistController.DeleteStockFromWatchlist)
	watchlistGroup.Post("/:watchlistId", middleware.AuthenticateUser, watchlistController.AddStockToWatchList)
	watchlistGroup.Get("/:watchlistId", middleware.AuthenticateUser, watchlistController.GetSingleWatchList)
	watchlistGroup.Get("/", middleware.AuthenticateUser, watchlistController.GetAllWatchlist)

}
