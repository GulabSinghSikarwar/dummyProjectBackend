package middleware

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
)

func AuthenticateUser(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	var tokenString string
	// var refreshToken string

	if strings.HasPrefix(authorization, "Bearer") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
		fmt.Println("through header :", tokenString)

	} else if jwtToken := c.Cookies("jwtToken"); jwtToken != "" {
		tokenString = jwtToken
		fmt.Println("through cookie :", jwtToken)
	}

	fmt.Println("cookies :", c.Cookies("jwtToken"))
	// If no token is found, return an error response
	if tokenString == "" {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "You are not logged in"})
	}

	// Load your JWT secret from your configuration
	config, err := ConfigFiles.LoadConfig(".")
	if err != nil {
		fmt.Println("err: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "Internal Server Error"})
	}

	// Parse the token using the secret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})

	// Handle parsing errors
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "Invalid token claim"})
	}

	// Check if the token is valid and not expired
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failed", "message": "Invalid token claim"})
	}

	// Check token expiration
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().UTC().After(expirationTime) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "failed", "message": "Token has expired"})
	}

	userID := fmt.Sprint(claims["sub"])

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	user_objectId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "something went wrong"})

	}
	err = database.OpenCollection(database.Client).FindOne(ctx, bson.M{"_id": user_objectId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "user with token does not exsist"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "failed", "message": "something went wrong"})

	}
	fmt.Println("user : after  token verify ", user)

	c.Set("Authorization", "Bearer "+tokenString)
	c.Set("Refresh-Token", "refresh-token")
	fmt.Println("tokenString : ..........", tokenString)

	fmt.Println("userId : ", user.ID)

	c.Locals("user", &user)
	return c.Next()
}
