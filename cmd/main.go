package main

import (
	docs "habrexclude/docs"
	"habrexclude/internal/config"
	"habrexclude/internal/handlers"
	"habrexclude/internal/middleware"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

// @title Habr Parser API
// @version 1.0
// @description REST API for parsing and analyzing content from habr.com/ru
// @description
// @description ## Core Features:
// @description • Parse articles, news, and posts from Habr
// @description • Filter content by categories, tags, and rating
// @description • Pagination and multiple sorting options
//
// @host localhost:3030
// @BasePath /

func main() {
	app := fiber.New(fiber.Config{
		GETOnly:      true,
		ServerHeader: "HABR EXCLUDE",
		AppName:      "HABR EXCLUDE",
	})
	app.Use(middleware.RateLimiter(1 * time.Second))
	app.Use(logger.New())

	baseLog := log.Default()
	baseLog.Println("Loading config...")
	config := config.New()

	baseLog.Println("Initializing handlers...")
	handlers.InitHandler(app, config, baseLog)

	if config.Mode == "dev" {
		if config.SwaggerHost != "" {
			docs.SwaggerInfo.Host = config.SwaggerHost
		}
		app.Get("/swagger/*", swaggo.HandlerDefault)
		app.Get("/docs/*", swaggo.New(swaggo.Config{
			Title:        "API UI",
			URL:          "./docs/swagger.json",
			DeepLinking:  false,
			DocExpansion: "none",
		}))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	baseLog.Println("Server started")

	go func() {
		<-c
		baseLog.Println("Shutting down server...")
		baseLog.Println("Server stopped")
		app.Shutdown()
	}()

	app.Listen(":3030")
}
