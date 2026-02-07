package main

import (
	"habrexclude/internal/config"
	"habrexclude/internal/handlers"
	"habrexclude/internal/middleware"
	"time"

	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	baseLog := log.Default()
	app := fiber.New(fiber.Config{
		GETOnly: true,
		AppName: "HABR EXCLUDE",
	})

	app.Use(middleware.RateLimiter(1 * time.Second))
	app.Use(logger.New())

	baseLog.Println("Loading config...")
	config := config.New()
	baseLog.Println("Config loaded, initializing handlers...")
	handlers.InitHandler(app, config, baseLog)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Server started")

	go func() {
		<-c
		baseLog.Println("Shutting down server...")
		baseLog.Println("Server stopped")
		app.Shutdown()
	}()

	app.Listen(":3030")
}
