package main

import (
	"fmt"

	"github.com/dummyProjectBackend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	fmt.Println(" server Started !! ")
	app.Use(logger.New())
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     "http://localhost:4200",
	// 	AllowHeaders:     "Origin, Content-Type, Accept,Authorization",
	// 	AllowMethods:     "GET, POST,DELETE,PUT,PATCH",
	// 	AllowCredentials: true,
	// }))
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Set("Access-Control-Allow-Headers", "Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Credentials", "true") // Allow credentials (e.g., cookies)

		// Handle preflight requests
		if c.Method() == fiber.MethodOptions {
			c.Set("Access-Control-Allow-Headers", c.Get("Access-Control-Request-Headers"))
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	routes.RegisterRoutes(app)
	// fmt.Println("client :", database.Client)

	app.Listen(":8000")
	// app.Server().ListenAndServe(":8000")

}
