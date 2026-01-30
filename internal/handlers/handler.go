package handlers

import "github.com/gofiber/fiber/v3"

type baseHandler struct{}

func InitHandler(app fiber.Router) {

	handler := baseHandler{}

	app.Get("/test", handler.GetArticles)
}

func (bs *baseHandler) GetArticles(c fiber.Ctx) error {

	


	return c.SendString("test")
}

// Get Posts

// Get News

// Search

// GetNote
