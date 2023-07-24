package watchlistController

import (
	"context"
	"fmt"
	"time"

	"github.com/dummyProjectBackend/database"
	"github.com/dummyProjectBackend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllWatchlist(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	fmt.Println(user)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	watchlistCollection := database.OpenWatchListCollection(database.Client, "watchlist")

	var allWatchLists []models.Watchlist
	filter := bson.M{"userId": user.ID}
	cursor, err := watchlistCollection.Find(ctx, filter)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no watchlist exsist with this credential"})

		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "something went wrong "})

		}
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var watchlist models.Watchlist
		if err := cursor.Decode(&watchlist); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "something went wrong "})

		}
		allWatchLists = append(allWatchLists, watchlist)
	}
	allStocks := GetAllStocks(&allWatchLists[0])

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "result": allWatchLists, "allStocks": allStocks})

}
func DeleteStockFromWatchlist(c *fiber.Ctx) error {

	stockId := c.Params("stockId")
	watchlistId := c.Params("watchlistId")

	watchlistCollection := database.OpenWatchListCollection(database.Client, "watchlist")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var existingWatchlist *models.Watchlist

	filter := bson.M{"_id": watchlistId}
	err := watchlistCollection.FindOne(ctx, filter).Decode(&existingWatchlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no watchlist exsist with this credential"})

		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "something went wrong "})

		}
	}
	objectID, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})

	}
	deleteStockFromWatchlist(existingWatchlist, objectID)

	update := bson.M{"$set": bson.M{"stocks": existingWatchlist.Stocks}}
	err = watchlistCollection.FindOneAndUpdate(ctx, filter, update).Decode(existingWatchlist)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no watchlist exsist with this credential"})

		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "something went wrong "})

		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "result": existingWatchlist})

}

func deleteStockFromWatchlist(existingWatchlist *models.Watchlist, targetStockId primitive.ObjectID) {
	index := -1

	for i, stockId := range existingWatchlist.Stocks {
		if stockId == targetStockId {
			index = i

			break

		}

	}
	if index >= 0 {
		existingWatchlist.Stocks = append(existingWatchlist.Stocks[:index], existingWatchlist.Stocks[:index+1]...)

	}

}
func GetAllStocks(existingWatchlist *models.Watchlist) []*models.Stock {
	var stocks []*models.Stock
	stockCollection := database.OpenWatchListCollection(database.Client, "stock")
	fmt.Println("exsisting watchlist : :", existingWatchlist)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, stockId := range existingWatchlist.Stocks {
		// fmt.Println("------------------stockId ---------------: ", stockId)
		filter := bson.M{"_id": stockId}
		var stock *models.Stock

		err := stockCollection.FindOne(ctx, filter).Decode(&stock)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			} else {
				return []*models.Stock{}
			}

		}
		stocks = append(stocks, stock)

	}
	return stocks

}

// Read
func GetSingleWatchList(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	userId := user.ID
	watchlistId := c.Params("watchlistId")
	WatchListCollection := database.OpenWatchListCollection(database.Client, "watchlist")

	var existingWatchlist *models.Watchlist
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"userId": userId, "_id": watchlistId}

	err := WatchListCollection.FindOne(ctx, filter).Decode(&existingWatchlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no watchlist exsist with such infomation"})
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
			}

		}

	}
	stocks := GetAllStocks(existingWatchlist)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"result": existingWatchlist,
		"stocks": stocks,
	})
}

// NOT TO USE
func getWatchList(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	userId := user.ID
	var payload *models.WatchlistGetReqBody

	err := c.BodyParser(&payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	errors := models.ValidateAddStockRequestBody(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": errors})

	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	WatchListCollection := database.OpenWatchListCollection(database.Client, "watchlist")

	var existingWatchlist *models.Watchlist

	filter := bson.M{"_id": payload.WatchlistId, "userId": userId}

	err = WatchListCollection.FindOne(ctx, filter).Decode(&existingWatchlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no watchlist exsist with such infomation"})
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
			}

		}

	}
	stocks := GetAllStocks(existingWatchlist)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"result": existingWatchlist,
		"stocks": stocks,
	})

}

// Create
func AddStockToWatchList(c *fiber.Ctx) error {

	user := c.Locals("user").(*models.User)
	watchlistId := c.Params("watchlistId")
	watchlistObjectId, err := primitive.ObjectIDFromHex(watchlistId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	var payload *models.AddStockRequestBody

	err = c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := models.ValidateAddStockRequestBody(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": errors})

	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	WatchListCollection := database.OpenWatchListCollection(database.Client, "watchlist")

	var existingWatchlist *models.Watchlist
	filter := bson.M{"_id": watchlistObjectId, "userId": user.ID}

	//  GET WATCHLIST WITH GIVEN ID

	err = WatchListCollection.FindOne(ctx, filter).Decode(&existingWatchlist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no watchlist exsist with such infomation"})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
		}

	}

	stockSymbol := payload.StockSymbol
	fmt.Println(" symbol : ", stockSymbol)
	var exsistingStock *models.Stock
	err = database.OpenWatchListCollection(database.Client, "stock").FindOne(ctx, bson.M{"ticker": stockSymbol}).Decode(&exsistingStock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no stock exsist with such infomation"})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed ::", "message": err.Error()})
		}

	}

	// payload_id, _ := primitive.ObjectIDFromHex(payload.S0)

	existingWatchlist.Stocks = append(existingWatchlist.Stocks, existingWatchlist.ID)
	filter = bson.M{
		"_id": watchlistObjectId, "userId": user.ID,
	}

	update := bson.M{"$set": bson.M{"stocks": existingWatchlist.Stocks}}
	fmt.Println("updated stock : ", existingWatchlist)

	result, err := WatchListCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "Something Went Wrong "})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "result": result, "stockAdded": exsistingStock})

}
