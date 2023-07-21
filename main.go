package main

import (
	"fmt"

	"github.com/dummyProjectBackend/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	fmt.Println(" server startsed ")

	routes.RegisterRoutes(app)

	app.Listen(":8000")

}
