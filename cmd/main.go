package main

import (
	// "fmt"
	"habrexclude/internal/config"
	"habrexclude/internal/handlers"
	// "habrexclude/internal/parser"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {

	config := config.New()
	baseLog := log.Default()
	app := fiber.New(fiber.Config{
		GETOnly: true,
		AppName: "HABR EXCLUDE",
	})

	handlers.InitHandler(app, config, baseLog)

	// test := parser.NewArticleFetcher(config)
	// res,err := test.GetAll(1)
	// fmt.Println(res[0], err)

	app.Use(logger.New())
	// app.Listen(":8080")
}
