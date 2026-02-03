package main

import (
	"habrexclude/internal/config"
	"habrexclude/internal/handlers"

	"log"
	"os"
	"os/signal"
	"syscall"

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

	app.Use(logger.New())

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Server started on :80")

	go func() {
		<-c
		baseLog.Println("Shutting down server...")
		baseLog.Println("Server stopped")
		app.Shutdown()
	}()

	app.Listen(":80", fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}
