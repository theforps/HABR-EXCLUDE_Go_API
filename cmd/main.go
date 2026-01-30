package main

import (
	"habrexclude/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

func main() {

	app := fiber.New(fiber.Config{
		GETOnly: true,
		AppName: "HABR EXCLUDE",
		
    })

	api := app.Group("/api")
	go handlers.InitHandler(api)

	app.Listen(":8080")
}
