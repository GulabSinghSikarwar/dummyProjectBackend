package authController

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	ConfigFiles "github.com/dummyProjectBackend/config_files"
	"github.com/dummyProjectBackend/database"
	"github.com/dummyProjectBackend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *models.User, tokenType string) (string, error) {
	//  signing method   assinging
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()

	//  maping claims
	claims := token.Claims.(jwt.MapClaims)

	if tokenType == "access" {

		claims["exp"] = time.Now().Add(time.Hour * 15).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	}
	claims["sub"] = user.ID
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	// generating tokenString
	configData, _ := ConfigFiles.LoadConfig(("."))
	tokenString, err := token.SignedString([]byte(configData.JwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func getAssociatedWatchList(user *models.User) *models.Watchlist {
	watchlistCollection := database.OpenWatchListCollection(database.Client, "watchlist")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	filter := bson.M{"userId": user.ID}
	var exsistingWatchlist *models.Watchlist

	err := watchlistCollection.FindOne(ctx, filter).Decode(&exsistingWatchlist)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil
			// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "errors": err.Error(), "message": "something went wrong"})
		} else {

			log.Fatal(err.Error())
		}
	}
	return exsistingWatchlist

}
func Logout(c *fiber.Ctx) {
	expired := time.Now().Add(-time.Hour * 24)

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Expires: expired,
		Value:   "",
	})

}

func SignupController(c *fiber.Ctx) error {

	var payload *models.SignUpInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})

	}
	fmt.Println(" payload :", *payload)
	// fmt

	errors := models.ValidateStruct(*payload)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": errors,
		})
	}

	if payload.Password != payload.ConfirmPasswod {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "passwod and confirm password are not same ",
		})

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})

	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	newUser := models.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: string(hashedPassword),
		ID:       primitive.NewObjectID(),
	}
	var exsistedUser models.User
	err = database.OpenCollection(database.Client).FindOne(ctx, bson.M{"email": payload.Email}).Decode(&exsistedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {

			// ADDING NEW USER
			result, err := database.OpenCollection(database.Client).InsertOne(ctx, newUser)

			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": err.Error()})

			}

			//  TYPE ASSERTION  OF USER ID

			// CREATING NEW INSTANCE OF WATCH IST
			newWatchlist := models.Watchlist{
				UserId: result.InsertedID.(primitive.ObjectID),
				Stocks: []primitive.ObjectID{},
				ID:     primitive.NewObjectID(),
			}

			// CREATING NEW WATCHLIST IN DB
			createdWatchListResult, err := database.OpenWatchListCollection(database.Client, "watchlist").InsertOne(ctx, newWatchlist)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": err.Error()})

			}

			return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "userCreated": result, "watchlistCreated": createdWatchListResult})

		}
		if err != mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "errors": err.Error(), "message": "something went wrong"})
		}

	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "Alredy User Exsist with given Email Id"})

}
func SignInController(c *fiber.Ctx) error {

	var payload *models.SignInInput
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": err.Error()})

	}
	error := models.ValidateStruct(payload)

	if error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "error": error})
	}

	var user *models.User
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()
	err = database.OpenCollection(database.Client).FindOne(ctx, bson.M{"email": payload.Email}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": " failed", "message": "No User exsist with the given email"})

		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "Something Went Wrong"})

	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "Invalid Email or Password"})
	}

	configData, _ := ConfigFiles.LoadConfig(".")

	//  creating  token

	refreshToken, err := GenerateToken(user, "access")
	accessToken, err := GenerateToken(user, "refresh")

	c.Cookie(&fiber.Cookie{
		Name:     "jwtToken",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   configData.JwtMaxage * 120,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   configData.JwtMaxage * 120,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	watchlist := getAssociatedWatchList(user)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"token":        accessToken,
		"refreshToken": refreshToken,
		"user":         user,
		"watchlist":    watchlist,
	})
}
