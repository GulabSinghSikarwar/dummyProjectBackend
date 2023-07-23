package authController

import (
	"context"
	"fmt"
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

	token := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()

	claims := token.Claims.(jwt.MapClaims)
	// claims["exp"] = now.Add(configData.JwtExpiresIn).Unix()
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["sub"] = user.ID
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(configData.JwtSecret))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "failed",
			"message": "something went wrong",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwtToken",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   configData.JwtMaxage * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"user":   user,
		"token":  tokenString,
	})
}
