package handlers

import (
	"habrexclude/internal/models"

	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type TestHandler struct {
	logger *log.Logger
	config *models.Config
}

func InitTestHandler(app *fiber.App, config *models.Config, log *log.Logger) {
	testHandler := TestHandler{
		logger: log,
		config: config,
	}

	api := app.Group("/test")
	api.Get("/test-server", testHandler.TestServer)
	api.Get("/test-habr", testHandler.TestHabr)
}

// TestHabr godoc
// @Summary Test connection
// @Description Test connection to HABR
// @Tags test
// @Success 200 {string} string
// @Failure 500 {string} string "Internal server error"
// @Router /test/test-habr [get]
func (h *TestHandler) TestHabr(c fiber.Ctx) error {
	response, err := http.Get(h.config.BaseUrl)
	if err != nil || response.StatusCode != 200 {
		h.logger.Printf("Error connecting to HABR: %v", err)
		c.Response().Header.SetStatusCode(http.StatusInternalServerError)
		return c.SendString("Lost connection to HABR")
	}

	c.Response().Header.SetStatusCode(http.StatusOK)
	return c.SendString("Test OK")
}

// TestServer godoc
// @Summary Test connection
// @Description Test connection to server
// @Tags test
// @Success 200 {string} string
// @Router /test/test-server [get]
func (h *TestHandler) TestServer(c fiber.Ctx) error {
	c.Response().Header.SetStatusCode(http.StatusOK)
	return c.SendString("Test OK")
}
