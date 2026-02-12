package handlers

import (
	"habrexclude/internal/models"
	"log"

	"github.com/gofiber/fiber/v3"
)

func InitHandlers(app *fiber.App, config *models.Config, logger *log.Logger) {
	InitBlockHandler(app, config, logger)
	InitTestHandler(app, config, logger)
}
