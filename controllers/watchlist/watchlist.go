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

func GetStockById(c *fiber.Ctx) error {

	stockId := c.Params("stockId")
	stockObjectId, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status ": "failed  ", "message  ": err.Error()})

	}
	stockCollection := database.OpenWatchListCollection(database.Client, "stock")

	var stock *models.Stock

	ctx := context.Background()
	filter := bson.M{"_id": stockObjectId}
	err = stockCollection.FindOne(ctx, filter).Decode(&stock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status ": "failed due to absence of required data ", "message  ": err.Error()})
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": " success  ", "stock": stock})

}

func GetAllStockInTheWatchList(c *fiber.Ctx) error {

	watchlistId := c.Params("watchlistId")
	watchlistObjectId, err := primitive.ObjectIDFromHex(watchlistId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "error occured while converting watchlist to objectId"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// watchlist collection
	watchListCollection := database.OpenWatchListCollection(database.Client, "watchlist")
	// stock collection
	stockCollection := database.OpenWatchListCollection(database.Client, "stock")

	var exsistingWatchlist *models.Watchlist

	filter := bson.M{"id": watchlistObjectId}

	err = watchListCollection.FindOne(ctx, filter).Decode(&exsistingWatchlist)
	if err != nil {

		if err == mongo.ErrNoDocuments {

			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no document Found "})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed , while finding watchlist ", "message": err.Error()})

		}

	}
	filter = bson.M{"_id": bson.M{"$in": exsistingWatchlist.Stocks}}

	// Find all documents that match the query.
	cursor, err := stockCollection.Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}

	// Iterate over the cursor and print the documents.
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}

		fmt.Println(result)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success  "})

}

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
	// var user models.User

	watchlistCollection := database.OpenWatchListCollection(database.Client, "watchlist")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()

	var existingWatchlist *models.Watchlist
	watchlistObjectId, err := primitive.ObjectIDFromHex(watchlistId)

	if err != nil {
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"status": "failed", "error": err.Error()})
	}

	filter := bson.M{"_id": watchlistObjectId}
	err = watchlistCollection.FindOne(ctx, filter).Decode(&existingWatchlist)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no watchlist exsist with this credential"})

		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "something went wrong "})

		}
	}

	var user models.User
	if u, ok := c.Locals("user").(models.User); ok {
		user = u
		if user.ID != existingWatchlist.UserId {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "error": "token id and watchlist userId do not match "})

		}
	}

	objectID, err := primitive.ObjectIDFromHex(stockId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})

	}
	fmt.Println(" stocks  before : ", existingWatchlist.Stocks)
	deleteStockFromWatchlist(existingWatchlist, objectID)
	fmt.Println(" stocks  after : ", existingWatchlist.Stocks)

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

	fmt.Println("STocks passesd : ", existingWatchlist.Stocks)
	fmt.Println(" target StockId : ", targetStockId)

	for i, stockId := range existingWatchlist.Stocks {
		fmt.Println(" single stock Id : ", stockId)
		fmt.Println(" compare : ", (targetStockId == stockId))
		if stockId == targetStockId {
			index = i

			break

		}

	}
	fmt.Println(" index : ", index)
	if len(existingWatchlist.Stocks) == 1 {
		existingWatchlist.Stocks = []primitive.ObjectID{}
	} else if index == 0 {
		existingWatchlist.Stocks = existingWatchlist.Stocks[1:]

	} else {
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

		fmt.Println("stockId :", stockId)
		filter := bson.M{"_id": stockId}
		var stock *models.Stock

		err := stockCollection.FindOne(ctx, filter).Decode(&stock)
		fmt.Println("stock with detail --> :", stock)
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
	// user := c.Locals("user").(*models.User)
	// userId := user.ID
	watchlistId := c.Params("watchlistId")
	WatchListCollection := database.OpenWatchListCollection(database.Client, "watchlist")
	watchlistObjectId, err := primitive.ObjectIDFromHex(watchlistId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": err.Error()})
	}

	var existingWatchlist *models.Watchlist
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"_id": watchlistObjectId}
	fmt.Println("filter : ", filter)

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

	filter = bson.M{"_id": bson.M{"$in": existingWatchlist.Stocks}}

	stockCollection := database.OpenWatchListCollection(database.Client, "stock")

	var result []*models.Stock

	csr, err := stockCollection.Find(context.Background(), filter)

	if err == mongo.ErrNoDocuments {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "no Stocks exsist with such infomation"})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
		}

	}
	err = csr.All(context.Background(), &result)
	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})
	}

	fmt.Println("results : ", result)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"result": existingWatchlist,
		"stocks": result,
		// "stocks": stocks,
	})
}

// NOT TO USE

// Create
func AddStockToWatchList(c *fiber.Ctx) error {

	user := c.Locals("user").(*models.User)
	watchlistId := c.Params("watchlistId")
	fmt.Println("id wid : ", watchlistId)
	watchlistObjectId, err := primitive.ObjectIDFromHex(watchlistId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail :::", "message": err.Error()})
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

	if user.ID != existingWatchlist.UserId {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "error": "token id and watchlist userId do not match "})

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
	for _, stockid := range existingWatchlist.Stocks {
		fmt.Println("existingWatchlist.ID ", existingWatchlist.ID, " stockid : ", stockid)
		if stockid == exsistingStock.ID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed ::", "message": "stock already exsist"})
		}
	}

	// payload_id, _ := primitive.ObjectIDFromHex(payload.S0)

	existingWatchlist.Stocks = append(existingWatchlist.Stocks, exsistingStock.ID)
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
