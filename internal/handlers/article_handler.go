package handlers

import (
	"habrexclude/internal/models"
	"log"

	"github.com/gofiber/fiber/v3"
)

type handler struct{
	config *models.Config
	logger *log.Logger
}

func InitHandler(app *fiber.App, conf *models.Config, log *log.Logger) {

	hand := handler{
		config: conf,
		logger: log,
	}

	api := app.Group("/api")
	api.Get("/test", hand.GetArticles)
}

func (h *handler) GetArticles(c fiber.Ctx) error {

	


	return c.SendString("test")
}

// Get Posts

// Get News

// Search

// GetNote
