package main

import (
	"fmt"

	"github.com/dummyProjectBackend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	fmt.Println(" server Started !! ")
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:4200",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST,DELETE,PUT,PATCH",
		AllowCredentials: true,
	}))

	routes.RegisterRoutes(app)
	// fmt.Println("client :", database.Client)

	app.Listen(":8000")
	// app.Server().ListenAndServe(":8000")

}
