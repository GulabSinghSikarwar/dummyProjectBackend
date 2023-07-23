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
	var tokenString string

	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer") {
		tokenString = strings.Trim(authorization, "Bearer")
	} else if c.Cookies("jwtToken") != "" {
		tokenString = c.Cookies("jwtToken")

	}
	if tokenString == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": "You are not  loggedIn "})
	}
	fmt.Println("token:String , : ", tokenString)
	config, _ := ConfigFiles.LoadConfig(".")

	token, err := jwt.Parse(
		tokenString,
		func(jwtToken *jwt.Token) (interface{}, error) {

			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", jwtToken.Header["alg"])
			}
			fmt.Println("token: ", jwtToken)

			// if !ok {
			// 	return nil, fmt.Errorf(" unexpected or unsopported  sigining method %s", jwtToken.Header["alg"])

			// }
			return []byte(config.JwtSecret), nil

		})

	fmt.Println("token ::: ", token, "  err::::", err.Error())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": " invalid token  claim"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	fmt.Println("token  --> claims::: ", claims, "  ::: ok ::: ", ok)

	if !ok || !token.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "failed", "message": " invalid token  claim"})

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

	c.Locals("user", &user)
	return c.Next()
}
