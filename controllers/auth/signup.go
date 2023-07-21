package authController

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dummyProjectBackend/database"
	"github.com/dummyProjectBackend/models"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// api/user/login

func SignUp(c *fiber.Ctx) error {

	// reqBody := c.Body()
	// fmt.Println(reqBody)
	// var body1 model.User

	// json.Unmarshal(reqBody, &body1)
	// fmt.Println(body1)
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  recieved user ", user)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// validator

	validate := validator.New()
	validationErr := validate.Struct(user)
	if validationErr != nil {
		// Handle validation errors
		// This block will be executed if there are validation errors
	}

	count, err := database.OpenCollection(database.Client, "user").CountDocuments(ctx, bson.M{"email": user.Email})

	if err != nil {
		return err

	}

	count, err = database.OpenCollection(database.Client, "user").CountDocuments(ctx, bson.M{"phone": user.Phone})

	if err != nil {
		return err

	}
	if count > 0 {

	}

	return c.SendString("Signup route")
}
